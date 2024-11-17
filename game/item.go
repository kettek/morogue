package game

import "github.com/kettek/morogue/id"

// ItemArchetype is effectively a blueprint for an item.
type ItemArchetype struct {
	ID    id.UUID
	Title string `msgpack:"T,omitempty"`
	Image string `msgpack:"i,omitempty"`
}

// Type returns "item".
func (a ItemArchetype) Type() string {
	return "item"
}

// GetID returns the ID of the archetype.
func (a ItemArchetype) GetID() id.UUID {
	return a.ID
}

// Item represents a generic item in the world.
type Item struct {
	Objectable
	Position
	Name string `msgpack:"n,omitempty"`
}

// Type returns "item"
func (o Item) Type() ObjectType {
	return "item"
}
