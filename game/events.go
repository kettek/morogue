package game

import (
	"encoding/json"

	"github.com/kettek/morogue/id"
)

// Event is the result of something happening on the server that is to be sent to the client. This includes sounds, position information, damage dealt, and more. Many events are as the result of client-sent Desires.
type Event interface {
	Type() string
}

// EventWrapper is for sending desires from the client to the server.
type EventWrapper struct {
	Type string          `json:"t"`
	Data json.RawMessage `json:"d"`
}

// WrapEvent wraps up an event to be sent over the wire.
func WrapEvent(e Event) (EventWrapper, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return EventWrapper{}, err
	}

	return EventWrapper{
		Type: e.Type(),
		Data: b,
	}, nil
}

// Event returns the event stored in the wrapper.
func (w *EventWrapper) Event() Event {
	switch w.Type {
	case (EventPosition{}).Type():
		var d EventPosition
		json.Unmarshal(w.Data, &d)
		return d
	case (EventSound{}).Type():
		var d EventSound
		json.Unmarshal(w.Data, &d)
		return d
	case (EventRemove{}).Type():
		var d EventRemove
		json.Unmarshal(w.Data, &d)
		return d
	case (EventAdd{}).Type():
		var d EventAdd
		json.Unmarshal(w.Data, &d)
		return d
	}
	return nil
}

// EventPosition represents a position update of something in a world location.
type EventPosition struct {
	WID  id.WID
	X, Y int
}

func (e EventPosition) Type() string {
	return "position"
}

// EventSound represents a sound emitted from a location. FromX and FromY are used to modify the visual offset of the sound. This makes it so when you bump into a wall or hit an enemy, the sound effect appears between the two points.
type EventSound struct {
	WID     id.WID `json:"wid,omitempty"`
	X       int    `json:"x,omitempty"`
	Y       int    `json:"y,omitempty"`
	FromX   int    `json:"fx,omitempty"`
	FromY   int    `json:"fy,omitempty"`
	Message string `json:"m,omitempty"`
}

func (e EventSound) Type() string {
	return "sound"
}

type EventRemove struct {
	WID id.WID `json:"wid,omitempty"`
}

func (e EventRemove) Type() string {
	return "remove"
}

type EventAdd struct {
	Object ObjectWrapper `json:"o,omitempty"`
}

func (e EventAdd) Type() string {
	return "add"
}
