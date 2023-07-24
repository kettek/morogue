package states

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/game"
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
	//
	archetypesList *widget.List
	//
	archetypes []archetype
}

type archetype struct {
	Archetype game.Archetype
	Image     *ebiten.Image
}

func NewCreate(connection net.Connection, msgCh chan net.Message) *Create {
	state := &Create{
		connection:  connection,
		messageChan: msgCh,
		ui: &ebitenui.UI{
			Container: widget.NewContainer(
				widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{0x22, 0x13, 0x1a, 0xff})),
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

func (state *Create) haveArchetype(archetype game.Archetype) bool {
	for _, arch := range state.archetypes {
		if arch.Archetype.Title == archetype.Title {
			return true
		}
	}
	return false
}

func (state *Create) acquireArchetypes(archetypes []game.Archetype) {
	for _, arch := range archetypes {
		if state.haveArchetype(arch) {
			continue
		}

		var arche archetype
		arche.Archetype = arch

		defer func() {
			state.archetypes = append(state.archetypes, arche)
		}()

		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/archetypes/"+arch.Image, nil)
		if err != nil {
			log.Println(err)
			// TODO: Show error image?
			continue
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			// TODO: Show error image?
			continue
		}

		if res.StatusCode != 200 {
			// TODO: Show error image?
			continue
		}

		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			// TODO: Show error image?
			continue
		}

		img, _, err := image.Decode(bytes.NewReader(resBody))
		if err != nil {
			// TODO: Show error image?
			continue
		}

		ebiImg := ebiten.NewImageFromImage(img)

		arche.Image = ebiImg
		fmt.Println("assigned imagie")
	}
}

func (state *Create) refreshArchetypes() {
	// FIXME: Archetypes need to be displayed as selectable rows with the following columns: Image | Title | Swole | Zooms | Brains | Funk
	buttonImage := &widget.ButtonImage{
		Idle:    eimage.NewNineSliceColor(color.RGBA{R: 170, G: 170, B: 180, A: 255}),
		Hover:   eimage.NewNineSliceColor(color.RGBA{R: 130, G: 130, B: 150, A: 255}),
		Pressed: eimage.NewNineSliceColor(color.RGBA{R: 100, G: 100, B: 120, A: 255}),
	}

	for _, arch := range state.archetypes {
		buttonStackedLayout := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewStackedLayout()),
			// instruct the container's anchor layout to center the button both horizontally and vertically;
			// since our button is a 2-widget object, we add the anchor info to the wrapping container
			// instead of the button
			widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			})),
		)
		// construct a pressable button
		button := widget.NewButton(
			// specify the images to use
			widget.ButtonOpts.Image(buttonImage),

			// add a handler that reacts to clicking the button
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				println("button clicked")
			}),
		)
		buttonStackedLayout.AddChild(button)
		// Put an image on top of the button, it will be centered.
		// If your image doesn't fit the button and there is no Y stretching support,
		// you may see a transparent rectangle inside the button.
		// To fix that, either use a separate button image (that can fit the image)
		// or add an appropriate stretching.
		buttonStackedLayout.AddChild(widget.NewGraphic(widget.GraphicOpts.Image(arch.Image)))

		state.ui.Container.AddChild(buttonStackedLayout)
	}

}

func (state *Create) Update(ctx ifs.RunContext) error {
	select {
	case msg := <-state.messageChan:
		switch m := msg.(type) {
		case net.ArchetypesMessage:
			state.acquireArchetypes(m.Archetypes)
			state.refreshArchetypes()
		}
	default:
	}

	state.ui.Update()

	return nil
}

func (state *Create) Draw(ctx ifs.DrawContext) {
	state.ui.Draw(ctx.Screen)
}
