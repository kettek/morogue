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

func (a *ArmorType) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case `light`:
		*a = ArmorTypeLight
	case `medium`:
		*a = ArmorTypeMedium
	case `heavy`:
		*a = ArmorTypeHeavy
	default:
		*a = ArmorTypeNone
	}
	return nil
}

// ArmorArchetype is effectively a blueprint for armour.
type ArmorArchetype struct {
	ID           id.UUID
	Title        string
	Image        string
	ArmorType    ArmorType
	MinArmor     int // Character proficiency with a weapon increases min up to max.
	MaxArmor     int
	ArmorPenalty int // Penalty to movement speed.
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

func (a *Armor) SetWID(wid id.WID) {
	a.WID = wid
}

func (a Armor) GetPosition() Position {
	return a.Position
}

func (a *Armor) SetPosition(p Position) {
	a.Position = p
}
