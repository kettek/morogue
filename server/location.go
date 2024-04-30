package server

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/gen"
	"github.com/kettek/morogue/id"
	"github.com/kettek/morogue/net"
)

type location struct {
	game.Location
	playerCharacters   []*game.Character // List of active player characters
	active             bool
	removable          bool // destroyable is used to allow a location to be removed.
	emptySince         time.Time
	turnCount          int  // Current turn count.
	turnActionCount    int  // Actions that have been taken by players this turn. OR, if the location is not in turns, a monotonic increment.
	turnActionLatch    int  // Trigger for processing a complete turn.
	turnActionOOCLatch int  // The latch for actions out of combat. This is generally equal to a second or 20 calls to process.
	inTurns            bool // Whether or not the location is currently processing the world in turns.
}

func newLocation() *location {
	return &location{
		turnActionOOCLatch: 20,
	}
}

func (l *location) addObject(o game.Object) {
	l.Objects = append(l.Objects, o)
}

func (l *location) removeObject(o game.Object) {
	for i, o2 := range l.Objects {
		if o2.GetWID() == o.GetWID() {
			l.Objects = append(l.Objects[:i], l.Objects[i+1:]...)
			return
		}
	}
}

func (l *location) addCharacter(character *game.Character) error {
	if l.Character(character.WID) != nil {
		return ErrCharacterAlreadyInLocation
	}

	// Find an open spot for them in wherever is appropriate.
	openCells := l.filterCells(func(c game.Cell) bool {
		return c.Blocks == game.MovementNone
	})
	if len(openCells) == 0 {
		return ErrCharacterCannotPlaceInLocation
	}
	// TODO: We want to also handle non-random spawning, such as through stairs, etc.
	// Add to location.
	spawnCell := openCells[rand.Intn(len(openCells)-1)]
	character.X = spawnCell.X
	character.Y = spawnCell.Y

	l.addObject(character)

	// Add the character's inventory.
	for _, o := range character.Inventory {
		l.addObject(o)
	}

	l.active = true

	// Add the character to the list of player characters.
	l.playerCharacters = append(l.playerCharacters, character)

	// Increase the turn latch.
	l.turnActionLatch += 1

	return nil
}

func (l *location) removeCharacter(wid id.WID) error {
	for i, o := range l.Objects {
		if char, ok := o.(*game.Character); ok && char.WID == wid {
			l.Objects = append(l.Objects[:i], l.Objects[i+1:]...)

			// Remove the character's inventory.
			for _, o := range char.Inventory {
				l.removeObject(o)
			}

			if len(l.Characters()) == 0 {
				l.active = false
				l.emptySince = time.Now()
			}

			// Remove the character from the list of player characters.
			for i, c := range l.playerCharacters {
				if c.WID == wid {
					l.playerCharacters = append(l.playerCharacters[:i], l.playerCharacters[i+1:]...)
					break
				}
			}

			// Decreate the turn latch.
			l.turnActionLatch -= 1

			return nil
		}
	}
	return ErrCharacterNotInLocation
}

func (l *location) moveCharacter(wid id.WID, dir game.MoveDirection) error {
	ch := l.Character(wid)
	if ch == nil {
		return ErrCharacterNotInLocation
	}
	x, y := dir.Position()
	x += ch.X
	y += ch.Y

	if cell, err := l.Cells.At(x, y); err != nil {
		return err
	} else if cell.Blocks == game.MovementAll {
		return ErrMovementBlocked
	}

	ch.X = x
	ch.Y = y

	// FIXME: This isn't the right place for this. There should be some sort of "actions" economy that is used to increase hunger.
	ch.Movable.MoveCounter++
	if ch.Movable.MoveCounter > 10 { // I guess 10 steps are reasonable enough for energy checks.
		ch.Movable.MoveCounter = 0
		energy := 1 + int(ch.Attributes.Swole/2) - int(ch.Attributes.Zooms/4)
		if energy < 1 {
			energy = 1
		}
		ch.UseEnergy(energy)
	}

	return nil
}

type wfcTile struct {
	ID     id.UUID
	Domain []id.UUID
}

func (l *location) generate(pid id.UUID, data *Data, wids *id.WIDGenerator) error {
	place, err := data.Places.ById(pid)
	if err != nil {
		return err
	}

	allPossibleTiles := []id.UUID{}

	for _, w := range place.WFC {
		allPossibleTiles = append(allPossibleTiles, w.ID)
	}

	w := place.Width.Roll()
	h := place.Height.Roll()

	var wfcTiles [][]wfcTile

	/*removeDomain := func(x, y int, id id.UUID) {
		if x < 0 || y < 0 || x >= len(wfcTiles) || y >= len(wfcTiles[0]) {
			return
		}
		for i, d := range wfcTiles[x][y].Domain {
			if d == id {
				wfcTiles[x][y].Domain = append(wfcTiles[x][y].Domain[:i], wfcTiles[x][y].Domain[i+1:]...)
				return
			}
		}
	}*/

	getLeastEntropy := func() (ex, ey int) {
		leastEntropy := 9999
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				if len(wfcTiles[x][y].Domain) > 0 {
					if len(wfcTiles[x][y].Domain) < leastEntropy {
						leastEntropy = len(wfcTiles[x][y].Domain)
						ex = x
						ey = y
					}
				}
			}
		}
		if leastEntropy == 9999 {
			return -1, -1
		}
		return
	}

	for x := 0; x < w; x++ {
		l.Cells = append(l.Cells, make([]game.Cell, 0))
		wfcTiles = append(wfcTiles, make([]wfcTile, 0))
		for y := 0; y < h; y++ {
			l.Cells[x] = append(l.Cells[x], game.Cell{})
			wfcTiles[x] = append(wfcTiles[x], wfcTile{
				Domain: allPossibleTiles,
			})
		}
	}

	placeFixture := func(f gen.Fixture, px, py int) error {
		w := f.Width()
		h := f.Height()

		if px+w > len(l.Cells) || py+h > len(l.Cells[0]) {
			return game.ErrOutOfBoundCell
		}

		for y, r := range f.Rows {
			for x, c := range r {
				if cid, ok := f.Keys[string(c)]; ok {
					l.Cells[px+x][py+y].TileID = &cid
					// Mark the tile as done for WFC purposes.
					wfcTiles[px+x][py+y].ID = cid
					wfcTiles[px+x][py+y].Domain = []id.UUID{}
				}
			}
		}
		return nil
	}

	for _, f := range place.Fixtures {
		count := f.Count.Roll()
		for i := 0; i < count; i++ {
			for tries := 0; tries < 10; tries++ {
				// TODO: Make this weighted.
				target := f.Targets[rand.Intn(len(f.Targets))]
				fixture, err := data.Fixtures.ById(target.ID)
				if err != nil {
					return err
				}
				var x int
				var y int
				if f.X.Max() == 0 {
					x = rand.Intn(w)
				} else {
					x = f.X.Roll()
				}
				if f.Y.Max() == 0 {
					y = rand.Intn(h)
				} else {
					y = f.Y.Roll()
				}
				if x >= 0 && x < w && y >= 0 && y < h {
					if err := placeFixture(fixture, x, y); err != nil {
						log.Println(errors.Join(err, fmt.Errorf("%s attempt %d", target.ID, tries)))
					} else {
						break
					}
				} else if tries+1 >= 10 {
					return fmt.Errorf("could not place fixture %s", target.ID)
				}
			}
		}
	}

	for {
		x, y := getLeastEntropy()
		if x == -1 && y == -1 {
			break
		}

		wfcTiles[x][y].ID = wfcTiles[x][y].Domain[rand.Intn(len(wfcTiles[x][y].Domain))]
		wfcTiles[x][y].Domain = []id.UUID{}
		l.Cells[x][y].TileID = &wfcTiles[x][y].ID
	}

	return nil
}

type cellLocation struct {
	X, Y int
	Cell game.Cell
}

func (l *location) filterCells(cb func(c game.Cell) bool) (cells []cellLocation) {
	for x, r := range l.Cells {
		for y, c := range r {
			if cb(c) {
				cells = append(cells, cellLocation{
					X:    x,
					Y:    y,
					Cell: c,
				})
			}
		}
	}
	return
}

func (l *location) process() (events []game.Event) {
	for _, c := range l.playerCharacters {
		events = append(events, l.processCharacter(c)...)
	}

	if !l.inTurns {
		l.turnActionCount++
	}

	if (!l.inTurns && l.turnActionCount >= l.turnActionOOCLatch) || (l.inTurns && l.turnActionCount >= l.turnActionLatch) {
		// TODO: Process non-player characters.
		for _, o := range l.Objects {
			switch o.(type) {
			default:
				// TODO
			}
		}

		l.turnActionCount = 0
		l.turnCount++

		// Only send turn events if we're actually in what we consider to be turns.
		if l.inTurns {
			events = append(events, game.EventTurn{
				Turn: l.turnCount,
			})
		}
	}
	return events
}

func (l *location) processCharacter(c *game.Character) (events []game.Event) {
	if c.Desire != nil {
		if l.inTurns {
			if c.SpentActions >= c.Actions {
				return nil
			}
			if c.SpentActions < c.Actions {
				c.SpentActions++
			}
			if c.SpentActions >= c.Actions {
				l.turnActionCount++
			}
		}
		switch d := c.Desire.(type) {
		case game.DesireMove:
			if err := l.moveCharacter(c.WID, d.Direction); err == nil {
				events = append(events, game.EventPosition{
					WID:      c.WID,
					Position: c.Position,
				})
				if c.MoveCounter == 0 {
					events = append(events, game.EventHunger{
						WID:    c.WID,
						Hunger: c.Hungerable.Hunger,
					})
				}
			} else {
				// Make bump sounds if the character is moving in the same direction as their last desire.
				if last, ok := c.LastDesire.(game.DesireMove); ok {
					if last.Direction == d.Direction {
						// Make the sfx come from the cell bumped into.
						x, y := d.Direction.Position()
						x += c.X
						y += c.Y
						if err == ErrMovementBlocked {
							events = append(events, game.EventSound{
								FromPosition: c.Position,
								Position:     game.Position{X: x, Y: y},
								Message:      "*bump*",
							})
						} else if err == game.ErrOutOfBoundCell {
							events = append(events, game.EventSound{
								FromPosition: c.Position,
								Position:     game.Position{X: x, Y: y},
								Message:      "*pmub*",
							})
						}
					}
				}
			}
		case game.DesireApply:
			if t := l.ObjectByWID(d.WID); t != nil {
				// The separation of edible from appliable logic is a bit annoying since they share so much.
				if _, isAppliable := t.(Appliable); isAppliable {
					var e game.Event
					if d.Apply {
						e = c.Apply(t, false)
					} else {
						e = c.Unapply(t, false)
					}
					if _, ok := e.(game.EventNotice); ok {
						c.Events = append(c.Events, e)
					} else if e != nil {
						events = append(events, e)
					}
				} else if _, isEdible := t.(Edible); isEdible {
					e := c.Apply(t, true)
					if _, ok := e.(game.EventNotice); ok {
						c.Events = append(c.Events, e)
					} else if e != nil {
						events = append(events, e)
					}
					if e, ok := e.(game.EventConsume); ok {
						if e.Finished {
							events = append(events, game.EventSound{
								FromPosition: c.GetPosition(),
								Position:     c.GetPosition(),
								Message:      "*burp*",
							})
							// Destroy the item.
							events = append(events, l.DestroyObject(t))
						} else {
							// TODO: Make eating foods make different sounds, such as "slurp", "crunch", etc.
							events = append(events, game.EventSound{
								FromPosition: c.GetPosition(),
								Position:     c.GetPosition(),
								Message:      "*munch*",
							})
						}
					}
				} else {
					c.Events = append(c.Events, game.EventNotice{
						Message: "You can't apply that.",
					})
				}
			}
		case game.DesirePickup:
			if t := l.ObjectByWID(d.WID); t != nil {
				if t.GetContainerWID() > 0 {
					c.Events = append(c.Events, game.EventNotice{
						Message: "You can't pick that up.",
					})
				} else {
					if t.GetPosition() != c.GetPosition() {
						c.Events = append(c.Events, game.EventNotice{
							Message: "You can't reach that.",
						})
					} else {
						events = append(events, c.Pickup(t))
						events = append(events, game.EventSound{
							FromPosition: c.Position,
							Position:     c.Position,
							Message:      "*snarf*",
						})
					}
				}
			}
		case game.DesireDrop:
			if t := l.ObjectByWID(d.WID); t != nil {
				e := c.Drop(t)
				if _, ok := e.(game.EventNotice); ok {
					c.Events = append(c.Events, e)
				} else {
					events = append(events, e)
					if e, ok := e.(game.EventDrop); ok {
						t.SetPosition(e.Position)
						events = append(events, game.EventSound{
							FromPosition: c.Position,
							Position:     c.Position,
							Message:      "*whump*",
						})
					}
				}
			}
		case game.DesireBash:
			if t := l.ObjectByWID(d.WID); t != nil {
				if hurtable, ok := t.(Hurtable); ok {
					// TODO: Maybe only take unarmed damage?
					damages := c.RollDamages()
					hurtable.TakeDamages(damages)
					events = append(events, game.EventDamages{
						From:    c.WID,
						Target:  d.WID,
						Damages: damages,
					})
					events = append(events, game.EventSound{
						FromPosition: c.Position,
						Position:     t.GetPosition(),
						Message:      "*thud*",
					})
				}
			} else {
				c.Events = append(c.Events, game.EventNotice{
					Message: "You kick at the air.",
				})
			}
		case game.DesireOpen:
			if t := l.ObjectByWID(d.WID); t != nil {
				if openable, isOpenable := t.(Openable); isOpenable {
					lockable, isLockable := t.(Lockable)
					if openable.IsOpened() {
						if err := openable.Close(); err == nil {
							events = append(events, game.EventSound{
								FromPosition: c.Position,
								Position:     t.GetPosition(),
								Message:      "*click*",
							})
							c.Events = append(c.Events, game.EventNotice{
								Message: "You close the door.",
							})
						} else if errors.Is(err, game.ErrAlreadyClosed) {
							c.Events = append(c.Events, game.EventNotice{
								Message: "It's already closed.",
							})
						}
					} else if !openable.IsOpened() {
						if isLockable && lockable.IsLocked() {
							c.Events = append(c.Events, game.EventNotice{
								Message: "It's locked.",
							})
						} else if err := openable.Open(); err == nil {
							c.Events = append(c.Events, game.EventNotice{
								Message: "You open the door.",
							})
							events = append(events, game.EventSound{
								FromPosition: c.Position,
								Position:     t.GetPosition(),
								Message:      "*creak*",
							})
						} else if errors.Is(err, game.ErrAlreadyOpen) {
							c.Events = append(c.Events, game.EventNotice{
								Message: "It's already open.",
							})
						}
					}
				}
			} else {
				c.Events = append(c.Events, game.EventNotice{
					Message: "There is nothing there to open.",
				})
			}
		}
		c.LastDesire = c.Desire
		c.Desire = nil
	}
	return
}

// startTurns is called when the location should start processing the world in terms of turns. This should be done when the players begin combat.
func (l *location) startTurns() {
	l.inTurns = true
	l.turnCount = 0
	l.turnActionCount = 0
	for _, c := range l.playerCharacters {
		c.SpentActions = 0
	}
}

// stopTurns is called when the location should stop processing the world in terms of turns. This should be called when the players end combat.
func (l *location) stopTurns() {
	l.turnCount = 0
	l.turnActionCount = 0
	l.inTurns = false
}

// DestroyObject removes the given object from the location and any container it may be in. EventRemove is returned, which should be sent from the server to the client.
func (l *location) DestroyObject(target game.Object) game.Event {
	if container := target.GetContainerWID(); container > 0 {
		if c := l.Character(container); c != nil {
			c.Drop(target)
		}
	}
	l.removeObject(target)
	return game.EventRemove{
		WID: target.GetWID(),
	}
}

// TODO: Allow returning world-munging requests.
func (l *location) handleClientMessage(cl *client, msg net.Message) error {
	switch m := msg.(type) {
	case net.DesireMessage:
		if ch := l.Character(cl.currentCharacter.WID); ch != nil {
			ch.Desire = m.Desire.Desire()
		} else {
			cl.conn.Write(net.DesireMessage{
				ResultCode: 404,
				Result:     ErrCharacterNotInLocation.Error(),
				WID:        m.WID,
			})
		}
	case net.InventoryMessage:
		// Allow sending for the character's own inventory.
		if cl.currentCharacter.WID == m.WID {
			cl.conn.Write(net.InventoryMessage{
				WID:       m.WID,
				Inventory: cl.currentCharacter.Inventory,
			})
		}
	case net.SkillsMessage:
		if cl.currentCharacter.WID == m.WID {
			cl.conn.Write(net.SkillsMessage{
				WID:    m.WID,
				Skills: cl.currentCharacter.Skills,
			})
		}
	}
	return nil
}

type locationConfig struct {
	ID    id.UUID
	Depth int
}

var (
	ErrCharacterNotInLocation         = errors.New("character is not in location")
	ErrCharacterAlreadyInLocation     = errors.New("character is already in location")
	ErrCharacterCannotPlaceInLocation = errors.New("character cannot be placed in location")
)

var (
	ErrMovementBlocked = errors.New("movement blocked")
)
