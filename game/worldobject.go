package game

import "github.com/kettek/morogue/id"

// WorldObject contains the basic information for an object to exist in the world.
type WorldObject struct {
	ArchetypeID id.UUID   `msgpack:"A,omitempty"`
	Archetype   Archetype `msgpack:"-" json:"-"`
	WID         id.WID    // ID assigned when entering da world.
	Container   id.WID    `msgpack:"c,omitempty"` // The container of the item, if any.
	ID          id.UUID   `msgpack:"id,omitempty"`
}

// GetWID returns the WID.
func (o WorldObject) GetWID() id.WID {
	return o.WID
}

// SetWID sets the WID.
func (o *WorldObject) SetWID(wid id.WID) {
	o.WID = wid
}

// GetArchetypeID returns the archetype.
func (o *WorldObject) GetArchetypeID() id.UUID {
	return o.ArchetypeID
}

// SetArchetype sets the archetype.
func (o *WorldObject) SetArchetype(a Archetype) {
	o.Archetype = a
}

// GetArchetype returns the archetype.
func (o *WorldObject) GetArchetype() Archetype {
	return o.Archetype
}
