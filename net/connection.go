package net

import (
	"context"
	"fmt"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"nhooyr.io/websocket"
)

type Connection struct {
	c *websocket.Conn
}

func NewConnection(c *websocket.Conn) *Connection {
	return &Connection{c}
}

func (conn *Connection) Connect() chan error {
	ch := make(chan error)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		c, _, err := websocket.Dial(ctx, "ws://localhost:8080/sockit", &websocket.DialOptions{
			Subprotocols: []string{"morogue"},
		})
		if err != nil {
			ch <- err
			return
		}
		c.SetReadLimit(1024 * 1024)

		conn.c = c

		fmt.Println("we in like flynn")

		ch <- nil
	}()

	return ch
}

func (conn *Connection) Close() {
	if conn.c == nil {
		return
	}
	conn.c.Close(websocket.StatusNormalClosure, "")
	conn.c = nil
}

func (conn *Connection) Loop() (chan Message, chan error) {
	mch := make(chan Message, 10)
	ech := make(chan error)
	go func() {
		for {
			var w Wrapper
			_, b, err := conn.c.Read(context.TODO())
			if err != nil {
				ech <- err
				return
			}
			if err := msgpack.Unmarshal(b, &w); err != nil {
				ech <- err
				return
			} else {
				if m := w.Message(); m != nil {
					mch <- m
				}
			}
		}
	}()
	return mch, ech
}

func (conn *Connection) Write(m Message) error {
	p, err := msgpack.Marshal(m)
	if err != nil {
		return err
	}

	w := Wrapper{
		Type: m.Type(),
		Data: p,
	}

	p, err = msgpack.Marshal(w)
	if err != nil {
		return err
	}

	return conn.c.Write(context.TODO(), websocket.MessageBinary, p)
}
