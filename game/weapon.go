package game

import (
	"fmt"
	"image/color"

	"github.com/kettek/morogue/id"
)

// WeaponType is the type the weapon is considered as.
type WeaponType uint8

// Our weapon types.
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
		return lc.T("range")
	case WeaponTypeThrown:
		return lc.T("thrown")
	case WeaponTypeMelee:
		return lc.T("melee")
	case WeaponTypeUnarmed:
		return lc.T("unarmed")
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
	Objectable
	Position
	Appliable
}

// Type returns the type of the item.
func (w Weapon) Type() ObjectType {
	return "weapon"
}
