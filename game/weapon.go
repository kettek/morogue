package game

import (
	"fmt"

	"github.com/kettek/morogue/id"
)

// WeaponType is the type the weapon is considered as.
type WeaponType uint8

const (
	WeaponTypeMelee WeaponType = iota
	WeaponTypeRange
	WeaponTypeThrown
	WeaponTypeUnarmed
)

func (a WeaponType) String() string {
	switch a {
	case WeaponTypeRange:
		return "range"
	case WeaponTypeThrown:
		return "thrown"
	case WeaponTypeMelee:
		return "melee"
	case WeaponTypeUnarmed:
		return "unarmed"
	default:
		return ""
	}
}

func (a *WeaponType) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case `"range"`:
		*a = WeaponTypeRange
	case `"thrown"`:
		*a = WeaponTypeThrown
	case `"melee"`:
		*a = WeaponTypeMelee
	case `"unarmed"`:
		*a = WeaponTypeUnarmed
	default:
		*a = WeaponTypeRange
	}
	return nil
}

// WeaponArchetype is effectively a blueprint for a weapon.
type WeaponArchetype struct {
	ID                 id.UUID
	Title              string     `msgpack:"T,omitempty"`
	Image              string     `msgpack:"i,omitempty"`
	PrimaryAttribute   Attribute  `msgpack:"p,omitempty"` // Primary attribute to draw damage from.
	SecondaryAttribute Attribute  `msgpack:"s,omitempty"` // Secondary attribute to draw 50% damage from.
	Description        string     `msgpack:"d,omitempty"`
	WeaponType         WeaponType `msgpack:"w,omitempty"`
	MinDamage          int        `msgpack:"m,omitempty"` // Character proficiency with a weapon increases min up to max.
	MaxDamage          int        `msgpack:"M,omitempty"`
	Slots              Slots      `msgpack:"S,omitempty"`
}

func (a WeaponArchetype) Type() string {
	return "weapon"
}

func (a WeaponArchetype) GetID() id.UUID {
	return a.ID
}

func (a WeaponArchetype) RangeString() string {
	if a.MinDamage == 0 {
		return fmt.Sprintf("〜%d", a.MaxDamage)
	}
	if a.MinDamage == a.MaxDamage {
		return fmt.Sprintf("%d", a.MaxDamage)
	}
	return fmt.Sprintf("%d〜%d", a.MinDamage, a.MaxDamage)
}

// Weapon is a weapon.
type Weapon struct {
	Position
	ArchetypeID id.UUID   `msgpack:"A,omitempty"`
	Archetype   Archetype `msgpack:"-" json:"-"`
	WID         id.WID
	Container   id.WID `msgpack:"c,omitempty"` // The container of the item, if any.
	Applied     bool   `msgpack:"a,omitempty"`
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

func (w *Weapon) GetArchetypeID() id.UUID {
	return w.ArchetypeID
}

func (w *Weapon) GetArchetype() Archetype {
	return w.Archetype
}

func (w *Weapon) SetArchetype(a Archetype) {
	w.Archetype = a
}
