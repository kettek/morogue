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
	"github.com/nfnt/resize"
	"golang.org/x/exp/slices"
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
	archetypesContainer *widget.Container
	//
	archetypes []archetype
	//
	face font.Face
	//
	tooltips map[string]*widget.Container
	//
	sortBy string
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
		tooltips: make(map[string]*widget.Container),
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
	state.face = face

	state.archetypesContainer = widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.RowLayoutOpts.Spacing(1),
		)),
	)

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
	state.ui.Container.AddChild(state.archetypesContainer)
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

		// Resize the image to 2x until ebitenui has scaling built-in.
		img = resize.Resize(uint(img.Bounds().Dx()*2), uint(img.Bounds().Dy()*2), img, resize.NearestNeighbor)

		ebiImg := ebiten.NewImageFromImage(img)

		arche.Image = ebiImg
	}
}

func (state *Create) refreshArchetypes() {
	state.archetypesContainer.RemoveChildren()

	// Heading
	{
		row := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			)),
		)

		el := widget.NewText(
			widget.TextOpts.Text("", state.face, color.White),
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
			widget.TextOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(50, 20),
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
			),
		)
		row.AddChild(el)

		parts := []string{"Archetype", "Swole", "Zooms", "Brains", "Funk"}
		for _, p := range parts {
			func(p string) {
				var c color.NRGBA
				var tooltip string
				switch p {
				case "Archetype":
					c = color.NRGBA{255, 255, 255, 255}
				case "Swole":
					c = game.ColorSwoleVibrant
					tooltip = game.AttributeSwoleDescription
				case "Zooms":
					c = game.ColorZoomsVibrant
					tooltip = game.AttributeZoomsDescription
				case "Brains":
					c = game.ColorBrainsVibrant
					tooltip = game.AttributeBrainsDescription
				case "Funk":
					c = game.ColorFunkVibrant
					tooltip = game.AttributeFunkDescription
				}

				tool := widget.NewTextToolTip(tooltip, state.face, color.White, eimage.NewNineSliceColor(color.NRGBA{R: 50, G: 50, B: 50, A: 255}))
				tool.Position = widget.TOOLTIP_POS_CURSOR_STICKY
				tool.Delay = time.Duration(time.Millisecond * 200)

				el := widget.NewText(
					widget.TextOpts.Text(p, state.face, c),
					widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
					widget.TextOpts.WidgetOpts(
						widget.WidgetOpts.MinSize(100, 20),
						widget.WidgetOpts.LayoutData(widget.RowLayoutData{
							Position: widget.RowLayoutPositionCenter,
						}),
						widget.WidgetOpts.ToolTip(tool),
						widget.WidgetOpts.MouseButtonPressedHandler(func(args *widget.WidgetMouseButtonPressedEventArgs) {
							state.sortBy = p
							state.refreshArchetypes()
						}),
					),
				)
				row.AddChild(el)
			}(p)
		}

		state.archetypesContainer.AddChild(row)
	}

	buttonImage := &widget.ButtonImage{
		Idle:    eimage.NewNineSliceColor(color.RGBA{R: 170, G: 170, B: 180, A: 255}),
		Hover:   eimage.NewNineSliceColor(color.RGBA{R: 130, G: 130, B: 150, A: 255}),
		Pressed: eimage.NewNineSliceColor(color.RGBA{R: 100, G: 100, B: 120, A: 255}),
	}

	var rowButtons []widget.RadioGroupElement

	slices.SortFunc(state.archetypes, func(a, b archetype) bool {
		if state.sortBy == "Archetype" {
			return a.Archetype.Title > b.Archetype.Title
		} else if state.sortBy == "Swole" {
			return a.Archetype.Swole > b.Archetype.Swole
		} else if state.sortBy == "Zooms" {
			return a.Archetype.Zooms > b.Archetype.Zooms
		} else if state.sortBy == "Brains" {
			return a.Archetype.Brains > b.Archetype.Brains
		} else if state.sortBy == "Funk" {
			return a.Archetype.Funk > b.Archetype.Funk
		}
		return false
	})

	for _, arch := range state.archetypes {
		func(arch archetype) {
			rowContainer := widget.NewContainer(
				widget.ContainerOpts.Layout(widget.NewStackedLayout()),
				widget.ContainerOpts.WidgetOpts(
					widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
						StretchHorizontal: true,
						StretchVertical:   true,
					}),
				),
			)
			rowContainerButton := widget.NewButton(
				widget.ButtonOpts.Image(buttonImage),

				// add a handler that reacts to clicking the button
				widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
					println("button clicked", arch.Archetype.Title)
				}),
			)
			rowButtons = append(rowButtons, rowContainerButton)
			rowContainer.AddChild(rowContainerButton)

			row := widget.NewContainer(
				widget.ContainerOpts.Layout(widget.NewRowLayout(
					widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
				)),
			)

			rowContainer.AddChild(row)

			graphicContainer := widget.NewContainer(
				widget.ContainerOpts.Layout(widget.NewStackedLayout()),
				widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{0, 0, 0, 255})),
				widget.ContainerOpts.WidgetOpts(
					widget.WidgetOpts.MinSize(50, 20),
					widget.WidgetOpts.LayoutData(widget.RowLayoutData{
						Stretch: true,
					}),
				),
			)

			graphic := widget.NewGraphic(
				widget.GraphicOpts.Image(arch.Image),
				widget.GraphicOpts.WidgetOpts(
					widget.WidgetOpts.LayoutData(
						widget.RowLayoutPositionCenter,
					),
				),
			)
			graphicContainer.AddChild(graphic)

			makeWidget := func(name string, clr color.NRGBA, tooltip string) *widget.Container {
				tool := widget.NewTextToolTip(tooltip, state.face, color.White, eimage.NewNineSliceColor(color.NRGBA{R: 50, G: 50, B: 50, A: 255}))
				tool.Position = widget.TOOLTIP_POS_CURSOR_STICKY
				tool.Delay = time.Duration(time.Millisecond * 200)

				box := widget.NewContainer(
					widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(clr)),
					widget.ContainerOpts.Layout(widget.NewStackedLayout()),
					widget.ContainerOpts.WidgetOpts(
						widget.WidgetOpts.MinSize(100, 20),
						widget.WidgetOpts.LayoutData(widget.RowLayoutData{
							Stretch: true,
						}),
						widget.WidgetOpts.ToolTip(tool),
					))
				content := widget.NewText(
					widget.TextOpts.Text(name, state.face, color.White),
					widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
					widget.TextOpts.WidgetOpts(
						widget.WidgetOpts.LayoutData(widget.RowLayoutData{
							Position: widget.RowLayoutPositionCenter,
						}),
					),
				)
				box.AddChild(content)
				return box
			}

			name := widget.NewText(
				widget.TextOpts.Text(arch.Archetype.Title, state.face, color.White),
				widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
				widget.TextOpts.WidgetOpts(
					widget.WidgetOpts.MinSize(100, 20),
					widget.WidgetOpts.LayoutData(widget.RowLayoutData{
						Position: widget.RowLayoutPositionCenter,
					}),
				),
			)

			swole := makeWidget(fmt.Sprintf("%d", int(arch.Archetype.Swole)), game.ColorSwole, game.AttributeSwoleDescription)

			zooms := makeWidget(fmt.Sprintf("%d", int(arch.Archetype.Zooms)), game.ColorZooms, game.AttributeZoomsDescription)

			brains := makeWidget(fmt.Sprintf("%d", int(arch.Archetype.Brains)), game.ColorBrains, game.AttributeBrainsDescription)

			funk := makeWidget(fmt.Sprintf("%d", int(arch.Archetype.Funk)), game.ColorFunk, game.AttributeFunkDescription)

			row.AddChild(graphicContainer)
			row.AddChild(name)
			row.AddChild(swole)
			row.AddChild(zooms)
			row.AddChild(brains)
			row.AddChild(funk)

			state.archetypesContainer.AddChild(rowContainer)
		}(arch)
	}

	widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(rowButtons...),
	)

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
