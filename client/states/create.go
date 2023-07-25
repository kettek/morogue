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
	"github.com/kettek/morogue/id"
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
	charactersSection   *widget.Container
	charactersContainer *widget.Container
	charactersControls  *widget.Container
	deleteWindow        *widget.Window
	//
	archetypesSection      *widget.Container
	archetypesContainer    *widget.Container
	archetypesControls     *widget.Container
	archetypesCreateImage  *widget.Graphic
	archetypesCreateName   *widget.TextInput
	archetypesCreateButton *widget.Button
	//
	archetypes []archetype
	//
	face font.Face
	//
	tooltips map[string]*widget.Container
	//
	sortBy            string
	selectedArchetype id.UUID
	characterToDelete string
	//
	traitsImage    *ebiten.Image
	swoleImage     *ebiten.Image
	zoomsImage     *ebiten.Image
	brainsImage    *ebiten.Image
	funkImage      *ebiten.Image
	archetypeImage *ebiten.Image
}

type character struct {
	Character game.Character
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
	// Load attributes images.
	if img, err := state.loadImage("images/swole.png", 2); err != nil {
		return err
	} else {
		state.swoleImage = img
	}
	if img, err := state.loadImage("images/zooms.png", 2); err != nil {
		return err
	} else {
		state.zoomsImage = img
	}
	if img, err := state.loadImage("images/brains.png", 2); err != nil {
		return err
	} else {
		state.brainsImage = img
	}
	if img, err := state.loadImage("images/funk.png", 2); err != nil {
		return err
	} else {
		state.funkImage = img
	}
	if img, err := state.loadImage("images/traits.png", 2); err != nil {
		return err
	} else {
		state.traitsImage = img
	}
	if img, err := state.loadImage("images/archetype.png", 2); err != nil {
		return err
	} else {
		state.archetypeImage = img
	}

	// load images for button states: idle, hover, and pressed
	buttonImages, _ := buttonImages()

	// load button text font
	face, _ := opentype.NewFace(ctx.Txt.Renderer.GetFont(), &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	state.face = face

	state.charactersContainer = widget.NewContainer(
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

	deleteContents := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{100, 100, 100, 255})),
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

	deleteText := widget.NewText(
		widget.TextOpts.Text("Confirm deletion", face, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	buttonImage := &widget.ButtonImage{
		Idle:    eimage.NewNineSliceColor(color.RGBA{R: 40, G: 40, B: 40, A: 255}),
		Hover:   eimage.NewNineSliceColor(color.RGBA{R: 50, G: 50, B: 50, A: 255}),
		Pressed: eimage.NewNineSliceColor(color.RGBA{R: 80, G: 80, B: 80, A: 255}),
	}
	deleteButton := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CursorHovered("delete"),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
		),
		widget.ButtonOpts.Text("delete", state.face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xff, 0x0, 0x0, 0xff},
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			state.connection.Write(net.DeleteCharacterMessage{
				Name: state.characterToDelete,
			})
			state.deleteWindow.Close()
		}),
	)

	deleteContents.AddChild(deleteText)
	deleteContents.AddChild(deleteButton)

	state.deleteWindow = widget.NewWindow(
		widget.WindowOpts.Contents(deleteContents),
		widget.WindowOpts.Modal(),
		widget.WindowOpts.CloseMode(widget.CLICK_OUT),
		//Set the minimum size the window can be
		widget.WindowOpts.MinSize(400, 200),
		//Set the maximum size a window can be
		widget.WindowOpts.MaxSize(400, 200),
		//Set the callback that triggers when a move is complete
	)

	{
		state.archetypesSection = widget.NewContainer(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
			),
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			)),
		)

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
		state.archetypesSection.AddChild(state.archetypesContainer)

		state.archetypesControls = widget.NewContainer(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
			),
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
				widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(10)),
				widget.RowLayoutOpts.Spacing(10),
			)),
		)

		state.archetypesCreateImage = widget.NewGraphic(
			widget.GraphicOpts.Image(nil),
			widget.GraphicOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
					Stretch:  true,
				}),
			),
		)
		state.archetypesControls.AddChild(state.archetypesCreateImage)

		state.archetypesCreateName = widget.NewTextInput(
			widget.TextInputOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(200, 20),
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
					Stretch:  true,
				}),
				widget.WidgetOpts.CursorHovered("text"),
			),
			widget.TextInputOpts.Image(&widget.TextInputImage{
				Idle:     eimage.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
				Disabled: eimage.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
			}),
			widget.TextInputOpts.Face(face),
			widget.TextInputOpts.Color(&widget.TextInputColor{
				Idle:          color.NRGBA{254, 255, 255, 255},
				Disabled:      color.NRGBA{R: 200, G: 200, B: 200, A: 255},
				Caret:         color.NRGBA{254, 255, 255, 255},
				DisabledCaret: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
			}),
			widget.TextInputOpts.Padding(widget.NewInsetsSimple(5)),
			widget.TextInputOpts.CaretOpts(
				widget.CaretOpts.Size(face, 2),
			),
			widget.TextInputOpts.Placeholder("character name"),
			widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
				state.doCreate()
			}),
			//This is called whenver there is a change to the text
			widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
				// TODO: Store it.
			}),
		)
		state.archetypesControls.AddChild(state.archetypesCreateName)

		state.archetypesCreateButton = widget.NewButton(
			widget.ButtonOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				}),
				widget.WidgetOpts.CursorHovered("interactive"),
			),
			widget.ButtonOpts.Image(buttonImages),
			widget.ButtonOpts.Text("create", face, &widget.ButtonTextColor{
				Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
			}),
			widget.ButtonOpts.TextPadding(widget.Insets{
				Left:   30,
				Right:  30,
				Top:    5,
				Bottom: 5,
			}),
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				state.doCreate()
			}),
		)
		state.archetypesControls.AddChild(state.archetypesCreateButton)

		state.archetypesSection.AddChild(state.archetypesControls)
	}

	state.logoutButton = widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.CursorHovered("interactive"),
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
	state.ui.Container.AddChild(state.charactersContainer)
	state.ui.Container.AddChild(state.archetypesSection)
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
		if arch.Archetype.UUID == archetype.UUID {
			return true
		}
	}
	return false
}

func (state *Create) loadImage(src string, scale float64) (*ebiten.Image, error) {
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

func (state *Create) populateCharacters(characters []game.Character) {
	state.charactersContainer.RemoveChildren()

	buttonImage := &widget.ButtonImage{
		Idle:    eimage.NewNineSliceColor(color.RGBA{R: 40, G: 40, B: 40, A: 255}),
		Hover:   eimage.NewNineSliceColor(color.RGBA{R: 50, G: 50, B: 50, A: 255}),
		Pressed: eimage.NewNineSliceColor(color.RGBA{R: 80, G: 80, B: 80, A: 255}),
	}

	var rowButtons []widget.RadioGroupElement
	for _, ch := range characters {
		func(ch game.Character) {
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
					fmt.Println("use char", ch)
					// TODO
				}),
				widget.ButtonOpts.WidgetOpts(
					widget.WidgetOpts.CursorHovered("interactive"),
				),
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

			var img *ebiten.Image
			for _, arch := range state.archetypes {
				if arch.Archetype.UUID == ch.Archetype {
					img = arch.Image
					break
				}
			}

			graphic := widget.NewGraphic(
				widget.GraphicOpts.Image(img),
				widget.GraphicOpts.WidgetOpts(
					widget.WidgetOpts.LayoutData(
						widget.RowLayoutPositionCenter,
					),
				),
			)
			graphicContainer.AddChild(graphic)

			name := widget.NewText(
				widget.TextOpts.Text(ch.Name, state.face, color.White),
				widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
				widget.TextOpts.WidgetOpts(
					widget.WidgetOpts.MinSize(100, 20),
					widget.WidgetOpts.LayoutData(widget.RowLayoutData{
						Position: widget.RowLayoutPositionCenter,
					}),
				),
			)

			delete := widget.NewButton(
				widget.ButtonOpts.Image(buttonImage),
				widget.ButtonOpts.WidgetOpts(
					widget.WidgetOpts.CursorHovered("delete"),
					widget.WidgetOpts.LayoutData(widget.RowLayoutData{
						Stretch: true,
					}),
				),
				widget.ButtonOpts.Text("delete", state.face, &widget.ButtonTextColor{
					Idle: color.NRGBA{0xff, 0x0, 0x0, 0xff},
				}),
				widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
					x, y := state.deleteWindow.Contents.PreferredSize()
					r := image.Rect(0, 0, x, y)
					r = r.Add(image.Point{100, 50})
					state.deleteWindow.SetLocation(r)
					state.ui.AddWindow(state.deleteWindow)
					state.characterToDelete = ch.Name
				}),
			)

			row.AddChild(graphicContainer)
			row.AddChild(name)
			row.AddChild(delete)

			state.charactersContainer.AddChild(rowContainer)
		}(ch)
	}
	widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(rowButtons...),
	)
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

		img, err := state.loadImage("archetypes/"+arch.Image, 2.0)
		if err != nil {
			// TODO: Show error image
			continue
		}

		arche.Image = img
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

		parts := []string{"Archetype", "Swole", "Zooms", "Brains", "Funk", "Traits"}
		for _, p := range parts {
			func(p string) {
				var c color.NRGBA
				var tooltip string
				var img *ebiten.Image
				width := 100
				switch p {
				case "Archetype":
					c = color.NRGBA{255, 255, 255, 255}
					tooltip = "Archetype is a collection of attributes and traits"
					img = state.archetypeImage
				case "Swole":
					c = game.ColorSwoleVibrant
					tooltip = game.AttributeSwoleDescription
					img = state.swoleImage
				case "Zooms":
					c = game.ColorZoomsVibrant
					tooltip = game.AttributeZoomsDescription
					img = state.zoomsImage
				case "Brains":
					c = game.ColorBrainsVibrant
					tooltip = game.AttributeBrainsDescription
					img = state.brainsImage
				case "Funk":
					c = game.ColorFunkVibrant
					tooltip = game.AttributeFunkDescription
					img = state.funkImage
				case "Traits":
					c = color.NRGBA{200, 200, 200, 255}
					img = state.traitsImage
					tooltip = "Traits are unique modifiers to an archetype"
					width = 200
				}

				tool := widget.NewTextToolTip(tooltip, state.face, color.White, eimage.NewNineSliceColor(color.NRGBA{R: 50, G: 50, B: 50, A: 255}))
				tool.Position = widget.TOOLTIP_POS_CURSOR_STICKY
				tool.Delay = time.Duration(time.Millisecond * 200)

				container := widget.NewContainer(
					widget.ContainerOpts.Layout(widget.NewStackedLayout()),
					widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{0, 0, 0, 255})),
					widget.ContainerOpts.WidgetOpts(
						widget.WidgetOpts.MinSize(width, 40),
						widget.WidgetOpts.LayoutData(widget.RowLayoutData{
							Stretch: true,
						}),
						widget.WidgetOpts.ToolTip(tool),
						widget.WidgetOpts.CursorHovered("interactive-tooltip"),
						widget.WidgetOpts.MouseButtonPressedHandler(func(args *widget.WidgetMouseButtonPressedEventArgs) {
							state.sortBy = p
							state.refreshArchetypes()
						}),
					),
				)

				if img != nil {
					graphic := widget.NewGraphic(
						widget.GraphicOpts.Image(img),
						widget.GraphicOpts.WidgetOpts(
							widget.WidgetOpts.LayoutData(
								widget.RowLayoutPositionCenter,
							),
						),
					)
					container.AddChild(graphic)
				} else {
					el := widget.NewText(
						widget.TextOpts.Text(p, state.face, c),
						widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
						widget.TextOpts.WidgetOpts(
							widget.WidgetOpts.LayoutData(widget.RowLayoutData{
								Position: widget.RowLayoutPositionCenter,
							}),
						),
					)
					container.AddChild(el)
				}
				row.AddChild(container)
			}(p)
		}

		state.archetypesContainer.AddChild(row)
	}

	buttonImage := &widget.ButtonImage{
		Idle:    eimage.NewNineSliceColor(color.RGBA{R: 40, G: 40, B: 40, A: 255}),
		Hover:   eimage.NewNineSliceColor(color.RGBA{R: 50, G: 50, B: 50, A: 255}),
		Pressed: eimage.NewNineSliceColor(color.RGBA{R: 80, G: 80, B: 80, A: 255}),
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
					state.selectedArchetype = arch.Archetype.UUID
					state.archetypesCreateImage.Image = arch.Image
				}),
				widget.ButtonOpts.WidgetOpts(
					widget.WidgetOpts.CursorHovered("interactive"),
				),
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

			makeWidget := func(name string, clr color.NRGBA, width int) *widget.Container {
				box := widget.NewContainer(
					widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(clr)),
					widget.ContainerOpts.Layout(widget.NewStackedLayout()),
					widget.ContainerOpts.WidgetOpts(
						widget.WidgetOpts.MinSize(width, 20),
						widget.WidgetOpts.LayoutData(widget.RowLayoutData{
							Stretch: true,
						}),
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

			swole := makeWidget(fmt.Sprintf("%d", int(arch.Archetype.Swole)), game.ColorSwole, 100)

			zooms := makeWidget(fmt.Sprintf("%d", int(arch.Archetype.Zooms)), game.ColorZooms, 100)

			brains := makeWidget(fmt.Sprintf("%d", int(arch.Archetype.Brains)), game.ColorBrains, 100)

			funk := makeWidget(fmt.Sprintf("%d", int(arch.Archetype.Funk)), game.ColorFunk, 100)

			var d string
			for _, t := range arch.Archetype.Traits {
				d += t + "\n"
			}
			desc := makeWidget(d, color.NRGBA{32, 32, 32, 0}, 200)

			row.AddChild(graphicContainer)
			row.AddChild(name)
			row.AddChild(swole)
			row.AddChild(zooms)
			row.AddChild(brains)
			row.AddChild(funk)
			row.AddChild(desc)

			state.archetypesContainer.AddChild(rowContainer)
		}(arch)
	}

	widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(rowButtons...),
	)

}

func (state *Create) doCreate() {
	name := state.archetypesCreateName.InputText
	id := state.selectedArchetype

	if name == "" {
		state.resultText.Label = "name must not be empty"
		return
	}
	if id.IsNil() {
		state.resultText.Label = "archetype must be selected"
		return
	}

	state.connection.Write(net.CreateCharacterMessage{
		Name:      name,
		Archetype: id,
	})
}

func (state *Create) Update(ctx ifs.RunContext) error {
	select {
	case msg := <-state.messageChan:
		switch m := msg.(type) {
		case net.ArchetypesMessage:
			state.acquireArchetypes(m.Archetypes)
			state.refreshArchetypes()
		case net.CharactersMessage:
			state.populateCharacters(m.Characters)
		case net.CreateCharacterMessage:
			state.resultText.Label = m.Result
		case net.DeleteCharacterMessage:
			if m.Result != "" {
				state.resultText.Label = m.Result
			}
		}
	default:
	}

	state.ui.Update()

	return nil
}

func (state *Create) Draw(ctx ifs.DrawContext) {
	state.ui.Draw(ctx.Screen)
}
