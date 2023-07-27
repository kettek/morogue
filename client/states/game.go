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
)

// Game represents the running game in the world.
type Game struct {
	data        *Data
	connection  net.Connection
	messageChan chan net.Message
	//
	ui *ebitenui.UI
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
}

// NewGame creates a new Game instance.
func NewGame(connection net.Connection, msgCh chan net.Message, data *Data) *Game {
	state := &Game{
		data:        data,
		connection:  connection,
		messageChan: msgCh,
		ui: &ebitenui.UI{
			Container: widget.NewContainer(
				widget.ContainerOpts.Layout(widget.NewRowLayout(
					widget.RowLayoutOpts.Direction(widget.DirectionVertical),
					widget.RowLayoutOpts.Spacing(20),
					widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20))),
				),
			),
		},
		locations: make(map[id.UUID]*game.Location),
	}
	state.binds.Init()
	state.mover.Init()
	state.scroller.Init()
	state.scroller.SetHandler(func(x, y int) {
		state.grid.SetOffset(x, y)
	})
	state.grid.SetCellSize(16*2, 16*2)
	state.grid.SetColor(color.NRGBA{255, 255, 255, 15})

	return state
}

func (state *Game) Begin(ctx ifs.RunContext) error {
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
		ID:         m.ID,
		Cells:      m.Cells,
		Characters: m.Characters,
		Mobs:       m.Mobs,
		Objects:    m.Objects,
	}

	// Request any tiles that we don't have.
	requestedTiles := make(map[id.UUID]struct{})
	for _, c := range m.Cells {
		for _, cell := range c {
			if cell.TileID != nil {
				if _, ok := requestedTiles[*cell.TileID]; !ok {
					if _, ok := state.data.tiles[*cell.TileID]; !ok {
						fmt.Println("requesting", *cell.TileID)
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
	w := len(state.location.Cells) * 16 * 2
	h := len(state.location.Cells[0]) * 16 * 2
	state.scroller.SetLimit(w, h)
	state.scroller.CenterTo(ctx.UI.Width/2-w/2, ctx.UI.Height/2-h/2)
	state.grid.SetSize(w, h)
}

func (state *Game) Update(ctx ifs.RunContext) error {
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
					src := "tiles/" + m.Tile.Image
					if img, err := state.data.loadImage(src, 2.0); err == nil {
						state.data.tileImages[m.Tile.ID] = img
					}
				}
			}
		case net.PositionMessage:
			if ch := state.location.Character(m.WID); ch != nil {
				ch.X = m.X
				ch.Y = m.Y
			}
		case net.OwnerMessage:
			state.characterWID = m.WID
		default:
			fmt.Println("TODO: Handle", m)
		}
	default:
	}
	state.binds.Update(ctx)

	state.ui.Update()

	state.scroller.Update(ctx)

	if dir := state.mover.Update(state.binds); dir != 0 {
		state.connection.Write(net.MoveMessage{
			WID:       state.characterWID,
			Direction: dir,
		})
	}

	return nil
}

func (state *Game) Draw(ctx ifs.DrawContext) {
	scrollOpts := ebiten.GeoM{}
	sx, sy := state.scroller.Scroll()
	scrollOpts.Translate(float64(sx), float64(sy))

	if state.location != nil {
		state.grid.Draw(ctx)

		// TODO: Replace map access and cells with a client-centric cell wrapper that contains the ebiten.image ptr directly for efficiency.
		for x, col := range state.location.Cells {
			for y, cell := range col {
				if cell.TileID == nil {
					continue
				}
				px := x * 16 * 2
				py := y * 16 * 2

				if px+sx > ctx.UI.Width || py+sy > ctx.UI.Height || px+sx+16*2 < 0 || py+sy+16*2 < 0 {
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
		for _, ch := range state.location.Characters {
			if img := state.data.archetypeImages[ch.Archetype]; img != nil {
				px := ch.X * 16 * 2
				py := ch.Y * 16 * 2
				opts := ebiten.DrawImageOptions{}
				opts.GeoM.Concat(scrollOpts)
				opts.GeoM.Translate(float64(px), float64(py))
				ctx.Screen.DrawImage(img, &opts)
			}
		}
	}

	state.ui.Draw(ctx.Screen)
}

func ToFloat64(a int, b int) (float64, float64) {
	return float64(a), float64(b)
}
