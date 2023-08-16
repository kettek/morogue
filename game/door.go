package game

import (
	"github.com/kettek/morogue/id"
)

// DoorArchetype is the archetype for a door.
type DoorArchetype struct {
	ID    id.UUID
	Title string `msgpack:"T,omitempty"`
	Image string `msgpack:"i,omitempty"`
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
	Position
	Lockable
	Hurtable
	ArchetypeID id.UUID   `msgpack:"A,omitempty"`
	Archetype   Archetype `msgpack:"-" json:"-"`
	WID         id.WID    // ID assigned when entering da world.
}

// Type returns the type of this object.
func (o Door) Type() ObjectType {
	return "door"
}

// GetWID returns the WID of this object.
func (o Door) GetWID() id.WID {
	return o.WID
}

// SetWID sets the WID of this object.
func (o *Door) SetWID(wid id.WID) {
	o.WID = wid
}

// GetPosition returns the position of this object.
func (o Door) GetPosition() Position {
	return o.Position
}

// SetPosition sets the position of this object.
func (o *Door) SetPosition(pos Position) {
	o.Position = pos
}

// GetArchetypeID returns the archetype ID of this object.
func (o Door) GetArchetypeID() id.UUID {
	return o.ArchetypeID
}

// SetArchetype sets the archetype of this object.
func (o *Door) SetArchetype(a Archetype) {
	o.Archetype = a
}

// GetArchetype returns the archetype of this object.
func (o Door) GetArchetype() Archetype {
	return o.Archetype
}
