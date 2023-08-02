package game

import "github.com/kettek/morogue/id"

// ArmorType is the type the armor is considered as.
type ArmorType uint8

const (
	ArmorTypeNone ArmorType = iota
	ArmorTypeLight
	ArmorTypeMedium
	ArmorTypeHeavy
)

// ArmorArchetype is effectively a blueprint for armour.
type ArmorArchetype struct {
	ID           id.UUID
	Title        string    `json:"T,omitempty"`
	Image        string    `json:"i,omitempty"`
	ArmorType    ArmorType `json:"t,omitempty"`
	MinArmor     int       `json:"m,omitempty"` // Character proficiency with a weapon increases min up to max.
	MaxArmor     int       `json:"M,omitempty"`
	ArmorPenalty int       `json:"p,omitempty"` // Penalty to movement speed.
}

func (a ArmorArchetype) Type() string {
	return "armor"
}

func (a ArmorArchetype) GetID() id.UUID {
	return a.ID
}

// Armor is a weapon.
type Armor struct {
	Position
	Archetype id.UUID `json:"A,omitempty"`
	WID       id.WID
	Container id.WID `json:"c,omitempty"` // The container of the item, if any.
	Applied   bool   `json:"a,omitempty"`
}

func (a Armor) Type() ObjectType {
	return "armor"
}

func (a Armor) GetWID() id.WID {
	return a.WID
}

func (a Armor) GetPosition() Position {
	return a.Position
}

func (a *Armor) SetPosition(p Position) {
	a.Position = p
}
