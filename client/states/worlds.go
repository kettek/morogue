package states

import (
	"fmt"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/id"
	"github.com/kettek/morogue/net"
)

// Worlds is the state for joining and creating worlds. It displays a list of
// worlds acquired from net.WorldsMessage and has a section for specifying
// world options and thereafter sending a net.CreateWorldMessage.
// Popping Worlds should return to the Create state.
type Worlds struct {
	data        *Data
	connection  net.Connection
	messageChan chan net.Message
	ui          *ebitenui.UI
	//
	backButton      *widget.Button
	controlsSection *widget.Container
	//
	selectedWorld id.UUID
	worldName     string
	//
	splitSection *widget.Container
	//
	worldsSection     *widget.Container
	worldsRowsHeader  *widget.Container
	worldsRowsContent *widget.Container
	worldsControls    *widget.Container
	//
	createSection  *widget.Container
	createContent  *widget.Container
	createControls *widget.Container
	//
	worlds []game.WorldInfo
}

// NewWorlds creates a new Worlds instance.
func NewWorlds(connection net.Connection, msgCh chan net.Message, data *Data) *Worlds {
	state := &Worlds{
		data:        data,
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

func (state *Worlds) Begin(ctx ifs.RunContext) error {
	state.controlsSection = widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{255, 0, 0, 255})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  true,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
	)

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
	state.controlsSection.AddChild(state.backButton)

	//
	state.splitSection = widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
		),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x22, 0x13, 0x1a, 0xff})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
	)

	state.worldsSection = widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(800, 20),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)

	state.worldsRowsHeader = widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
	)

	state.worldsRowsContent = widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)

	state.worldsControls = widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
	)

	state.worldsSection.AddChild(state.worldsRowsHeader)
	state.worldsSection.AddChild(state.worldsRowsContent)
	state.worldsSection.AddChild(state.worldsControls)

	//
	state.createSection = widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(300, 20),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)

	state.createContent = widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)

	nameInput := widget.NewTextInput(
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
		widget.TextInputOpts.Placeholder("world name"),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			state.worldName = args.InputText
		}),
	)
	state.createContent.AddChild(nameInput)

	state.createControls = widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
	)

	createButton := widget.NewButton(
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
			state.connection.Write(net.CreateWorldMessage{
				Name: state.worldName,
			})
		}),
	)
	state.createControls.AddChild(createButton)

	state.createSection.AddChild(state.createContent)
	state.createSection.AddChild(state.createControls)

	//
	state.splitSection.AddChild(state.worldsSection)
	state.splitSection.AddChild(state.createSection)

	state.ui.Container.AddChild(state.controlsSection)
	state.ui.Container.AddChild(state.splitSection)
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

func (state *Worlds) populate(ctx ifs.RunContext, worlds []game.WorldInfo) {
	state.worlds = worlds

	state.worldsRowsContent.RemoveChildren()

	var rowButtons []widget.RadioGroupElement
	for _, w := range state.worlds {
		func(w game.WorldInfo) {
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
					//
				}),
				widget.ButtonOpts.WidgetOpts(
					widget.WidgetOpts.CustomData(w.ID), // Store name for sync reference.
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

			name := widget.NewText(
				widget.TextOpts.Text(w.Name, ctx.UI.HeadlineFace, color.White),
				widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
				widget.TextOpts.WidgetOpts(
					widget.WidgetOpts.MinSize(100, 20),
					widget.WidgetOpts.LayoutData(widget.RowLayoutData{
						Position: widget.RowLayoutPositionCenter,
					}),
				),
			)

			row.AddChild(name)

			state.worldsRowsContent.AddChild(rowContainer)
		}(w)
	}
	// TODO: radio group
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
			if m.ResultCode == 200 {
				state.populate(ctx, m.Worlds)
			}
			if m.Result == "" {
				// TODO: Show info
			}
		case net.CreateWorldMessage:
			fmt.Println("handle result of create", m)
		case net.JoinWorldMessage:
			if m.ResultCode == 200 {
				ctx.Sm.Push(NewGame(state.connection, state.messageChan, state.data))
				return nil
			}
			if m.Result == "" {
				// TODO: Show info
			}
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
