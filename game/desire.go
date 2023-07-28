package game

import (
	"encoding/json"

	"github.com/kettek/morogue/id"
)

// Desire represents a desire of the client to make their character do something. These often cause Events to be sent back to the client on location updating.
type Desire interface {
	Type() string
}

// DesireWrapper is for sending desires from the client to the server.
type DesireWrapper struct {
	Type string          `json:"t"`
	Data json.RawMessage `json:"d"`
}

// Desire returns the desire stored in the wrapper.
func (w *DesireWrapper) Desire() Desire {
	switch w.Type {
	case (DesireMove{}).Type():
		var d DesireMove
		json.Unmarshal(w.Data, &d)
		return d
	case (DesireApply{}).Type():
		var d DesireApply
		json.Unmarshal(w.Data, &d)
		return d
	}
	return nil
}

// DesireMove represents the desire to move in a cardinal direction.
type DesireMove struct {
	Direction MoveDirection `json:"d,omitempty"`
}

func (d DesireMove) Type() string {
	return "move"
}

// DesireApply represents the desire to apply a particular object.
type DesireApply struct {
	WID id.WID `json:"wid,omitempty"`
}

func (d DesireApply) Type() string {
	return "apply"
}
