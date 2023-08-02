package game

import "github.com/kettek/morogue/id"

// Archetype is an interface that all archetypes must implement.
type Archetype interface {
	GetID() id.UUID
	Type() string
}
