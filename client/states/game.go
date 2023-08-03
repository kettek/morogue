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
	binds     clgame.Binds
	scroller  clgame.Scroller
	grid      clgame.Grid
	mover     clgame.Mover
	sounds    clgame.Sounds
	inventory clgame.Inventory
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

	// Request any tiles that we don't have.
	requestedTiles := make(map[id.UUID]struct{})
	for _, c := range m.Cells {
		for _, cell := range c {
			if cell.TileID != nil {
				if _, ok := requestedTiles[*cell.TileID]; !ok {
					if _, ok := state.data.tiles[*cell.TileID]; !ok {
						state.connection.Write(net.TileMessage{
							ID: *cell.TileID,
						})
					}
					requestedTiles[*cell.TileID] = struct{}{}
				}
			}
		}
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
		case net.TileMessage:
			if m.ResultCode == 200 {
				// Store the tile and request the image.
				state.data.tiles[m.Tile.ID] = m.Tile
				if _, ok := state.data.tileImages[m.Tile.ID]; !ok {
					src := "archetypes/" + m.Tile.Image
					if img, err := state.data.LoadImage(src, ctx.Game.Zoom); err == nil {
						state.data.tileImages[m.Tile.ID] = img
					}
				}
			}
		case net.ArchetypesMessage:
			for _, a := range m.Archetypes {
				state.data.archetypes[a.GetID()] = a
				if _, err := state.data.EnsureImage(a, ctx.Game.Zoom); err != nil {
					fmt.Println("Error loading archetype image:", err)
				}
			}
			// This isn't exactly efficient.
			state.inventory.Refresh(ctx)
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
				state.inventory.SetInventory(&character.Inventory)
				state.inventory.Refresh(ctx)
			}
		case net.OwnerMessage:
			state.characterWID = m.WID
			if character := state.Character(); character != nil {
				state.ensureObjects(m.Inventory)
				character.Inventory = m.Inventory

				character.Skills = m.Skills
				state.inventory.SetInventory(&character.Inventory)
				state.inventory.Refresh(ctx)
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
					ch.Apply(o)
					if ch == state.Character() {
						fmt.Println("You applied an item")
						state.inventory.Refresh(ctx)
					} else {
						fmt.Printf("%s applied an item\n", ch.Name)
					}
				}
			}
		}
	case game.EventPickup:
		if o := state.location.ObjectByWID(evt.WID); o != nil {
			if picker := state.location.ObjectByWID(evt.Picker); picker != nil {
				o.SetPosition(game.Position{X: -1, Y: -1}) // FIXME: This feels wrong to use -1, -1 to signify hidden.
				if ch, ok := picker.(*game.Character); ok {
					ch.Pickup(o)
					if ch == state.Character() {
						fmt.Println("You picked up an item")
						state.inventory.Refresh(ctx)
					} else {
						fmt.Printf("%s picked up an item\n", ch.Name)
					}
				}
			}
		}
	case game.EventDrop:
		if o := state.location.ObjectByWID(evt.WID); o != nil {
			o.SetPosition(evt.Position) // Set the object's position to the dropped position.
			if dropper := state.location.ObjectByWID(evt.Dropper); dropper != nil {
				if ch, ok := dropper.(*game.Character); ok {
					ch.Drop(o)
					if ch == state.Character() {
						fmt.Println("You dropped an item")
						state.inventory.Refresh(ctx)
					} else {
						fmt.Printf("%s dropped an item\n", ch.Name)
					}
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

				if img := state.data.tileImages[*cell.TileID]; img != nil {
					opts := ebiten.DrawImageOptions{}
					opts.GeoM.Concat(scrollOpts)
					opts.GeoM.Translate(float64(px), float64(py))
					ctx.Screen.DrawImage(img, &opts)
				}
			}
		}
		// Draw characters
		for _, char := range state.location.Characters() {
			if img := state.data.archetypeImages[char.Archetype]; img != nil {
				px := char.X * cw
				py := char.Y * ch
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
