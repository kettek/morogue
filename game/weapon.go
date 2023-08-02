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
	Title              string    `json:"T,omitempty"`
	Image              string    `json:"i,omitempty"`
	PrimaryAttribute   Attribute `json:"p,omitempty"` // Primary attribute to draw damage from.
	SecondaryAttribute Attribute `json:"s,omitempty"` // Secondary attribute to draw 50% damage from.
	MinDamage          int       `json:"m,omitempty"` // Character proficiency with a weapon increases min up to max.
	MaxDamage          int       `json:"M,omitempty"`
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
	Archetype id.UUID `json:"A,omitempty"`
	WID       id.WID
	Container id.WID `json:"c,omitempty"` // The container of the item, if any.
	Applied   bool   `json:"a,omitempty"`
}

func (w Weapon) Type() ObjectType {
	return "weapon"
}

func (w Weapon) GetWID() id.WID {
	return w.WID
}

// GetPosition returns the position of the item.
func (w Weapon) GetPosition() Position {
	return w.Position
}

// SetPosition sets the position of the item.
func (w *Weapon) SetPosition(p Position) {
	w.Position = p
}
