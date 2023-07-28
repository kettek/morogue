package game

import (
	"encoding/json"

	"github.com/kettek/morogue/id"
)

type Event interface {
	Type() string
}

// EventWrapper is for sending desires from the client to the server.
type EventWrapper struct {
	Type string          `json:"t"`
	Data json.RawMessage `json:"d"`
}

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

	}
	return nil
}

type EventPosition struct {
	WID  id.WID
	X, Y int
}

func (e EventPosition) Type() string {
	return "position"
}

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
