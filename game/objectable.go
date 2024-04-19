package game

import (
	"github.com/kettek/morogue/id"
)

// Objectable contains the basic information for an object to exist in the world.
type Objectable struct {
	ArchetypeID id.UUID   `msgpack:"A,omitempty"`
	Archetype   Archetype `msgpack:"-" json:"-"`
	WID         id.WID    `json:"-"`                       // ID assigned when entering da world.
	Container   id.WID    `msgpack:"c,omitempty" json:"-"` // The container of the item, if any.
	//ID          id.UUID   `msgpack:"id,omitempty"`
}

// GetWID returns the WID.
func (o Objectable) GetWID() id.WID {
	return o.WID
}

// SetWID sets the WID.
func (o *Objectable) SetWID(wid id.WID) {
	o.WID = wid
}

// GetArchetypeID returns the archetype.
func (o *Objectable) GetArchetypeID() id.UUID {
	return o.ArchetypeID
}

// SetArchetype sets the archetype.
func (o *Objectable) SetArchetype(a Archetype) {
	o.Archetype = a
}

// GetArchetype returns the archetype.
func (o *Objectable) GetArchetype() Archetype {
	return o.Archetype
}

// SetContainerWID sets the container.
func (o *Objectable) SetContainerWID(wid id.WID) {
	o.Container = wid
}

// GetContainerWID returns the container.
func (o Objectable) GetContainerWID() id.WID {
	return o.Container
}
