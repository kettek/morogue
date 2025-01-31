package states

import (
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/config"
	"github.com/kettek/morogue/net"

	// Just for "L"
	"github.com/kettek/morogue/locale"
)

// Server is the first state when connecting to a server. It provides
// the ability to login or register an account with the server.
type Server struct {
	//connection  net.Connection
	//messageChan chan net.Message
	ui *ebitenui.UI
	//
	inputs   *widget.Container
	controls *widget.Container
	//
	serverInput *widget.TextInput
	resultText  *widget.Text
	joinButton  *widget.Button
	lc          locale.Localizer
	// Pointers to Connect's fields.
	connMode *string
	connChan *chan error
	conn     *net.Connection
}

// NewServer creates a new Server instance.
func NewServer(mode *string, connection *net.Connection, connChan *chan error) *Server {
	state := &Server{
		//connection:  connection,
		//messageChan: msgCh,
		ui: &ebitenui.UI{
			Container: widget.NewContainer(
				widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x22, 0x13, 0x1a, 0xff})),
				widget.ContainerOpts.Layout(widget.NewAnchorLayout(
					widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(20)),
				)),
			),
		},
		lc:       locale.Get(locale.Locale()),
		connMode: mode,
		conn:     connection,
		connChan: connChan,
	}
	return state
}

func (state *Server) Begin(ctx ifs.RunContext) error {
	innerContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(20),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20))),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	state.serverInput = widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
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
		widget.TextInputOpts.Placeholder(state.lc.T("server")),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			state.checkInputs()
			ctx.Cfg.LastServer = args.InputText
			config.Save()
		}),
	)
	state.serverInput.SetText(ctx.Cfg.LastServer)

	state.inputs = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionStart,
				Stretch:  true,
			}),
		),
	)

	state.controls = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
			widget.WidgetOpts.MinSize(400, 20),
		),
	)

	state.joinButton = widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.CursorHovered("interactive"),
		),
		widget.ButtonOpts.Image(ctx.UI.ButtonImage),
		widget.ButtonOpts.Text(state.lc.T("login"), ctx.UI.HeadlineFace, ctx.UI.ButtonTextColor),
		widget.ButtonOpts.TextPadding(ctx.UI.ButtonPadding),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			state.doJoin(ctx)
		}),
	)

	state.resultText = widget.NewText(
		widget.TextOpts.Text(state.lc.T("join a server."), ctx.UI.BodyCopyFace, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)

	innerContainer.AddChild(state.inputs)
	innerContainer.AddChild(state.controls)
	innerContainer.AddChild(state.resultText)

	state.ui.Container.AddChild(innerContainer)

	state.showJoin()

	return nil
}

func (state *Server) checkInputs() {
	if state.serverInput.GetText() == "" {
		state.resultText.Label = state.lc.T("server must not be empty")
		return
	}
	state.resultText.Label = ""
}

func (state *Server) doJoin(ctx ifs.RunContext) {
	if state.serverInput.GetText() == "" {
		return
	}
	state.resultText.Label = state.lc.T("joining...")
	*state.connMode = modeConnecting
	*state.connChan = state.conn.Connect(state.serverInput.GetText())
	ctx.Sm.Pop()
	/*state.connection.Write(net.ServerMessage{
		User:     state.usernameInput.GetText(),
		Password: state.passwordInput.GetText(),
	})*/
}

func (state *Server) showJoin() {
	state.inputs.RemoveChildren()
	state.inputs.AddChild(state.serverInput)

	state.controls.RemoveChildren()
	state.controls.AddChild(state.joinButton)
}

func (state *Server) Return(interface{}) error {
	state.resultText.Label = state.lc.T("ok")

	return nil
}

func (state *Server) Leave() error {
	return nil
}

func (state *Server) End() (interface{}, error) {
	return nil, nil
}

func (state *Server) Update(ctx ifs.RunContext) error {
	/*select {
	case msg := <-state.messageChan:
		switch m := msg.(type) {
		case net.ServerMessage:
			if m.ResultCode == 200 {
				state.resultText.Label = state.lc.T("logged in!")
				ctx.Sm.Push(NewCreate(state.connection, state.messageChan))
				return nil
			} else if m.ResultCode == 404 {
				state.showRegister()
				state.resultText.Label = state.lc.T("Confirm your password to register.")
			} else {
				state.resultText.Label = m.Result
			}
		case net.RegisterMessage:
			if m.ResultCode == 200 {
				state.resultText.Label = state.lc.T("logged in!")
				ctx.Sm.Push(NewCreate(state.connection, state.messageChan))
				return nil
			} else {
				state.showServer()
				state.resultText.Label = m.Result
			}
		}
		fmt.Println("got eem", msg)
	default:
	}*/

	state.ui.Update()

	return nil
}

func (state *Server) Draw(ctx ifs.DrawContext) {
	state.ui.Draw(ctx.Screen)
}
