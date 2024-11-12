package game

import (
	"github.com/kettek/morogue/id"
	"github.com/vmihailenco/msgpack/v5"
)

// Desire represents a desire of the client to make their character do something. These often cause Events to be sent back to the client on location updating.
type Desire interface {
	Type() string
}

// DesireWrapper is for sending desires from the client to the server.
type DesireWrapper struct {
	Type string             `msgpack:"t"`
	Data msgpack.RawMessage `msgpack:"d"`
}

// Desire returns the desire stored in the wrapper.
func (w *DesireWrapper) Desire() Desire {
	switch w.Type {
	case (DesireMove{}).Type():
		var d DesireMove
		msgpack.Unmarshal(w.Data, &d)
		return d
	case (DesireApply{}).Type():
		var d DesireApply
		msgpack.Unmarshal(w.Data, &d)
		return d
	case (DesirePickup{}).Type():
		var d DesirePickup
		msgpack.Unmarshal(w.Data, &d)
		return d
	case (DesireDrop{}).Type():
		var d DesireDrop
		msgpack.Unmarshal(w.Data, &d)
		return d
	case (DesireBash{}).Type():
		var d DesireBash
		msgpack.Unmarshal(w.Data, &d)
		return d
	case (DesireOpen{}).Type():
		var d DesireOpen
		msgpack.Unmarshal(w.Data, &d)
		return d
	}
	return nil
}

// DesireMove represents the desire to move in a cardinal direction.
type DesireMove struct {
	Direction MoveDirection `msgpack:"d,omitempty"`
}

func (d DesireMove) Type() string {
	return "move"
}

// DesireApply represents the desire to apply or unapply a particular object.
type DesireApply struct {
	WID   id.WID `msgpack:"wid,omitempty"`
	Apply bool   `msgpack:"a,omitempty"` // Whether to apply or unapply.
}

func (d DesireApply) Type() string {
	return "apply"
}

// DesirePickup represents the desire to pick up a particular object.
type DesirePickup struct {
	WID id.WID `msgpack:"wid,omitempty"`
}

func (d DesirePickup) Type() string {
	return "pickup"
}

// DesireDrop represents the desire to drop a particular object.
type DesireDrop struct {
	WID id.WID `msgpack:"wid,omitempty"`
}

func (d DesireDrop) Type() string {
	return "drop"
}

// DesireBash represents the desire to bash a particular object or direction.
type DesireBash struct {
	WID       id.WID        `msgpack:"wid,omitempty"`
	Direction MoveDirection `msgpack:"d,omitempty"`
}

func (d DesireBash) Type() string {
	return "bash"
}

// DesireOpen represents the desire to open or close a particular object.
type DesireOpen struct {
	WID  id.WID `msgpack:"wid,omitempty"`
	Open bool   `msgpack:"o,omitempty"` // Whether to open or close.
}

func (d DesireOpen) Type() string {
	return "open"
}

// DesirePing represents the desire to ping a location or WID to other players.
type DesirePing struct {
	WID      id.WID   `msgpack:"wid,omitempty"`
	Position Position `msgpack:"p,omitempty"`
}

func (d DesirePing) Type() string {
	return "ping"
}
