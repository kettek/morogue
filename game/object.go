package game

import "github.com/kettek/morogue/id"

type Object struct {
	ID   id.UUID `json:"id,omitempty"`
	Name string  `json:"n,omitempty"`
}
