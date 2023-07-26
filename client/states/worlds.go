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

type Worlds struct {
	connection  net.Connection
	messageChan chan net.Message
	ui          *ebitenui.UI
	backButton  *widget.Button
}

func NewWorlds(connection net.Connection, msgCh chan net.Message) *Worlds {
	state := &Worlds{
		connection:  connection,
		messageChan: msgCh,
		ui: &ebitenui.UI{
			Container: widget.NewContainer(
				widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x22, 0x13, 0x1a, 0xff})),
				widget.ContainerOpts.Layout(widget.NewRowLayout(
					widget.RowLayoutOpts.Direction(widget.DirectionVertical),
					widget.RowLayoutOpts.Spacing(20),
					widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)))),
			),
		},
	}
	return state
}

func (state *Worlds) Begin(ctx ifs.RunContext) error {
	state.backButton = widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.CursorHovered("interactive"),
		),
		widget.ButtonOpts.Image(ctx.UI.ButtonImage),
		widget.ButtonOpts.Text("back", ctx.UI.HeadlineFace, ctx.UI.ButtonTextColor),
		widget.ButtonOpts.TextPadding(ctx.UI.ButtonPadding),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			state.connection.Write(net.UnjoinCharacterMessage{})
		}),
	)

	state.ui.Container.AddChild(state.backButton)
	return nil
}

func (state *Worlds) Return(interface{}) error {
	return nil
}

func (state *Worlds) Leave() error {
	return nil
}

func (state *Worlds) End() (interface{}, error) {
	return nil, nil
}

func (state *Worlds) Update(ctx ifs.RunContext) error {
	select {
	case msg := <-state.messageChan:
		switch m := msg.(type) {
		case net.UnjoinCharacterMessage:
			if m.ResultCode == 200 {
				ctx.Sm.Pop()
				return nil
			}
		case net.WorldsMessage:
			fmt.Println("populate worlds", m)
		default:
			fmt.Println(m)
		}
		fmt.Println("got eem", msg)
	default:
	}

	state.ui.Update()

	return nil
}

func (state *Worlds) Draw(ctx ifs.DrawContext) {
	state.ui.Draw(ctx.Screen)
}
