package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kettek/morogue/net"
	"nhooyr.io/websocket"
)

// protoServer is the WebSocket echo server implementation.
// It ensures the client speaks the echo subprotocol and
// only allows one message every 100ms with a 10 message burst.
type protoServer struct {
	// logf controls where logs are sent.
	logf func(f string, v ...interface{})
	//
	serveMux http.ServeMux
}

func newProtoServer() *protoServer {
	p := &protoServer{
		logf: log.Printf,
	}

	p.serveMux.Handle("/", http.FileServer(http.Dir("./static")))
	p.serveMux.HandleFunc("/sockit", p.handleSockit)

	return p
}

func (s *protoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *protoServer) handleSockit(w http.ResponseWriter, r *http.Request) {
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

	cl := net.NewConnection(c)

	for {
		var w net.Wrapper

		_, b, err := c.Read(r.Context())
		if err == nil {
			err = json.Unmarshal(b, &w)
			if err != nil {
				s.logf("failed to unmarshal with %v: %v", r.RemoteAddr, err)
				return
			}
			if m := w.Message(); m != nil {
				// TODO: Send m to client instance for handling.
				fmt.Println("client msg", m, m.Type())
				cl.Write(&net.PingMessage{})
			}
		}

		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			s.logf("failed to read with %v: %v", r.RemoteAddr, err)
			// TODO: Handle err reason, if it exists.
			return
		}
	}
}
