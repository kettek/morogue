package game

import (
	"encoding/json"

	"github.com/kettek/morogue/id"
)

type Desire interface {
	Type() string
}

type CharacterDesire struct {
	WID    id.WID
	Desire Desire
}

// DesireWrapper is for sending desires from the client to the server.
type DesireWrapper struct {
	Type string          `json:"t"`
	Data json.RawMessage `json:"d"`
}

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

type DesireMove struct {
	Direction MoveDirection `json:"d,omitempty"`
}

func (d DesireMove) Type() string {
	return "move"
}

type DesireApply struct {
	WID id.WID `json:"wid,omitempty"`
}

func (d DesireApply) Type() string {
	return "apply"
}
