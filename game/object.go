package game

import "github.com/kettek/morogue/id"

type Object struct {
	Position
	WID  id.WID  // ID assigned when entering da world.
	ID   id.UUID `json:"id,omitempty"`
	Name string  `json:"n,omitempty"`
}
