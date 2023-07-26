package game

import "github.com/kettek/morogue/id"

type Location struct {
	ID         id.UUID     `json:"id,omitempty"`
	Cells      Cells       `json:"c,omitempty"`
	Mobs       []Mob       `json:"m,omitempty"`
	Characters []Character `json:"ch,omitempty"`
	Objects    []Object    `json:"o,omitempty"`
}
