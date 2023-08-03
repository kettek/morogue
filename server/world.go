package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/id"
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
	wids              id.WIDGenerator
	locations         []*location
	clientChan        chan *client
	clientRemoveChan  chan *client
	addToUniverseChan chan *client
	quitChan          chan struct{}
}

func newWorld(d *Data) *world {
	wid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	w := &world{
		info: game.WorldInfo{
			ID: id.UUID(wid),
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

// assignWIDs assigns a WID to an object and all of its children. This is done when any object is added to the world.
func (w *world) assignWIDs(o game.Object) {
	o.SetWID(w.wids.Next())
	switch o := o.(type) {
	case *game.Character:
		for _, o2 := range o.Inventory {
			w.assignWIDs(o2)
		}
		/*case *game.Container:
		for _, i := range o.GetItems() {
			w.assignWIDs(i)
		}*/
	}
}

func (w *world) loop(addToUniverseChan chan *client, clientRemoveChan chan *client) {
	w.clientRemoveChan = clientRemoveChan
	w.addToUniverseChan = addToUniverseChan
	ticker := time.NewTicker(50 * time.Millisecond)

	// TODO: Ensure a starting location is being created.
	lid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	start := &location{}
	start.ID = id.UUID(lid)
	start.active = true
	err = start.generate()
	if err != nil {
		fmt.Println("OH NO", err)
	}
	w.locations = append(w.locations, start)

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
			var char *game.Character
			for _, ch := range cl.account.Characters {
				if ch.Name == cl.character {
					char = ch
					break
				}
			}
			// Boot back to universe if char is nil... TODO: Add error message.
			if char == nil {
				w.addToUniverseChan <- cl
				return
			}

			// Add client as character to world. Boot back to universe if placement failed. TODO: Add error message.
			if err := start.addCharacter(char); err != nil {
				w.addToUniverseChan <- cl
				return
			}
			// Assign WIDs to character and their inventory.
			w.assignWIDs(char)
			cl.currentLocation = start
			cl.currentCharacter = char

			// Send starting location to client.
			cl.conn.Write(net.LocationMessage{
				ID:      start.ID,
				Cells:   start.Cells,
				Objects: start.Objects,
			})

			// Send client their character owner message
			cl.conn.Write(net.OwnerMessage{
				WID:       char.WID,
				Inventory: char.Inventory,
				Skills:    char.Skills,
			})

			// Send create to clients in location.
			if evt, err := game.WrapEvent(game.EventAdd{
				Object: cl.currentCharacter,
			}); err == nil {
				cls := w.clientsInLocation(start)
				for _, cl := range cls {
					cl.conn.Write(net.EventMessage{
						Event: evt,
					})
				}
			}
			// Now add the client to the list.
			w.clients = append(w.clients, cl)
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
			// This shouldn't ever be nil, but let's be safe.
			if cl.currentCharacter != nil {
				for _, l := range w.locations {
					if err := l.removeCharacter(cl.currentCharacter.WID); err == nil {
						// Send remove to clients in location, excluding the removed client.
						if evt, err := game.WrapEvent(game.EventRemove{
							WID: cl.currentCharacter.WID,
						}); err == nil {
							cls := w.clientsInLocation(l)
							for _, cl2 := range cls {
								if cl == cl2 {
									continue
								}
								cl2.conn.Write(net.EventMessage{
									Event: evt,
								})
							}
						}
						break
					}
				}
			}
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
			if err != errRemoveLocationFromWorld {
				fmt.Println(err)
			}
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
		case net.ArchetypesMessage:
			var archetypes []game.Archetype
			for _, uuid := range m.IDs {
				if a := w.data.Archetype(uuid); a == nil {
					cl.conn.Write(net.ArchetypeMessage{
						ResultCode: 404,
						ID:         uuid,
					})
				} else {
					archetypes = append(archetypes, a)
				}
			}
			if len(archetypes) > 0 {
				cl.conn.Write(net.ArchetypesMessage{
					Archetypes: archetypes,
				})
			}
		default:
			// For all other messages, pass off handling to client's current location.
			if cl.currentLocation != nil {
				cl.currentLocation.handleClientMessage(cl, msg)
			} else {
				fmt.Println("message sent to the ether", msg)
			}
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

func (w *world) clientsInLocation(l *location) []*client {
	var locationClients []*client
	for _, cl := range w.clients {
		if l.Character(cl.currentCharacter.WID) != nil {
			locationClients = append(locationClients, cl)
		}
	}
	return locationClients
}

func (w *world) processLocation(l *location) error {
	// TODO: Probably add remove/clear timer for particular location types?
	if !l.active {
		if l.removable && time.Since(l.emptySince) > 5*time.Minute {
			return errRemoveLocationFromWorld
		}
		return nil
	}

	locationClients := w.clientsInLocation(l)

	// Process events for all the clients in this location.
	events := l.process()

	// Convert & send private client events.
	for _, cl := range locationClients {
		if cl.currentCharacter.Events != nil {
			var eventsMessage net.EventsMessage
			for _, event := range cl.currentCharacter.Events {
				if evt, err := game.WrapEvent(event); err == nil {
					eventsMessage.Events = append(eventsMessage.Events, evt)
				}
			}
			cl.conn.Write(eventsMessage)
		}
		cl.currentCharacter.Events = nil
	}

	// Convert events to be sent to clients.
	var eventsMessage net.EventsMessage
	for _, event := range events {
		if evt, err := game.WrapEvent(event); err == nil {
			eventsMessage.Events = append(eventsMessage.Events, evt)
		}
	}
	// Send events to clients.
	if len(eventsMessage.Events) > 0 {
		for _, cl := range locationClients {
			cl.conn.Write(eventsMessage)
		}
	}

	return nil
}

var (
	errRemoveLocationFromWorld = errors.New("this is also not an error lol")
)
