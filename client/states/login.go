package states

import (
	"fmt"

	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/net"
)

type Login struct {
	connection  net.Connection
	messageChan chan net.Message
}

func NewLogin(connection net.Connection, msgCh chan net.Message) *Login {
	state := &Login{
		connection:  connection,
		messageChan: msgCh,
	}
	return state
}

func (state *Login) Begin() error {
	state.connection.Write(&net.PingMessage{})
	return nil
}

func (state *Login) Return(interface{}) error {
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
		fmt.Println("got eem", msg)
	default:
	}

	return nil
}

func (state *Login) Draw(ctx ifs.DrawContext) {

}
