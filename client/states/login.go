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

type Login struct {
	connection  net.Connection
	messageChan chan net.Message
	ui          *ebitenui.UI
	//
	usernameInput               *widget.TextInput
	passwordInput               *widget.TextInput
	confirmInput                *widget.TextInput
	resultText                  *widget.Text
	loginButton, registerButton *widget.Button
}

func NewLogin(connection net.Connection, msgCh chan net.Message) *Login {
	state := &Login{
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

func (state *Login) Begin(ctx ifs.RunContext) error {
	state.connection.Write(&net.PingMessage{})

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
		widget.TextInputOpts.Placeholder("username"),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			state.checkInputs()
		}),
	)

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
		widget.TextInputOpts.Placeholder("password"),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			state.checkInputs()
		}),
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
		widget.TextInputOpts.Placeholder("confirm"),
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
		widget.ButtonOpts.Text("login", ctx.UI.HeadlineFace, ctx.UI.ButtonTextColor),
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
		widget.ButtonOpts.Text("register", ctx.UI.HeadlineFace, ctx.UI.ButtonTextColor),
		widget.ButtonOpts.TextPadding(ctx.UI.ButtonPadding),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			state.doRegister()
		}),
	)

	state.resultText = widget.NewText(
		widget.TextOpts.Text("login. you will be prompted to register if username does not exist.", ctx.UI.BodyCopyFace, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)

	state.showLogin()

	return nil
}

func (state *Login) checkInputs() {
	if state.usernameInput.InputText == "" {
		state.resultText.Label = "username must not be empty"
		return
	}
	if state.passwordInput.InputText == "" {
		state.resultText.Label = "password must not be empty"
		return
	}
	hasConfirm := false
	for _, ch := range state.ui.Container.Children() {
		if ch == state.confirmInput {
			hasConfirm = true
			break
		}
	}
	if hasConfirm && state.passwordInput.InputText != state.confirmInput.InputText {
		state.resultText.Label = "passwords must match"
		return
	}
	state.resultText.Label = ""
}

func (state *Login) doLogin() {
	if state.usernameInput.InputText == "" || state.passwordInput.InputText == "" {
		return
	}
	state.resultText.Label = "logging in..."
	state.connection.Write(net.LoginMessage{
		User:     state.usernameInput.InputText,
		Password: state.passwordInput.InputText,
	})
}

func (state *Login) doRegister() {
	if state.usernameInput.InputText == "" || state.passwordInput.InputText == "" || state.passwordInput.InputText != state.confirmInput.InputText {
		return
	}
	state.resultText.Label = "registering..."
	state.connection.Write(net.RegisterMessage{
		User:     state.usernameInput.InputText,
		Password: state.passwordInput.InputText,
	})
}

func (state *Login) showLogin() {
	state.ui.Container.RemoveChild(state.resultText)
	state.ui.Container.RemoveChild(state.confirmInput)
	state.ui.Container.RemoveChild(state.registerButton)
	state.ui.Container.AddChild(state.usernameInput)
	state.ui.Container.AddChild(state.passwordInput)
	state.ui.Container.AddChild(state.loginButton)

	state.ui.Container.AddChild(state.resultText)
}

func (state *Login) showRegister() {
	state.ui.Container.RemoveChild(state.resultText)
	state.ui.Container.RemoveChild(state.loginButton)
	state.ui.Container.AddChild(state.confirmInput)
	state.ui.Container.AddChild(state.registerButton)

	state.ui.Container.AddChild(state.resultText)
}

func (state *Login) Return(interface{}) error {
	state.resultText.Label = "...and so you return."
	state.ui.Container.RemoveChildren()
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
				state.resultText.Label = "logged in!"
				ctx.Sm.Push(NewCreate(state.connection, state.messageChan))
				return nil
			} else if m.ResultCode == 404 {
				state.showRegister()
				state.resultText.Label = "Confirm your password to register."
			} else {
				state.resultText.Label = m.Result
			}
		case net.RegisterMessage:
			if m.ResultCode == 200 {
				state.resultText.Label = "logged in!"
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
