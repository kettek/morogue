package main

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/net"
)

// world represents an entire game world state that runs in its own goroutine.
// clients can join and leave the world via passed in channels.
type world struct {
	info              game.WorldInfo
	clients           []*client
	password          string
	live              bool
	data              *Data
	locations         []*location
	clientChan        chan *client
	clientRemoveChan  chan *client
	addToUniverseChan chan *client
	quitChan          chan struct{}
}

func newWorld(d *Data) *world {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	w := &world{
		info: game.WorldInfo{
			ID: id,
		},
		data:       d,
		quitChan:   make(chan struct{}),
		clientChan: make(chan *client, 2),
	}

	return w
}

func (w *world) generateLocation( /*locationInfo*/ ) {
	// TODO: Generate and add location to locations.
}

func (w *world) loop(addToUniverseChan chan *client, clientRemoveChan chan *client) {
	w.clientRemoveChan = clientRemoveChan
	w.addToUniverseChan = addToUniverseChan
	ticker := time.NewTicker(100 * time.Millisecond)

	// TODO: Ensure a starting location is being created.
	lid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	start := location{}
	start.ID = lid
	err = start.generate()
	if err != nil {
		fmt.Println("OH NO", err)
	}

	w.live = true
	for w.live {
		// Select for quit and/or client add channels.
		select {
		case <-w.quitChan:
			w.live = false
			for _, cl := range w.clients {
				w.addToUniverseChan <- cl
			}
			return
		case cl := <-w.clientChan:
			w.clients = append(w.clients, cl)
			// Send starting location to client.
			cl.conn.Write(net.LocationMessage{
				ID:         start.ID,
				Mobs:       start.Mobs,
				Cells:      start.Cells,
				Characters: start.Characters,
				Objects:    start.Objects,
			})
		default:
		}
		// Select for timer delay.
		w.update()
		<-ticker.C
	}
}

func (w *world) update() error {
	// Process clients.
	i := 0
	for _, cl := range w.clients {
		if err := w.updateClient(cl); err == nil {
			w.clients[i] = cl
			i++
		} else {
			fmt.Println(err)
		}
	}
	for j := i; j < len(w.clients); j++ {
		w.clients[j] = nil
	}
	w.clients = w.clients[:i]

	// Process locations.
	i = 0
	for _, l := range w.locations {
		if err := w.processLocation(l); err == nil {
			w.locations[i] = l
			i++
		} else {
			fmt.Println(err)
		}
	}
	for j := i; j < len(w.locations); j++ {
		w.locations[j] = nil
	}
	w.locations = w.locations[:i]

	return nil
}

func (w *world) updateClient(cl *client) error {
	select {
	case msg := <-cl.msgChan:
		switch m := msg.(type) {
		case net.TileMessage:
			if t, err := w.data.Tile(m.ID); err != nil {
				cl.conn.Write(net.TileMessage{
					ResultCode: 404,
					Result:     err.Error(),
					ID:         m.ID,
				})
			} else {
				cl.conn.Write(net.TileMessage{
					ResultCode: 200,
					ID:         m.ID,
					Tile:       t,
				})
			}
		default:
			fmt.Println("TODO: handle client message", m)
		}
		// TODO: If the location the client is traveling to is not done, send progress reports to client.
	case err := <-cl.closedChan:
		w.clientRemoveChan <- cl
		fmt.Println("client yeeted from world context", err)
		return err
	default:
	}
	return nil
}

func (w *world) processLocation(l *location) error {
	// TODO: Probably add remove/clear timer for particular location types?
	if !l.active {
		return nil
	}
	return l.process()
}
