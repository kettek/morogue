package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kettek/morogue/net"
	"nhooyr.io/websocket"
)

type socketServer struct {
	// logf controls where logs are sent.
	logf func(f string, v ...interface{})
	//
	serveMux http.ServeMux
	//
	newClientChan chan client
	checkChan     chan struct{}
}

func newSocketServer(newClientChan chan client, checkChan chan struct{}) *socketServer {
	p := &socketServer{
		logf:          log.Printf,
		newClientChan: newClientChan,
		checkChan:     checkChan,
	}

	p.serveMux.Handle("/", http.FileServer(http.Dir("./static")))
	p.serveMux.HandleFunc("/sockit", p.handleSockit)

	return p
}

func (s *socketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *socketServer) handleSockit(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"morogue"},
	})
	if err != nil {
		s.logf("%v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	if c.Subprotocol() != "morogue" {
		c.Close(websocket.StatusPolicyViolation, "client must speak the morogue subprotocol")
		return
	}

	conn := net.NewConnection(c)
	client := client{
		conn:       conn,
		msgChan:    make(chan net.Message, 10),
		closedChan: make(chan error, 2),
	}

	s.newClientChan <- client

	for {
		var w net.Wrapper

		_, b, err := c.Read(r.Context())
		if err == nil {
			err = json.Unmarshal(b, &w)
			if err != nil {
				s.logf("failed to unmarshal with %v: %v", r.RemoteAddr, err)
			} else if m := w.Message(); m != nil {
				client.msgChan <- m
				s.checkChan <- struct{}{}
				conn.Write(&net.PingMessage{})
			}
		}

		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			client.closedChan <- err
			return
		}
		if err != nil {
			client.closedChan <- err
			s.checkChan <- struct{}{}
			return
		}
	}
}
