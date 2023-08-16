package states

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"time"

	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/id"
	"github.com/kettek/morogue/net"
	"golang.org/x/exp/slices"
)

// Create is a rather massive and messy state that controls character selection, creation,
// and deletion. Popping Create should return to the Login state.
type Create struct {
	data *Data
	//
	connection  net.Connection
	messageChan chan net.Message
	ui          *ebitenui.UI
	//
	logoutButton *widget.Button
	resultText   *widget.Text
	//
	charactersSection      *widget.Container
	charactersContainer    *widget.Container
	charactersRadioGroup   *widget.RadioGroup
	charactersControls     *widget.Container
	charactersJoinButton   *widget.Button
	charactersDeleteButton *widget.Button
	deleteWindow           *widget.Window
	//
	archetypesSection      *widget.Container
	archetypesContainer    *widget.Container
	archetypesRadioGroup   *widget.RadioGroup
	archetypesControls     *widget.Container
	archetypesCreateImage  *widget.Graphic
	archetypesCreateName   *widget.TextInput
	archetypesCreateButton *widget.Button
	//
	characters []*game.Character
	archetypes []archetype
	//
	tooltips map[string]*widget.Container
	//
	sortBy            string
	selectedArchetype id.UUID
	selectedCharacter string
	//
	traitsImage    *ebiten.Image
	swoleImage     *ebiten.Image
	zoomsImage     *ebiten.Image
	brainsImage    *ebiten.Image
	funkImage      *ebiten.Image
	archetypeImage *ebiten.Image
}

type archetype struct {
	Archetype game.CharacterArchetype
	Image     *ebiten.Image
}

// NewCreate creates a new Create instance.
func NewCreate(connection net.Connection, msgCh chan net.Message) *Create {
	state := &Create{
		data:        NewData(),
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
	if img, err := state.data.LoadImage("images/swole.png", ctx.Game.Zoom); err != nil {
		return err
	} else {
		state.swoleImage = img
	}
	if img, err := state.data.LoadImage("images/zooms.png", ctx.Game.Zoom); err != nil {
		return err
	} else {
		state.zoomsImage = img
	}
	if img, err := state.data.LoadImage("images/brains.png", ctx.Game.Zoom); err != nil {
		return err
	} else {
		state.brainsImage = img
	}
	if img, err := state.data.LoadImage("images/funk.png", ctx.Game.Zoom); err != nil {
		return err
	} else {
		state.funkImage = img
	}
	if img, err := state.data.LoadImage("images/traits.png", ctx.Game.Zoom); err != nil {
		return err
	} else {
		state.traitsImage = img
	}
	if img, err := state.data.LoadImage("images/archetype.png", ctx.Game.Zoom); err != nil {
		return err
	} else {
		state.archetypeImage = img
	}

	state.charactersSection = widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)

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
	state.charactersSection.AddChild(state.charactersContainer)

	state.charactersControls = widget.NewContainer(
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
	state.charactersSection.AddChild(state.charactersControls)

	state.charactersJoinButton = widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.CursorHovered("interactive"),
		),
		widget.ButtonOpts.Image(ctx.UI.ButtonImage),
		widget.ButtonOpts.Text("join", ctx.UI.HeadlineFace, ctx.UI.ButtonTextColor),
		widget.ButtonOpts.TextPadding(ctx.UI.ButtonPadding),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			state.doJoin()
		}),
	)
	state.charactersControls.AddChild(state.charactersJoinButton)

	state.charactersDeleteButton = widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.CursorHovered("interactive"),
		),
		widget.ButtonOpts.Image(ctx.UI.ButtonImage),
		widget.ButtonOpts.Text("delete", ctx.UI.HeadlineFace, ctx.UI.ButtonTextColor),
		widget.ButtonOpts.TextPadding(ctx.UI.ButtonPadding),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			state.doDelete()
		}),
	)
	state.charactersControls.AddChild(state.charactersDeleteButton)

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
		widget.TextOpts.Text("Confirm deletion", ctx.UI.BodyCopyFace, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	deleteButton := widget.NewButton(
		widget.ButtonOpts.Image(ctx.UI.ButtonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CursorHovered("delete"),
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
		),
		widget.ButtonOpts.Text("delete", ctx.UI.HeadlineFace, ctx.UI.ButtonTextColor),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			state.connection.Write(net.DeleteCharacterMessage{
				Name: state.selectedCharacter,
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
			widget.TextInputOpts.Image(ctx.UI.TextInputImage),
			widget.TextInputOpts.Face(ctx.UI.BodyCopyFace),
			widget.TextInputOpts.Color(ctx.UI.TextInputColor),
			widget.TextInputOpts.Padding(ctx.UI.TextInputPadding),
			widget.TextInputOpts.CaretOpts(
				widget.CaretOpts.Size(ctx.UI.BodyCopyFace, 2),
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
			widget.ButtonOpts.Image(ctx.UI.ButtonImage),
			widget.ButtonOpts.Text("create", ctx.UI.HeadlineFace, ctx.UI.ButtonTextColor),
			widget.ButtonOpts.TextPadding(ctx.UI.ButtonPadding),
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				state.doCreate()
			}),
		)
		state.archetypesControls.AddChild(state.archetypesCreateButton)

		state.archetypesSection.AddChild(state.archetypesControls)
	}

	state.logoutButton = widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.CursorHovered("interactive"),
		),
		widget.ButtonOpts.Image(ctx.UI.ButtonImage),
		widget.ButtonOpts.Text("logout", ctx.UI.HeadlineFace, ctx.UI.ButtonTextColor),
		widget.ButtonOpts.TextPadding(ctx.UI.ButtonPadding),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			state.connection.Write(net.LogoutMessage{})
			ctx.Sm.Pop()
		}),
	)

	state.resultText = widget.NewText(
		widget.TextOpts.Text("Create a new hero or select a previous one.", ctx.UI.BodyCopyFace, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)

	state.ui.Container.AddChild(state.resultText)
	state.ui.Container.AddChild(state.charactersSection)
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
		if arch.Archetype.GetID() == archetype.GetID() {
			return true
		}
	}
	return false
}

func (state *Create) populateCharacters(ctx ifs.RunContext, characters []*game.Character) {
	state.characters = characters
	state.charactersContainer.RemoveChildren()

	var rowButtons []widget.RadioGroupElement
	for _, ch := range characters {
		func(ch *game.Character) {
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
				widget.ButtonOpts.Image(ctx.UI.ButtonImage),
				widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
					state.selectedCharacter = ch.Name
				}),
				widget.ButtonOpts.WidgetOpts(
					widget.WidgetOpts.CustomData(ch.Name), // Store name for sync reference.
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
				if arch.Archetype.GetID() == ch.ArchetypeID {
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
				widget.TextOpts.Text(ch.Name, ctx.UI.HeadlineFace, color.White),
				widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
				widget.TextOpts.WidgetOpts(
					widget.WidgetOpts.MinSize(100, 20),
					widget.WidgetOpts.LayoutData(widget.RowLayoutData{
						Position: widget.RowLayoutPositionCenter,
					}),
				),
			)

			row.AddChild(graphicContainer)
			row.AddChild(name)

			state.charactersContainer.AddChild(rowContainer)
		}(ch)
	}
	state.charactersRadioGroup = widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(rowButtons...),
	)
}

func (state *Create) acquireArchetypes(ctx ifs.RunContext, archetypes []game.Archetype) {
	for _, arch := range archetypes {
		if state.haveArchetype(arch) {
			continue
		}

		switch arch := arch.(type) {
		case game.CharacterArchetype:
			var arche archetype
			arche.Archetype = arch

			defer func() {
				state.archetypes = append(state.archetypes, arche)
			}()

			img, err := state.data.LoadImage("archetypes/"+arch.Image, ctx.Game.Zoom)
			if err != nil {
				// TODO: Show error image
				continue
			}
			state.data.archetypeImages[arch.ID] = img

			arche.Image = img
		}
	}
}

func (state *Create) refreshArchetypes(ctx ifs.RunContext) {
	state.archetypesContainer.RemoveChildren()

	// Heading
	{
		row := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			)),
		)

		el := widget.NewText(
			widget.TextOpts.Text("", ctx.UI.HeadlineFace, color.White),
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

				tool := widget.NewTextToolTip(tooltip, ctx.UI.BodyCopyFace, color.White, eimage.NewNineSliceColor(color.NRGBA{R: 50, G: 50, B: 50, A: 255}))
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
							state.refreshArchetypes(ctx)
							state.syncUI()
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
						widget.TextOpts.Text(p, ctx.UI.HeadlineFace, c),
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
				widget.ButtonOpts.Image(ctx.UI.ButtonImage),
				// add a handler that reacts to clicking the button
				widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
					state.selectedArchetype = arch.Archetype.ID
					state.archetypesCreateImage.Image = arch.Image
				}),
				widget.ButtonOpts.WidgetOpts(
					widget.WidgetOpts.CustomData(arch.Archetype.ID), // Store name for sync reference.
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
					widget.TextOpts.Text(name, ctx.UI.HeadlineFace, color.White),
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
				widget.TextOpts.Text(arch.Archetype.Title, ctx.UI.HeadlineFace, color.White),
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
				d += t.String() + "\n"
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

	state.archetypesRadioGroup = widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(rowButtons...),
	)
}

func (state *Create) syncUI() {
	// Set a default selected character if none exists.
	hasSelected := false
	for _, ch := range state.characters {
		if ch.Name == state.selectedCharacter {
			hasSelected = true
			break
		}
	}
	if !hasSelected && len(state.characters) > 0 {
		state.selectedCharacter = state.characters[0].Name
	}
	// Set a default selected archetype.
	hasSelected = false
	for _, a := range state.archetypes {
		if a.Archetype.ID == state.selectedArchetype {
			hasSelected = true
			break
		}
	}
	if !hasSelected && len(state.archetypes) > 0 {
		state.selectedArchetype = state.archetypes[0].Archetype.ID
		state.archetypesCreateImage.Image = state.archetypes[0].Image
	}

	// Synchronize the selected character. The button, which doubles as the radio, is the first child of the containing row, so we find it, then do appropriate checks until we can set the active radio to it.
	for _, w := range state.charactersContainer.Children() {
		w, ok := w.(*widget.Container)
		if !ok || len(w.Children()) == 0 {
			continue
		}
		btn, ok := w.Children()[0].(*widget.Button)
		if !ok {
			continue
		}
		n, ok := (btn.GetWidget().CustomData).(string)
		if !ok {
			continue
		}
		if n == state.selectedCharacter {
			state.charactersRadioGroup.SetActive(btn)
			break
		}
	}
	// Synchronize the selected archetype row. The button, which doubles as the radio, is the first child of the containing row, so we find it, then do appropriate checks until we can set the active radio to it.
	for _, w := range state.archetypesContainer.Children() {
		w, ok := w.(*widget.Container)
		if !ok || len(w.Children()) == 0 {
			continue
		}
		btn, ok := w.Children()[0].(*widget.Button)
		if !ok {
			continue
		}
		u, ok := (btn.GetWidget().CustomData).(id.UUID)
		if !ok {
			continue
		}
		if u == state.selectedArchetype {
			state.archetypesRadioGroup.SetActive(btn)
			break
		}
	}
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

func (state *Create) doJoin() {
	state.connection.Write(net.JoinCharacterMessage{
		Name: state.selectedCharacter,
	})
}

func (state *Create) doDelete() {
	x, y := state.deleteWindow.Contents.PreferredSize()
	r := image.Rect(0, 0, 0, 0)
	pt := input.GetWindowSize()
	r = r.Add(image.Point{pt.X / 2, pt.Y / 2})
	r = r.Sub(image.Point{x / 2, y / 2})

	state.deleteWindow.SetLocation(r)
	state.ui.AddWindow(state.deleteWindow)
}

func (state *Create) Update(ctx ifs.RunContext) error {
	select {
	case msg := <-state.messageChan:
		switch m := msg.(type) {
		case net.ArchetypesMessage:
			for _, a := range m.Archetypes {
				state.data.archetypes[a.GetID()] = a
			}
			state.acquireArchetypes(ctx, m.Archetypes)
			state.refreshArchetypes(ctx)
			state.syncUI()
		case net.CharactersMessage:
			state.populateCharacters(ctx, m.Characters)
			state.syncUI()
		case net.CreateCharacterMessage:
			if m.ResultCode == 200 {
				// Set selected character on success and clear name field.
				state.selectedCharacter = state.archetypesCreateName.InputText
				state.archetypesCreateName.InputText = ""
			}
			if m.Result != "" {
				state.resultText.Label = m.Result
			}
		case net.DeleteCharacterMessage:
			if m.Result != "" {
				state.resultText.Label = m.Result
			}
		case net.JoinCharacterMessage:
			if m.ResultCode == 200 {
				ctx.Sm.Push(NewWorlds(state.connection, state.messageChan, state.data))
			} else {
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
