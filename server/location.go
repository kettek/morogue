package main

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/gen"
	"github.com/kettek/morogue/id"
	"github.com/kettek/morogue/net"
)

type location struct {
	game.Location
	active bool
}

func (l *location) addCharacter(character *game.Character) error {
	if l.hasCharacter(character.WID) {
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
	l.Characters = append(l.Characters, character)
	return nil
}

func (l *location) hasCharacter(wid id.WID) bool {
	for _, char := range l.Characters {
		if char.WID == wid {
			return true
		}
	}
	return false
}

func (l *location) removeCharacter(wid id.WID) error {
	for i, char := range l.Characters {
		if char.WID == wid {
			l.Characters = append(l.Characters[:i], l.Characters[i+1:]...)
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
	return nil
}

func (l *location) generate() error {
	// FIXME: Pull from somewhere.
	for x := 0; x < 10; x++ {
		l.Cells = append(l.Cells, make([]game.Cell, 0))
		for y := 0; y < 10; y++ {
			l.Cells[x] = append(l.Cells[x], game.Cell{})
		}
	}

	gen.Generate(gen.Styles[gen.StyleBox], gen.Config{
		Width:  10,
		Height: 10,
		Cell: func(x, y int) gen.Cell {
			return &l.Cells[x][y]
		},
		SetCell: func(x, y int, cell gen.Cell) {
			if cell.Flags().Has("blocked") {
				l.Cells[x][y].Blocks = game.MovementAll
				if tid, err := id.UID(id.Tile, "stone-wall"); err == nil {
					l.Cells[x][y].TileID = &tid
				}
			}
		},
	})
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
	for _, c := range l.Characters {
		if c.Desire != nil {
			switch d := c.Desire.(type) {
			case game.DesireMove:
				if err := l.moveCharacter(c.WID, d.Direction); err == nil {
					events = append(events, game.EventPosition{
						WID: c.WID,
						X:   c.X,
						Y:   c.Y,
					})
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
									FromX:   c.X,
									FromY:   c.Y,
									X:       x,
									Y:       y,
									Message: "*bump*",
								})
							} else if err == game.ErrOutOfBoundCell {
								events = append(events, game.EventSound{
									FromX:   c.X,
									FromY:   c.Y,
									X:       x,
									Y:       y,
									Message: "*pmub*",
								})
							}
						}
					}
				}
			}
			c.LastDesire = c.Desire
			c.Desire = nil
		}
	}
	for _, m := range l.Mobs {
		fmt.Println("TODO: Handle mob", m)
	}
	for _, o := range l.Objects {
		fmt.Println("TODO: Handle object", o)
	}
	return events
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
