package states

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/morogue/client/ifs"
	clgame "github.com/kettek/morogue/client/states/game"
	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/id"
	"github.com/kettek/morogue/net"
	"github.com/nfnt/resize"
)

// Game represents the running game in the world.
type Game struct {
	connection  net.Connection
	messageChan chan net.Message
	//
	ui *ebitenui.UI
	//
	locations map[id.UUID]*game.Location
	location  *game.Location // current
	//
	tiles      map[id.UUID]game.Tile
	tileImages map[id.UUID]*ebiten.Image
	//
	scroller clgame.Scroller
	grid     clgame.Grid
}

// NewGame creates a new Game instance.
func NewGame(connection net.Connection, msgCh chan net.Message) *Game {
	state := &Game{
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
		locations:  make(map[id.UUID]*game.Location),
		tiles:      make(map[id.UUID]game.Tile),
		tileImages: make(map[id.UUID]*ebiten.Image),
	}
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
					if _, ok := state.tiles[*cell.TileID]; !ok {
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
				state.tiles[m.Tile.ID] = m.Tile
				if _, ok := state.tileImages[m.Tile.ID]; !ok {
					src := "tiles/" + m.Tile.Image
					if img, err := state.loadImage(src, 2.0); err == nil {
						state.tileImages[m.Tile.ID] = img
					}
				}
			}
		default:
			fmt.Println("TODO: Handle", m)
		}
	default:
	}

	state.ui.Update()

	state.scroller.Update(ctx)

	return nil
}

func (state *Game) loadImage(src string, scale float64) (*ebiten.Image, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/"+src, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(resBody))
	if err != nil {
		return nil, err
	}

	// Resize the image to 2x until ebitenui has scaling built-in.
	img = resize.Resize(uint(float64(img.Bounds().Dx())*scale), uint(float64(img.Bounds().Dy())*scale), img, resize.NearestNeighbor)

	return ebiten.NewImageFromImage(img), nil
}

func (state *Game) Draw(ctx ifs.DrawContext) {
	scrollOpts := ebiten.GeoM{}
	scrollOpts.Translate(ToFloat64(state.scroller.Scroll()))

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
				if img := state.tileImages[*cell.TileID]; img != nil {
					opts := ebiten.DrawImageOptions{}
					opts.GeoM.Concat(scrollOpts)
					opts.GeoM.Translate(float64(px), float64(py))
					ctx.Screen.DrawImage(img, &opts)
				}
			}
		}
	}

	state.ui.Draw(ctx.Screen)
}

func ToFloat64(a int, b int) (float64, float64) {
	return float64(a), float64(b)
}
