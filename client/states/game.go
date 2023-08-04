package states

import (
	"fmt"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/morogue/client/ifs"
	clgame "github.com/kettek/morogue/client/states/game"
	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/id"
	"github.com/kettek/morogue/net"
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
	characterWID id.WID
	//
	binds    clgame.Binds
	scroller clgame.Scroller
	grid     clgame.Grid
	mover    clgame.Mover
	sounds   clgame.Sounds
	//
	inventory clgame.Inventory
	below     clgame.Below
	hotbar    clgame.Hotbar
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
	}

	/*state.innerContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout()
		),
	)*/

	state.binds.Init()
	state.mover.Init()
	state.scroller.Init()
	state.scroller.SetHandler(func(x, y int) {
		state.grid.SetOffset(x, y)
		state.sounds.SetOffset(x, y)
	})
	state.grid.SetColor(color.NRGBA{255, 255, 255, 15})

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

	// TODO: Hide this.
	inventoryContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	inventoryContainerInner := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(8)),
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
	inventoryContainer.AddChild(inventoryContainerInner)
	state.inventory.Init(inventoryContainerInner, ctx)

	belowContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	belowContainerInner := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(8)),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal:  false,
				StretchVertical:    false,
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
			}),
		),
	)
	belowContainer.AddChild(belowContainerInner)
	state.below.Init(belowContainerInner, ctx)

	hotbarContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	hotbarContainerInner := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(8)),
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
	hotbarContainer.AddChild(hotbarContainerInner)
	state.hotbar.Init(hotbarContainerInner, ctx, &state.binds)

	state.ui.Container.AddChild(inventoryContainer)
	state.ui.Container.AddChild(belowContainer)
	state.ui.Container.AddChild(hotbarContainer)

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

func (state *Game) ensureObjects(objects game.Objects) {
	var missingArchetypes []id.UUID
	for _, obj := range objects {
		if _, ok := state.data.archetypes[obj.GetArchetype()]; !ok {
			missingArchetypes = append(missingArchetypes, obj.GetArchetype())
		}
	}
	if len(missingArchetypes) > 0 {
		state.connection.Write(net.ArchetypesMessage{
			IDs: missingArchetypes,
		})
	}
}

func (state *Game) Update(ctx ifs.RunContext) error {
	state.ui.Update()
	select {
	case msg := <-state.messageChan:
		switch m := msg.(type) {
		case net.LocationMessage:
			state.setLocationFromMessage(m)
			state.travelTo(m.ID)
			state.syncUIToLocation(ctx)
		case net.ArchetypeMessage:
			fmt.Println(msg)
		case net.ArchetypesMessage:
			for _, a := range m.Archetypes {
				state.data.archetypes[a.GetID()] = a
				if _, err := state.data.EnsureImage(a, ctx.Game.Zoom); err != nil {
					fmt.Println("Error loading archetype image:", err)
				}
			}
			// This isn't exactly efficient.
			state.refreshInventory(ctx)
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
				}
				state.refreshInventory(ctx)
			}
		case net.OwnerMessage:
			state.characterWID = m.WID
			if character := state.Character(); character != nil {
				state.ensureObjects(m.Inventory)

				character.Inventory = m.Inventory
				// Transform the inventory to reference our local objects.
				for i, o := range character.Inventory {
					realObj := state.location.ObjectByWID(o.GetWID())
					character.Inventory[i] = realObj
				}
				state.refreshInventory(ctx)

				character.Skills = m.Skills
			}
		default:
			fmt.Println("TODO: Handle", m.Type())
		}
	default:
	}
	state.binds.Update(ctx)

	state.scroller.Update(ctx)

	state.hotbar.Update(ctx)

	state.sounds.Update()

	if state.location != nil {
		if character := state.Character(); character != nil {
			if dir := state.mover.Update(state.binds); dir != 0 {
				state.sendDesire(state.characterWID, game.DesireMove{
					Direction: dir,
				})
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
		} else {
			state.location.Objects.RemoveByWID(evt.Object.GetWID())
			state.location.Objects.Add(evt.Object)
		}
	case game.EventRemove:
		state.location.Objects.RemoveByWID(evt.WID)
	case game.EventPosition:
		if c := state.location.Character(evt.WID); c != nil {
			c.X = evt.X
			c.Y = evt.Y
		}
	case game.EventSound:
		state.sounds.Add(evt.Message, evt.Position, evt.FromPosition)
	case game.EventNotice:
		fmt.Println(evt.Message)
	case game.EventApply:
		if o := state.location.ObjectByWID(evt.WID); o != nil {
			if applier := state.location.ObjectByWID(evt.Applier); applier != nil {
				if ch, ok := applier.(*game.Character); ok {
					if evt.Applied {
						ch.Apply(o)
						if ch == state.Character() {
							fmt.Println("You applied an item")
							state.refreshInventory(ctx)
						} else {
							fmt.Printf("%s applied an item\n", ch.Name)
						}
					} else {
						ch.Unapply(o)
						if ch == state.Character() {
							fmt.Println("You unapplied an item")
							state.refreshInventory(ctx)
						} else {
							fmt.Printf("%s unapplied an item\n", ch.Name)
						}
					}
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
				} else {
					fmt.Printf("%s dropped an item\n", ch.Name)
				}
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

func (state *Game) Draw(ctx ifs.DrawContext) {
	scrollOpts := ebiten.GeoM{}
	sx, sy := state.scroller.Scroll()
	scrollOpts.Translate(float64(sx), float64(sy))

	if state.location != nil {
		state.grid.Draw(ctx)

		cw := int(float64(ctx.Game.CellWidth) * ctx.Game.Zoom)
		ch := int(float64(ctx.Game.CellHeight) * ctx.Game.Zoom)

		// TODO: Replace map access and cells with a client-centric cell wrapper that contains the ebiten.image ptr directly for efficiency.
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

				if img := state.data.archetypeImages[*cell.TileID]; img != nil {
					opts := ebiten.DrawImageOptions{}
					opts.GeoM.Concat(scrollOpts)
					opts.GeoM.Translate(float64(px), float64(py))
					ctx.Screen.DrawImage(img, &opts)
				}
			}
		}
		// Draw characters
		for _, o := range state.location.Objects {
			// Ignore out of bounds objects.
			pos := o.GetPosition()
			if pos.X < 0 || pos.Y < 0 {
				continue
			}
			if img := state.data.archetypeImages[o.GetArchetype()]; img != nil {
				px := pos.X * cw
				py := pos.Y * ch
				opts := ebiten.DrawImageOptions{}
				opts.GeoM.Concat(scrollOpts)
				opts.GeoM.Translate(float64(px), float64(py))
				ctx.Screen.DrawImage(img, &opts)
			}
		}

		state.sounds.Draw(ctx)
	}

	state.ui.Draw(ctx.Screen)
}

func ToFloat64(a int, b int) (float64, float64) {
	return float64(a), float64(b)
}
