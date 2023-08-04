package game

import (
	"fmt"

	"github.com/kettek/morogue/id"
)

// ArmorType is the type the armor is considered as.
type ArmorType uint8

const (
	ArmorTypeNone ArmorType = iota
	ArmorTypeLight
	ArmorTypeMedium
	ArmorTypeHeavy
)

func (a ArmorType) String() string {
	switch a {
	case ArmorTypeNone:
		return "petty"
	case ArmorTypeLight:
		return "light"
	case ArmorTypeMedium:
		return "medium"
	case ArmorTypeHeavy:
		return "heavy"
	default:
		return ""
	}
}

func (a *ArmorType) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case `"petty"`:
		*a = ArmorTypeNone
	case `"light"`:
		*a = ArmorTypeLight
	case `"medium"`:
		*a = ArmorTypeMedium
	case `"heavy"`:
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
	Description  string
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

func (a ArmorArchetype) RangeString() string {
	if a.MinArmor == 0 {
		return fmt.Sprintf("〜%d", a.MaxArmor)
	}
	return fmt.Sprintf("%d〜%d", a.MinArmor, a.MaxArmor)
}

// Armor is a weapon.
type Armor struct {
	Position
	Archetype id.UUID `msgpack:"A,omitempty"`
	WID       id.WID
	Container id.WID `msgpack:"c,omitempty"` // The container of the item, if any.
	Applied   bool   `msgpack:"a,omitempty"`
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

func (a *Armor) GetArchetype() id.UUID {
	return a.Archetype
}
