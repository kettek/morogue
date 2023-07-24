package states

import (
	"fmt"
	"image/color"

	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/net"
)

const (
	modeConnecting = "connecting"
	modeFailed     = "failed"
	modeSuccess    = "success"
)

type Connect struct {
	connection     net.Connection
	connectionChan chan error
	loopChan       chan error
	messageChan    chan net.Message
	mode           string
	result         string
}

func (state *Connect) Begin() error {
	state.mode = modeConnecting
	state.connectionChan = state.connection.Connect()
	return nil
}

func (state *Connect) Return(interface{}) error {
	return nil
}

func (state *Connect) Leave() error {
	return nil
}

func (state *Connect) End() (interface{}, error) {
	return nil, nil
}

func (state *Connect) Update(ctx ifs.RunContext) error {
	select {
	case err := <-state.connectionChan:
		if err != nil {
			state.mode = modeFailed
			state.result = err.Error()
			return nil
		}
		state.mode = modeSuccess
		state.messageChan, state.loopChan = state.connection.Loop()
		state.connection.Write(&net.PingMessage{})
	case err := <-state.loopChan:
		state.mode = modeFailed
		state.result = err.Error()
	default:
		//
	}

	if state.mode == modeSuccess {
		select {
		case msg := <-state.messageChan:
			fmt.Println("got eem", msg)
		default:
		}
	}

	return nil
}

func (state *Connect) Draw(ctx ifs.DrawContext) {
	// background color
	ctx.Screen.Fill(color.NRGBA{0, 0, 0, 255})

	// get screen center position and text content
	bounds := ctx.Screen.Bounds() // assumes origin (0, 0)
	x, y := bounds.Dx()/2, bounds.Dy()/2

	// draw the text
	ctx.Txt.Draw(ctx.Screen, state.mode, x, y)
	y += int(ctx.Txt.Utils().GetLineHeight())
	var last int
	for i := 0; i <= len(state.result); {
		if i >= len(state.result) {
			ctx.Txt.Draw(ctx.Screen, state.result[last:i], x, y)
			break
		} else if state.result[i] == '\n' {
			ctx.Txt.Draw(ctx.Screen, state.result[last:i], x, y)
			y += int(ctx.Txt.Utils().GetLineHeight())
		}
		i++
	}
}
