package game

import "github.com/kettek/morogue/id"

// BagArchetype is for supporting our baggage.
type BagArchetype struct {
	ID          id.UUID
	Title       string `msgpack:"T,omitempty"`
	Description string `msgpack:"d,omitempty"`
	Image       string `msgpack:"i,omitempty"`
	Capacity    int    `msgpack:"c,omitempty"`
	Limit       int    `msgpack:"l,omitempty"`
}

// Type returns "bag".
func (a BagArchetype) Type() string {
	return "bag"
}

// GetID returns the ID of the archetype.
func (a BagArchetype) GetID() id.UUID {
	return a.ID
}

// Bag represents a bag object in the world.
type Bag struct {
	Objectable
	Position
	Containerable
	Name string `msgpack:"n,omitempty"`
}

// Type returns "bag"
func (o Bag) Type() ObjectType {
	return "bag"
}
