package states

import (
	"fmt"
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

// Login is the first state when connecting to a server. It provides
// the ability to login or register an account with the server.
type Login struct {
	connection  net.Connection
	messageChan chan net.Message
	ui          *ebitenui.UI
	//
	inputs   *widget.Container
	controls *widget.Container
	//
	usernameInput                           *widget.TextInput
	passwordInput                           *widget.TextInput
	confirmInput                            *widget.TextInput
	resultText                              *widget.Text
	loginButton, registerButton, backButton *widget.Button
	lc                                      locale.Localizer
}

// NewLogin creates a new Login instance.
func NewLogin(connection net.Connection, msgCh chan net.Message) *Login {
	state := &Login{
		connection:  connection,
		messageChan: msgCh,
		ui: &ebitenui.UI{
			Container: widget.NewContainer(
				widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x22, 0x13, 0x1a, 0xff})),
				widget.ContainerOpts.Layout(widget.NewAnchorLayout(
					widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(20)),
				)),
			),
		},
		lc: locale.Get("en-us"),
	}
	return state
}

func (state *Login) Begin(ctx ifs.RunContext) error {
	state.connection.Write(net.PingMessage{})

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

	state.usernameInput = widget.NewTextInput(
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
		widget.TextInputOpts.Placeholder(state.lc.T("username")),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			state.checkInputs()
			ctx.Cfg.LastUsername = args.InputText
			config.Save()
		}),
	)
	state.usernameInput.SetText(ctx.Cfg.LastUsername)

	state.passwordInput = widget.NewTextInput(
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
		widget.TextInputOpts.Secure(true),
		widget.TextInputOpts.Placeholder(state.lc.T("password")),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			state.checkInputs()
			ctx.Cfg.LastPassword = args.InputText
			config.Save()
		}),
	)
	state.passwordInput.SetText(ctx.Cfg.LastPassword) // FIXME: Only store the hash and send it to the server.

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

	state.confirmInput = widget.NewTextInput(
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
		widget.TextInputOpts.Secure(true),
		widget.TextInputOpts.Placeholder(state.lc.T("confirm")),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			state.checkInputs()
		}),
	)

	state.loginButton = widget.NewButton(
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
			state.doLogin()
		}),
	)

	state.registerButton = widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.CursorHovered("interactive"),
		),
		widget.ButtonOpts.Image(ctx.UI.ButtonImage),
		widget.ButtonOpts.Text(state.lc.T("register"), ctx.UI.HeadlineFace, ctx.UI.ButtonTextColor),
		widget.ButtonOpts.TextPadding(ctx.UI.ButtonPadding),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			state.doRegister()
		}),
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
		widget.ButtonOpts.Text(state.lc.T("back"), ctx.UI.HeadlineFace, ctx.UI.ButtonTextColor),
		widget.ButtonOpts.TextPadding(ctx.UI.ButtonPadding),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			state.showLogin()
		}),
	)

	state.resultText = widget.NewText(
		widget.TextOpts.Text(state.lc.T("login. you will be prompted to register if username does not exist."), ctx.UI.BodyCopyFace, color.White),
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

	state.showLogin()

	return nil
}

func (state *Login) checkInputs() {
	if state.usernameInput.GetText() == "" {
		state.resultText.Label = state.lc.T("username must not be empty")
		return
	}
	if state.passwordInput.GetText() == "" {
		state.resultText.Label = state.lc.T("password must not be empty")
		return
	}
	hasConfirm := false
	for _, ch := range state.ui.Container.Children() {
		if ch == state.confirmInput {
			hasConfirm = true
			break
		}
	}
	if hasConfirm && state.passwordInput.GetText() != state.confirmInput.GetText() {
		state.resultText.Label = state.lc.T("passwords must match")
		return
	}
	state.resultText.Label = ""
}

func (state *Login) doLogin() {
	if state.usernameInput.GetText() == "" || state.passwordInput.GetText() == "" {
		return
	}
	state.resultText.Label = state.lc.T("logging in...")
	state.connection.Write(net.LoginMessage{
		User:     state.usernameInput.GetText(),
		Password: state.passwordInput.GetText(),
	})
}

func (state *Login) doRegister() {
	if state.usernameInput.GetText() == "" || state.passwordInput.GetText() == "" || state.passwordInput.GetText() != state.confirmInput.GetText() {
		return
	}
	state.resultText.Label = state.lc.T("registering...")
	state.connection.Write(net.RegisterMessage{
		User:     state.usernameInput.GetText(),
		Password: state.passwordInput.GetText(),
	})
}

func (state *Login) showLogin() {
	state.inputs.RemoveChildren()
	state.inputs.AddChild(state.usernameInput)
	state.inputs.AddChild(state.passwordInput)

	state.controls.RemoveChildren()
	state.controls.AddChild(state.loginButton)
}

func (state *Login) showRegister() {
	state.inputs.RemoveChildren()
	state.inputs.AddChild(state.usernameInput)
	state.inputs.AddChild(state.passwordInput)
	state.inputs.AddChild(state.confirmInput)

	state.controls.RemoveChildren()
	state.controls.AddChild(state.backButton)
	state.controls.AddChild(state.registerButton)

	state.ui.Container.AddChild(state.resultText)
}

func (state *Login) Return(interface{}) error {
	state.resultText.Label = state.lc.T("...and so you return.")
	state.showLogin()

	return nil
}

func (state *Login) Leave() error {
	return nil
}

func (state *Login) End() (interface{}, error) {
	return nil, nil
}

func (state *Login) Update(ctx ifs.RunContext) error {
	select {
	case msg := <-state.messageChan:
		switch m := msg.(type) {
		case net.LoginMessage:
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
				state.showLogin()
				state.resultText.Label = m.Result
			}
		}
		fmt.Println("got eem", msg)
	default:
	}

	state.ui.Update()

	return nil
}

func (state *Login) Draw(ctx ifs.DrawContext) {
	state.ui.Draw(ctx.Screen)
}
