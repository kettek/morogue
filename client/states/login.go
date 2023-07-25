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

	// load images for button states: idle, hover, and pressed
	buttonImages, _ := buttonImages()

	// load button text font
	face, _ := opentype.NewFace(ctx.Txt.Renderer.GetFont(), &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	state.usernameInput = widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			//Set the layout information to center the textbox in the parent
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
			widget.WidgetOpts.CursorHovered("text"),
		),

		//Set the Idle and Disabled background image for the text input
		//If the NineSlice image has a minimum size, the widget will use that or
		// widget.WidgetOpts.MinSize; whichever is greater
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
			Disabled: image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
		}),

		//Set the font face and size for the widget
		widget.TextInputOpts.Face(face),

		//Set the colors for the text and caret
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.NRGBA{254, 255, 255, 255},
			Disabled:      color.NRGBA{R: 200, G: 200, B: 200, A: 255},
			Caret:         color.NRGBA{254, 255, 255, 255},
			DisabledCaret: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		}),

		//Set how much padding there is between the edge of the input and the text
		widget.TextInputOpts.Padding(widget.NewInsetsSimple(5)),

		//Set the font and width of the caret
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(face, 2),
		),

		//This text is displayed if the input is empty
		widget.TextInputOpts.Placeholder("username"),

		//This is called when the user hits the "Enter" key.
		//There are other options that can configure this behavior
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Submitted: ", args.InputText)
		}),

		//This is called whenver there is a change to the text
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
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
			Disabled: image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
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
		//This parameter indicates that the inputted text should be hidden
		widget.TextInputOpts.Secure(true),

		widget.TextInputOpts.Placeholder("password"),
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Submitted: ", args.InputText)
		}),
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
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
			Disabled: image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
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
		//This parameter indicates that the inputted text should be hidden
		widget.TextInputOpts.Secure(true),

		widget.TextInputOpts.Placeholder("confirm"),
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Submitted: ", args.InputText)
		}),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			state.checkInputs()
		}),
	)

	state.loginButton = widget.NewButton(
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
		widget.ButtonOpts.Text("login", face, &widget.ButtonTextColor{
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
			state.doLogin()
		}),
	)

	state.registerButton = widget.NewButton(
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
		widget.ButtonOpts.Text("register", face, &widget.ButtonTextColor{
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
			state.doRegister()
		}),
	)

	state.resultText = widget.NewText(
		widget.TextOpts.Text("login. you will be prompted to register if username does not exist.", face, color.White),
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
