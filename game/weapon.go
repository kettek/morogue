package game

import (
	"fmt"
	"image/color"

	"github.com/kettek/morogue/id"
)

// WeaponType is the type the weapon is considered as.
type WeaponType uint8

const (
	WeaponTypeNone WeaponType = iota
	WeaponTypeMelee
	WeaponTypeRange
	WeaponTypeThrown
	WeaponTypeUnarmed
)

// String returns the string representation of the weapon type.
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

// Color returns the color associated with the weapon type.
func (a WeaponType) Color() color.Color {
	switch a {
	case WeaponTypeRange:
		return color.NRGBA{R: 50, G: 250, B: 50, A: 255}
	case WeaponTypeThrown:
		return color.NRGBA{R: 50, G: 250, B: 150, A: 255}
	case WeaponTypeMelee:
		return color.NRGBA{R: 250, G: 50, B: 50, A: 255}
	case WeaponTypeUnarmed:
		return color.NRGBA{R: 250, G: 150, B: 50, A: 255}
	default:
		return color.NRGBA{R: 200, G: 200, B: 200, A: 255}
	}
}

// UnmarshalJSON unmarshals the JSON representation of the weapon type.
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

// Type returns the type of the archetype.
func (a WeaponArchetype) Type() string {
	return "weapon"
}

// GetID returns the ID of the archetype.
func (a WeaponArchetype) GetID() id.UUID {
	return a.ID
}

// RangeString returns the string representation of the damage range.
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
	Appliable
	ArchetypeID id.UUID   `msgpack:"A,omitempty"`
	Archetype   Archetype `msgpack:"-" json:"-"`
	WID         id.WID
	Container   id.WID `msgpack:"c,omitempty"` // The container of the item, if any.
}

// Type returns the type of the item.
func (w Weapon) Type() ObjectType {
	return "weapon"
}

// GetWID returns the WID of the item.
func (w Weapon) GetWID() id.WID {
	return w.WID
}

// SetWID sets the WID of the item.
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

// GetArchetypeID returns the ID of the archetype.
func (w *Weapon) GetArchetypeID() id.UUID {
	return w.ArchetypeID
}

// GetArchetype returns the archetype.
func (w *Weapon) GetArchetype() Archetype {
	return w.Archetype
}

// SetArchetype sets the archetype.
func (w *Weapon) SetArchetype(a Archetype) {
	w.Archetype = a
}
