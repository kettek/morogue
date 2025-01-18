package states

import (
	"fmt"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/goro/pathing"
	"github.com/kettek/morogue/client/ifs"
	clgame "github.com/kettek/morogue/client/states/game"
	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/id"
	"github.com/kettek/morogue/locale"
	"github.com/kettek/morogue/net"
	"github.com/tinne26/etxt"
	"github.com/vmihailenco/msgpack/v5"
)

// Game represents the running game in the world.
type Game struct {
	data        *Data
	connection  net.Connection
	messageChan chan net.Message
	//
	ui             *ebitenui.UI
	innerContainer *widget.Container
	//
	locations map[id.UUID]*game.Location
	location  *game.Location // current
	//
	characterWID          id.WID
	objectsMissingArchs   []id.WID
	lockCameraToCharacter bool
	showGrid              bool
	pendingDesire         game.Desire
	pathingSteps          []pathing.Step
	//
	binds    clgame.Binds
	scroller clgame.Scroller
	grid     clgame.Grid
	actioner clgame.Actioner
	pather   clgame.Pather
	pinger   clgame.Pinger
	sounds   clgame.Sounds
	kickers  clgame.Kickers
	//
	inventory clgame.Inventory
	below     clgame.Below
	hotbar    clgame.Hotbar
	statbar   clgame.Statbar
	lc        locale.Localizer
}

// NewGame creates a new Game instance.
func NewGame(connection net.Connection, msgCh chan net.Message, data *Data) *Game {
	state := &Game{
		data:        data,
		connection:  connection,
		messageChan: msgCh,
		ui: &ebitenui.UI{
			Container: widget.NewContainer(
				widget.ContainerOpts.Layout(widget.NewStackedLayout()),
			),
		},
		locations: make(map[id.UUID]*game.Location),
		lc:        locale.Get(locale.Locale()),
	}

	/*state.innerContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout()
		),
	)*/

	state.binds.Init()
	state.actioner.Init()
	state.scroller.Init()
	state.scroller.SetHandler(func(x, y int) {
		state.grid.SetOffset(x, y)
		state.sounds.SetOffset(x, y)
		state.kickers.SetOffset(x, y)
		state.pather.SetOffset(x, y)
		state.pinger.SetOffset(x, y)
	})
	state.grid.SetColor(color.NRGBA{255, 255, 255, 30})
	state.grid.SetHeldHandler(func(x, y int) {
		path := pathing.NewPathFromFunc(len(state.location.Cells), len(state.location.Cells[0]), func(x, y int) uint32 {
			cell := state.location.Cells[x][y]
			if cell.Blocks == game.MovementWalk || cell.Blocks == game.MovementAll {
				return pathing.MaximumCost
			}
			return 0
		}, pathing.AlgorithmAStar)
		path.AllowDiagonals(true)
		state.pather.Steps = path.Compute(state.Character().X, state.Character().Y, x, y)
	})
	state.pinger.Init(nil, ifs.RunContext{})
	state.pinger.PingLocation = func(x, y int, kind string) {
		state.sendDesire(state.characterWID, game.DesirePing{
			Position: game.Position{
				X: x,
				Y: y,
			},
			Kind: kind,
		})
	}

	state.inventory.Data = data
	state.inventory.DropItem = func(wid id.WID) {
		state.sendDesire(state.characterWID, game.DesireDrop{
			WID: wid,
		})
	}
	state.inventory.ApplyItem = func(wid id.WID, apply bool) {
		state.sendDesire(state.characterWID, game.DesireApply{
			WID:   wid,
			Apply: apply,
		})
	}
	state.inventory.PickupItem = func(wid id.WID) {
		state.sendDesire(state.characterWID, game.DesirePickup{
			WID: wid,
		})
	}

	state.below.Data = data
	state.below.PickupItem = func(wid id.WID) {
		state.sendDesire(state.characterWID, game.DesirePickup{
			WID: wid,
		})
	}
	state.below.ApplyItem = func(wid id.WID) {
		state.sendDesire(state.characterWID, game.DesireApply{
			WID: wid,
		})
	}
	state.below.DropItem = func(wid id.WID) {
		state.sendDesire(state.characterWID, game.DesireDrop{
			WID: wid,
		})
	}

	return state
}

func (state *Game) Begin(ctx ifs.RunContext) error {
	cw := int(float64(ctx.Game.CellWidth) * ctx.Game.Zoom)
	ch := int(float64(ctx.Game.CellHeight) * ctx.Game.Zoom)
	state.grid.SetCellSize(cw, ch)

	{
		invBelowContainer := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		)

		invBelowContainerInner := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(widget.Insets{Top: 100}),
			)),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					StretchHorizontal:  false,
					StretchVertical:    false,
					HorizontalPosition: widget.AnchorLayoutPositionEnd,
					VerticalPosition:   widget.AnchorLayoutPositionStart,
				}),
			),
		)

		inventoryContainer := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		)
		inventoryContainerInner := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout(
				widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(8)),
			)),
		)
		inventoryContainer.AddChild(inventoryContainerInner)
		state.inventory.Init(inventoryContainerInner, ctx)

		belowContainer := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		)
		belowContainerInner := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout(
				widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(8)),
			)),
		)
		belowContainer.AddChild(belowContainerInner)
		state.below.Init(belowContainerInner, ctx, &state.binds)

		invBelowContainerInner.AddChild(inventoryContainer)
		invBelowContainerInner.AddChild(belowContainer)

		invBelowContainer.AddChild(invBelowContainerInner)

		state.ui.Container.AddChild(invBelowContainer)
	}

	{
		hotbarAndStatbarContainer := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		)
		hotbarAndStatbarContainerInner := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			)),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					StretchHorizontal:  false,
					StretchVertical:    false,
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionEnd,
				}),
			),
		)
		hotbarAndStatbarContainer.AddChild(hotbarAndStatbarContainerInner)

		state.ui.Container.AddChild(hotbarAndStatbarContainer)

		hotbarContainer := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					StretchHorizontal:  true,
					StretchVertical:    false,
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionEnd,
				}),
			),
		)
		hotbarContainerInner := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout(
				widget.AnchorLayoutOpts.Padding(widget.Insets{Bottom: 10}),
			)),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					StretchHorizontal:  true,
					StretchVertical:    false,
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionEnd,
				}),
			),
		)
		hotbarContainer.AddChild(hotbarContainerInner)
		state.hotbar.Init(hotbarContainerInner, ctx, &state.binds)

		hotbarAndStatbarContainerInner.AddChild(hotbarContainer)

		statbarContainer := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		)
		statbarContainerInner := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					StretchHorizontal:  true,
					StretchVertical:    false,
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionEnd,
				}),
			),
		)
		statbarContainer.AddChild(statbarContainerInner)
		state.statbar.Init(statbarContainerInner, ctx)

		hotbarAndStatbarContainerInner.AddChild(statbarContainer)
	}

	return nil
}

func (state *Game) Return(interface{}) error {
	return nil
}

func (state *Game) Leave() error {
	return nil
}

func (state *Game) End() (interface{}, error) {
	return nil, nil
}

func (state *Game) setLocationFromMessage(m net.LocationMessage) {
	l := &game.Location{
		ID:      m.ID,
		Cells:   m.Cells,
		Objects: m.Objects,
	}

	// Request objects we don't have.
	state.ensureObjects(m.Objects)

	// Request tile objects.
	var missingTiles []id.UUID
	missingTiles2 := make(map[id.UUID]bool)
	for _, r := range l.Cells {
		for _, c := range r {
			if c.TileID == nil {
				continue
			}

			if _, ok := missingTiles2[*c.TileID]; !ok {
				missingTiles2[*c.TileID] = true
				missingTiles = append(missingTiles, *c.TileID)
			}
		}
	}

	if len(missingTiles) > 0 {
		state.connection.Write(net.ArchetypesMessage{
			IDs: missingTiles,
		})
	}

	state.locations[m.ID] = l
}

func (state *Game) travelTo(id id.UUID) {
	l, ok := state.locations[id]
	if !ok {
		// TODO: Maybe request the location?
		return
	}
	state.location = l
}

func (state *Game) syncUIToLocation(ctx ifs.RunContext) {
	if state.location == nil {
		return
	}
	w := int(float64(len(state.location.Cells)*ctx.Game.CellWidth) * ctx.Game.Zoom)
	h := int(float64(len(state.location.Cells[0])*ctx.Game.CellHeight) * ctx.Game.Zoom)
	state.scroller.SetLimit(w, h)
	state.scroller.CenterTo(ctx.UI.Width/2-w/2, ctx.UI.Height/2-h/2)
	state.grid.SetSize(w, h)
}

func (state *Game) assignObjectArchetype(obj game.Object) {
	if a := state.data.archetypes[obj.GetArchetypeID()]; a == nil {
		state.objectsMissingArchs = append(state.objectsMissingArchs, obj.GetWID())
	} else {
		obj.SetArchetype(a)
	}
}

func (state *Game) centerCameraOn(ctx ifs.RunContext, o game.Object) {
	pos := o.GetPosition()
	x := -int(float64(pos.X*ctx.Game.CellWidth) * ctx.Game.Zoom)
	y := -int(float64(pos.Y*ctx.Game.CellHeight) * ctx.Game.Zoom)
	state.scroller.CenterTo(x, y)
}

func (state *Game) ensureObjects(objects game.Objects) {
	var missingArchetypes []id.UUID
	for _, obj := range objects {
		if _, ok := state.data.archetypes[obj.GetArchetypeID()]; !ok {
			missingArchetypes = append(missingArchetypes, obj.GetArchetypeID())
		}
		state.assignObjectArchetype(obj)
	}
	if len(missingArchetypes) > 0 {
		state.connection.Write(net.ArchetypesMessage{
			IDs: missingArchetypes,
		})
	}
}

func (state *Game) Update(ctx ifs.RunContext) error {
	state.ui.Update()
	state.grid.Update(ctx)
	state.pinger.Update(ctx)
	select {
	case msg := <-state.messageChan:
		switch m := msg.(type) {
		case net.LocationMessage:
			state.setLocationFromMessage(m)
			state.travelTo(m.ID)
			state.syncUIToLocation(ctx)
			if ch := state.Character(); ch != nil {
				state.centerCameraOn(ctx, ch)
			}
		case net.ArchetypeMessage:
			fmt.Println(msg)
		case net.ArchetypesMessage:
			for _, a := range m.Archetypes {
				state.data.archetypes[a.GetID()] = a
				if _, err := state.data.EnsureImage(a, ctx.Game.Zoom); err != nil {
					fmt.Println("Error loading archetype image:", err)
				}
				// Update any objects awaiting this archetype.
				i := 0
				for _, wid := range state.objectsMissingArchs {
					if o := state.location.ObjectByWID(wid); o != nil {
						if o.GetArchetypeID() == a.GetID() {
							state.assignObjectArchetype(o)
							continue
						}
					}
					state.objectsMissingArchs[i] = wid
					i++
				}
				state.objectsMissingArchs = state.objectsMissingArchs[:i]
			}
			// This isn't exactly efficient.
			state.refreshInventory(ctx)
			state.refreshStatbar(ctx)
		case net.EventsMessage:
			for _, evt := range m.Events {
				state.handleEvent(evt.Event(), ctx)
			}
		case net.EventMessage:
			state.handleEvent(m.Event.Event(), ctx)
		case net.InventoryMessage:
			if character := state.Character(); character != nil {
				state.ensureObjects(m.Inventory)

				character.Inventory = m.Inventory
				// Transform the inventory to reference our local objects.
				for i, o := range character.Inventory {
					realObj := state.location.ObjectByWID(o.GetWID())
					character.Inventory[i] = realObj
					// Also assign the object's container to be the player.
					o.SetContainerWID(state.characterWID)
				}
				state.refreshInventory(ctx)
			}
		case net.OwnerMessage:
			state.characterWID = m.WID
			if character := state.Character(); character != nil {
				state.ensureObjects(m.Inventory)

				character.Attributes = m.Attributes

				character.Inventory = m.Inventory
				// Transform the inventory to reference our local objects.
				for i, o := range character.Inventory {
					realObj := state.location.ObjectByWID(o.GetWID())
					character.Inventory[i] = realObj
					// Also assign the object's container to be the player.
					o.SetContainerWID(state.characterWID)
				}
				state.refreshInventory(ctx)
				state.refreshStatbar(ctx)

				character.Skills = m.Skills

				// Center the camera on the character.
				state.centerCameraOn(ctx, character)
			}
		default:
			fmt.Println("TODO: Handle", m.Type())
		}
	default:
	}
	state.binds.Update(ctx)

	state.scroller.Update(ctx)

	state.hotbar.Update(ctx)
	state.below.Update(ctx)

	state.sounds.Update()
	state.kickers.Update()

	if state.location != nil {
		if character := state.Character(); character != nil {
			if state.binds.IsActionHeld("lock-camera") == 0 {
				state.lockCameraToCharacter = !state.lockCameraToCharacter
			}
			if state.binds.IsActionHeld("snap-camera") >= 0 {
				state.centerCameraOn(ctx, character)
			}
			if state.binds.IsActionHeld("toggle-grid") == 0 {
				state.showGrid = !state.showGrid
			}
			if desire := state.actioner.Update(state.binds); desire != nil {
				state.sendDesire(state.characterWID, desire)
				state.pather.Steps = nil
			} else if desire := state.pather.Update(character); desire != nil {
				state.sendDesire(state.characterWID, desire)
			}

			// This isn't great, but whatever.
			var belowObjects game.Objects
			for _, o := range state.location.Objects {
				if o == character {
					continue
				}
				if o.GetPosition() == character.GetPosition() {
					belowObjects = append(belowObjects, o)
				}
			}
			state.below.Refresh(ctx, belowObjects)
		}
	}

	return nil
}

func (state *Game) handleEvent(e game.Event, ctx ifs.RunContext) {
	if state.location == nil {
		return
	}
	switch evt := e.(type) {
	case game.EventAdd:
		if ch := state.location.ObjectByWID(evt.Object.GetWID()); ch == nil {
			state.location.Objects.Add(evt.Object)
			state.assignObjectArchetype(evt.Object)
		} else {
			state.location.Objects.RemoveByWID(evt.Object.GetWID())
			state.location.Objects.Add(evt.Object)
			state.assignObjectArchetype(evt.Object)
		}
	case game.EventRemove:
		if o := state.location.ObjectByWID(evt.WID); o != nil {
			if co := o.GetContainerWID(); co > 0 {
				if c := state.location.ObjectByWID(co); c != nil {
					switch c := c.(type) {
					case *game.Character:
						// Just use drop to remove from character.
						c.Drop(o)
						// Refresh the inventory if its the current player.
						if c == state.Character() {
							state.refreshInventory(ctx)
						}
					}
				}
			}
		}
		state.location.Objects.RemoveByWID(evt.WID)
	case game.EventPosition:
		if c := state.location.Character(evt.WID); c != nil {
			c.X = evt.X
			c.Y = evt.Y
			if c == state.Character() && state.lockCameraToCharacter {
				pos := c.GetPosition()
				x := -int(float64(pos.X*ctx.Game.CellWidth) * ctx.Game.Zoom)
				y := -int(float64(pos.Y*ctx.Game.CellHeight) * ctx.Game.Zoom)
				state.scroller.CenterTo(x, y)
			}
		}
	case game.EventSound:
		state.sounds.Add(evt.Message, evt.Position, evt.FromPosition)
	case game.EventNotice:
		fmt.Println(state.lc.T(evt.Message, evt.Args))
		fmt.Println(evt.Message, evt.Args, state.lc.T(evt.Message, evt.Args))
	case game.EventApply:
		if o := state.location.ObjectByWID(evt.WID); o != nil {
			if applier := state.location.ObjectByWID(evt.Applier); applier != nil {
				if ch, ok := applier.(*game.Character); ok {
					if evt.Applied {
						ch.Apply(o, true)
						if ch == state.Character() {
							fmt.Println("You applied an item")
							state.refreshInventory(ctx)
							state.refreshStatbar(ctx)
						} else {
							fmt.Printf("%s applied an item\n", ch.Name)
						}
					} else {
						ch.Unapply(o, true)
						if ch == state.Character() {
							fmt.Println("You unapplied an item")
							state.refreshInventory(ctx)
							state.refreshStatbar(ctx)
						} else {
							fmt.Printf("%s unapplied an item\n", ch.Name)
						}
					}
				}
			}
		}
	case game.EventConsume:
		if o := state.location.ObjectByWID(evt.WID); o != nil {
			if consumer := state.location.ObjectByWID(evt.Consumer); consumer != nil {
				if ch, ok := consumer.(*game.Character); ok {
					if food, ok := o.GetArchetype().(game.FoodArchetype); ok {
						ch.Apply(o, true)
						if ch == state.Character() {
							fmt.Printf("You consumed %s\n", food.Title)
							state.refreshInventory(ctx)
							state.refreshStatbar(ctx)
						} else {
							fmt.Printf("%s consumed %s\n", ch.Name, food.Title)
						}
					}
				}
			}
		}
	case game.EventHunger:
		if o := state.location.ObjectByWID(evt.WID); o != nil {
			if ch := state.location.Character(evt.WID); ch != nil {
				ch.Hunger = evt.Hunger
				if ch == state.Character() {
					state.refreshStatbar(ctx)
				}
			}
		}
	case game.EventPickup:
		if o := state.location.ObjectByWID(evt.WID); o != nil {
			if picker := state.location.ObjectByWID(evt.Picker); picker != nil {
				if ch, ok := picker.(*game.Character); ok {
					ch.Pickup(o)
					if ch == state.Character() {
						fmt.Println("You picked up an item")
						state.refreshInventory(ctx)
					} else {
						fmt.Printf("%s picked up an item\n", ch.Name)
					}
				}
			}
		}
	case game.EventDrop:
		if o := state.location.ObjectByWID(evt.Object.GetWID()); o == nil {
			state.location.Objects.Add(evt.Object)
		} else {
			state.location.Objects.RemoveByWID(evt.Object.GetWID())
			state.location.Objects.Add(evt.Object)
		}

		o := state.location.ObjectByWID(evt.Object.GetWID())
		o.SetPosition(evt.Position) // Set the object's position to the dropped position.
		if dropper := state.location.ObjectByWID(evt.Dropper); dropper != nil {
			if ch, ok := dropper.(*game.Character); ok {
				ch.Drop(o)
				if ch == state.Character() {
					fmt.Println("You dropped an item")
					state.refreshInventory(ctx)
					state.refreshStatbar(ctx)
				} else {
					fmt.Printf("%s dropped an item\n", ch.Name)
				}
			}
		}
	case game.EventTurn:
		fmt.Println("turn", evt.Turn)
	case game.EventPing:
		state.pinger.Add(evt.Position, evt.Kind)
	case game.EventDamages:
		fmt.Println("TODO: Handle damages", evt)
	case game.EventHealth:
		if o := state.location.ObjectByWID(evt.Target); o != nil {
			if ch := state.location.Character(evt.Target); ch != nil {
				diff := evt.Health - ch.Health
				ch.Health = evt.Health
				if ch == state.Character() {
					state.refreshStatbar(ctx)
				}
				msg := fmt.Sprintf("%d", diff)
				clr := color.NRGBA{255, 0, 0, 255}
				if diff > 0 {
					msg = "+" + msg
					clr = color.NRGBA{0, 255, 0, 255}
				}
				state.kickers.Add(clgame.Kicker{
					Message: msg,
					Position: game.Position{
						X: ch.X,
						Y: ch.Y,
					},
					Lifetime: 60,
					Color:    clr,
				})
			}
		}
	}
}

func (state *Game) sendDesire(wid id.WID, d game.Desire) error {
	b, err := msgpack.Marshal(d)
	if err != nil {
		return err
	}

	state.connection.Write(net.DesireMessage{
		WID: wid,
		Desire: game.DesireWrapper{
			Type: d.Type(),
			Data: b,
		},
	})
	return nil
}

func (state *Game) Character() *game.Character {
	if state.location != nil {
		if character := state.location.Character(state.characterWID); character != nil {
			return character
		}
	}
	return nil
}

func (state *Game) refreshInventory(ctx ifs.RunContext) {
	state.inventory.Refresh(ctx, state.Character().Inventory)
}

func (state *Game) refreshStatbar(ctx ifs.RunContext) {
	if ch := state.Character(); ch != nil {
		if a := state.data.archetypes[state.Character().ArchetypeID]; a != nil {
			state.statbar.Refresh(ctx, ch, a.(game.CharacterArchetype))
		}
	}
}

func (state *Game) Draw(ctx ifs.DrawContext) {
	scrollOpts := ebiten.GeoM{}
	sx, sy := state.scroller.Scroll()
	scrollOpts.Translate(float64(sx), float64(sy))

	if state.location != nil {
		cw := int(float64(ctx.Game.CellWidth) * ctx.Game.Zoom)
		ch := int(float64(ctx.Game.CellHeight) * ctx.Game.Zoom)

		// TODO: Replace map access and cells with a client-centric cell wrapper that contains the ebiten.image ptr directly for efficiency.
		// Draw tiles.
		for x, col := range state.location.Cells {
			for y, cell := range col {
				if cell.TileID == nil {
					continue
				}
				px := x * cw
				py := y * ch

				if px+sx > ctx.UI.Width || py+sy > ctx.UI.Height || px+sx+cw < 0 || py+sy+ch < 0 {
					continue
				}

				// This is kind of goofy, but modulate the color of the tile based on its position.
				v := float32((x+y)%8) / 70.0

				if img := state.data.archetypeImages[*cell.TileID]; img != nil {
					opts := ebiten.DrawImageOptions{}
					opts.GeoM.Concat(scrollOpts)
					opts.GeoM.Translate(float64(px), float64(py))
					opts.ColorScale.Scale(1.0-v, 1.0-v, 1.0-v, 1.0)
					ctx.Screen.DrawImage(img, &opts)
				}
			}
		}

		// Draw grid
		if state.showGrid {
			state.grid.Draw(ctx)
		}

		// Draw pathing
		if state.Character() != nil {
			state.pather.Draw(ctx, state.Character())
		}

		// Draw characters
		ctx.Txt.Save()
		ctx.Txt.SetAlign(etxt.VertCenter | etxt.HorzCenter)
		for _, o := range state.location.Objects {
			// Ignore out of bounds objects.
			pos := o.GetPosition()
			if pos.X < 0 || pos.Y < 0 {
				continue
			}
			if img := state.data.archetypeImages[o.GetArchetypeID()]; img != nil {
				px := pos.X * cw
				py := pos.Y * ch
				opts := ebiten.DrawImageOptions{}
				opts.GeoM.Concat(scrollOpts)
				opts.GeoM.Translate(float64(px), float64(py))
				ctx.Screen.DrawImage(img, &opts)
			}
			switch o := o.(type) {
			case *game.Character:
				px := float64(pos.X*cw + cw/2)
				py := float64(pos.Y*ch - ch/5)
				ctx.Txt.SetOutlineColor(color.NRGBA{0, 0, 0, 128})
				ctx.Txt.SetColor(color.NRGBA{255, 255, 255, 128})
				ctx.Txt.DrawWithOutline(ctx.Screen, o.Name, int(px+float64(sx)), int(py+float64(sy)))
			}
		}
		ctx.Txt.Restore()

		state.sounds.Draw(ctx)
		state.kickers.Draw(ctx)
		state.pinger.Draw(ctx)
	}

	state.ui.Draw(ctx.Screen)
}

func ToFloat64(a int, b int) (float64, float64) {
	return float64(a), float64(b)
}
