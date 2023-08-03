package game

import "github.com/kettek/morogue/id"

// WeaponType is the type the weapon is considered as.
type WeaponType uint8

const (
	WeaponTypeRange WeaponType = iota
	WeaponTypeMelee
	WeaponTypeUnarmed
)

// WeaponArchetype is effectively a blueprint for a weapon.
type WeaponArchetype struct {
	ID                 id.UUID
	Title              string    `msgpack:"T,omitempty"`
	Image              string    `msgpack:"i,omitempty"`
	PrimaryAttribute   Attribute `msgpack:"p,omitempty"` // Primary attribute to draw damage from.
	SecondaryAttribute Attribute `msgpack:"s,omitempty"` // Secondary attribute to draw 50% damage from.
	MinDamage          int       `msgpack:"m,omitempty"` // Character proficiency with a weapon increases min up to max.
	MaxDamage          int       `msgpack:"M,omitempty"`
}

func (a WeaponArchetype) Type() string {
	return "weapon"
}

func (a WeaponArchetype) GetID() id.UUID {
	return a.ID
}

// Weapon is a weapon.
type Weapon struct {
	Position
	Archetype id.UUID `msgpack:"A,omitempty"`
	WID       id.WID
	Container id.WID `msgpack:"c,omitempty"` // The container of the item, if any.
	Applied   bool   `msgpack:"a,omitempty"`
}

func (w Weapon) Type() ObjectType {
	return "weapon"
}

func (w Weapon) GetWID() id.WID {
	return w.WID
}

func (w *Weapon) SetWID(wid id.WID) {
	w.WID = wid
}

// GetPosition returns the position of the item.
func (w Weapon) GetPosition() Position {
	return w.Position
}

// SetPosition sets the position of the item.
func (w *Weapon) SetPosition(p Position) {
	w.Position = p
}

func (w *Weapon) GetArchetype() id.UUID {
	return w.Archetype
}
