package states

import (
	"fmt"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/net"
)

// Game represents the running game in the world.
type Game struct {
	connection  net.Connection
	messageChan chan net.Message
	//
	ui *ebitenui.UI
}

// NewGame creates a new Game instance.
func NewGame(connection net.Connection, msgCh chan net.Message) *Game {
	state := &Game{
		connection:  connection,
		messageChan: msgCh,
		ui: &ebitenui.UI{
			Container: widget.NewContainer(
				widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x22, 0x13, 0x1a, 0xff})),
				widget.ContainerOpts.Layout(widget.NewRowLayout(
					widget.RowLayoutOpts.Direction(widget.DirectionVertical),
					widget.RowLayoutOpts.Spacing(20),
					widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20))),
				),
			),
		},
	}
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

func (state *Game) Update(ctx ifs.RunContext) error {
	select {
	case msg := <-state.messageChan:
		switch m := msg.(type) {
		default:
			fmt.Println("TODO: Handle", m)
		}
	default:
	}

	state.ui.Update()

	return nil
}

func (state *Game) Draw(ctx ifs.DrawContext) {
	state.ui.Draw(ctx.Screen)
}
