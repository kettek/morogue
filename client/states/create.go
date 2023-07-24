package states

import (
	"fmt"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/net"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Create struct {
	connection  net.Connection
	messageChan chan net.Message
	ui          *ebitenui.UI
	//
	logoutButton *widget.Button
	resultText   *widget.Text
}

func NewCreate(connection net.Connection, msgCh chan net.Message) *Create {
	state := &Create{
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

func (state *Create) Begin(ctx ifs.RunContext) error {
	// load images for button states: idle, hover, and pressed
	buttonImages, _ := buttonImages()

	// load button text font
	face, _ := opentype.NewFace(ctx.Txt.Renderer.GetFont(), &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	state.logoutButton = widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImages),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("logout", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			state.connection.Write(net.LogoutMessage{})
			ctx.Sm.Pop()
		}),
	)

	state.resultText = widget.NewText(
		widget.TextOpts.Text("Create a new hero or select a previous one.", face, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)

	state.ui.Container.AddChild(state.resultText)
	state.ui.Container.AddChild(state.logoutButton)

	return nil
}

func (state *Create) Return(interface{}) error {
	return nil
}

func (state *Create) Leave() error {
	return nil
}

func (state *Create) End() (interface{}, error) {
	return nil, nil
}

func (state *Create) Update(ctx ifs.RunContext) error {
	select {
	case msg := <-state.messageChan:
		switch m := msg.(type) {
		case net.ArchetypesMessage:
			fmt.Println("Populate archetypes", m)
		}
	default:
	}

	state.ui.Update()

	return nil
}

func (state *Create) Draw(ctx ifs.DrawContext) {
	state.ui.Draw(ctx.Screen)
}
