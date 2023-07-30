package net

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
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
		c.SetReadLimit(1024 * 1024)
		if err != nil {
			ch <- err
			return
		}

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
			if err := wsjson.Read(context.TODO(), conn.c, &w); err != nil {
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
	p, err := json.Marshal(m)
	if err != nil {
		return err
	}

	w := Wrapper{
		Type: m.Type(),
		Data: p,
	}

	return wsjson.Write(context.TODO(), conn.c, w)
}
