package server

import (
	"log"
	"net/http"

	"github.com/kettek/morogue/net"
	"github.com/vmihailenco/msgpack/v5"
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

// NewSocketServer returns a new socketServer, providing functionality for serving web assets, archetypes, images, and the game websocket server.
func NewSocketServer(newClientChan chan client, checkChan chan struct{}) http.Handler {
	p := &socketServer{
		logf:          log.Printf,
		newClientChan: newClientChan,
		checkChan:     checkChan,
	}

	p.serveMux.Handle("/", http.FileServer(http.Dir("./static")))
	p.serveMux.Handle("/archetypes/", http.StripPrefix("/archetypes/", http.FileServer(http.Dir("./archetypes"))))
	p.serveMux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	p.serveMux.HandleFunc("/sockit", p.handleSockit)

	return p
}

// ServeHTTP does what you think it dfoes.
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
			err = msgpack.Unmarshal(b, &w)
			if err != nil {
				s.logf("failed to unmarshal with %v: %v", r.RemoteAddr, err)
			} else if m := w.Message(); m != nil {
				client.msgChan <- m
				s.checkChan <- struct{}{}
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
