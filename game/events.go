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
	case (EventPickup{}).Type():
		var d EventPickup
		json.Unmarshal(w.Data, &d)
		return d
	case (EventDrop{}).Type():
		var d EventDrop
		json.Unmarshal(w.Data, &d)
		return d
	case (EventApply{}).Type():
		var d EventApply
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

// Type returns "position"
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

// Type returns "sound"
func (e EventSound) Type() string {
	return "sound"
}

// EventRemove removes an object with the given WID from the current location.
type EventRemove struct {
	WID id.WID `json:"wid,omitempty"`
}

// Type returns "remove"
func (e EventRemove) Type() string {
	return "remove"
}

// EventAdd adds the provided object.
type EventAdd struct {
	Object Object `json:"o,omitempty"`
}

// Type returns "add"
func (e EventAdd) Type() string {
	return "add"
}

// eventAdd is used internally as the real structure for JSON marshal/unmarshal.
// This is done so as to have the resulting json from EventAdd contain proper
// fields rather than a direct ObjectWrapper object. That is to say:
// event: {o: {t: "type", d: ...}} rather than {t: "type", d: ...}
// This is so if eventAdd ever needs more fields we can add them and also have
// the expected event->fields structure remain constant amonst all events.
type eventAdd struct {
	Object ObjectWrapper `json:"o,omitempty"`
}

// MarshalJSON marshals EventAdd into eventAdd.
func (e EventAdd) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(e.Object)
	if err != nil {
		return nil, err
	}

	e2 := eventAdd{
		Object: ObjectWrapper{
			Type: e.Object.Type(),
			Data: b,
		},
	}

	return json.Marshal(e2)
}

// UnmarshalJSON unmarshals EventAdd from eventAdd.
func (e *EventAdd) UnmarshalJSON(b []byte) error {
	var e2 eventAdd

	if err := json.Unmarshal(b, &e2); err != nil {
		return err
	}
	o, err := e2.Object.Object()
	if err != nil {
		return err
	}
	e.Object = o

	return nil
}

// EventApply notifies the client that the given item was applied.
type EventApply struct {
	WID     id.WID
	Applied bool `json:"a,omitempty"`
}

// Type returns "apply".
func (e EventApply) Type() string {
	return "apply"
}

// EventPickup notifies the client that the given item was picked up.
type EventPickup struct {
	WID     id.WID
	IsYours bool `json:"y,omitempty"` // IsYours determines if the recipient of the pickup event is the one who picked it up. This is used for the client to add it to their inventory.
}

func (e EventPickup) Type() string {
	return "pickup"
}

// EventDrop notifies the client that the given item was dropped.
type EventDrop struct {
	WID     id.WID
	IsYours bool `json:"y,omitempty"` // IsYours determines if the recipient of the drop event is the one who dropped it. This is used for the client to remove it from their inventory.
	X, Y    int  // The position the item is dropped to. Generally this is the same location as the dropper.
}

func (e EventDrop) Type() string {
	return "drop"
}
