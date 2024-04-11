package game

import (
	"github.com/kettek/morogue/id"
)

// DoorArchetype is the archetype for a door.
type DoorArchetype struct {
	ID        id.UUID
	Title     string    `msgpack:"T,omitempty"`
	Image     string    `msgpack:"i,omitempty"`
	BlockType BlockType `msgpack:"-"`
	Health    int       `msgpack:"-"`
	MaxHealth int       `msgpack:"-"`
}

// Type returns the type of this archetype.
func (d DoorArchetype) Type() string {
	return "door"
}

// GetID returns the ID of this archetype.
func (d DoorArchetype) GetID() id.UUID {
	return d.ID
}

// Door is a door.
type Door struct {
	WorldObject
	Position
	Hurtable
	Lockable
	Blockable
	Openable
}

// Type returns the type of this object.
func (o Door) Type() ObjectType {
	return "door"
}
