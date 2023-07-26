package game

import "github.com/kettek/morogue/id"

type Location struct {
	ID         id.UUID     `json:"id,omitempty"`
	Tiles      Tiles       `json:"t,omitempty"`
	Mobs       []Mob       `json:"m,omitempty"`
	Characters []Character `json:"c,omitempty"`
	Objects    []Object    `json:"o,omitempty"`
}
