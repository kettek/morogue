package game

import (
	"fmt"
	"image/color"

	"github.com/kettek/morogue/id"
)

// ArmorType is the type the armor is considered as.
type ArmorType uint8

// Our armor types.
const (
	ArmorTypeNone ArmorType = iota
	ArmorTypeLight
	ArmorTypeMedium
	ArmorTypeHeavy
)

// String returns the string representation of the armor type.
func (a ArmorType) String() string {
	switch a {
	case ArmorTypeNone:
		return lc.T("petty")
	case ArmorTypeLight:
		return lc.T("light")
	case ArmorTypeMedium:
		return lc.T("medium")
	case ArmorTypeHeavy:
		return lc.T("heavy")
	default:
		return ""
	}
}

// Color returns the color associated with the armor type.
func (a ArmorType) Color() color.Color {
	switch a {
	case ArmorTypeNone:
		return color.NRGBA{R: 200, G: 200, B: 200, A: 255}
	case ArmorTypeLight:
		return color.NRGBA{R: 100, G: 100, B: 200, A: 255}
	case ArmorTypeMedium:
		return color.NRGBA{R: 200, G: 200, B: 100, A: 255}
	case ArmorTypeHeavy:
		return color.NRGBA{R: 200, G: 100, B: 100, A: 255}
	default:
		return color.NRGBA{R: 200, G: 200, B: 200, A: 255}
	}
}

// UnmarshalJSON unmarshals the JSON representation of the armor type.
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
	ID          id.UUID
	Title       string
	Image       string
	Description string
	ArmorType   ArmorType
	MinArmor    int // Character proficiency with a weapon increases min up to max.
	MaxArmor    int
	MovePenalty int   // Penalty to movement speed.
	Slots       Slots `msgpack:"S,omitempty"`
}

// Type returns the type of the archetype.
func (a ArmorArchetype) Type() string {
	return "armor"
}

// GetID returns the ID of the archetype.
func (a ArmorArchetype) GetID() id.UUID {
	return a.ID
}

// RangeString returns the armor range of the archetype.
func (a ArmorArchetype) RangeString() string {
	if a.MinArmor == 0 {
		return fmt.Sprintf("〜%d", a.MaxArmor)
	}
	if a.MinArmor == a.MaxArmor {
		return fmt.Sprintf("%d", a.MaxArmor)
	}
	return fmt.Sprintf("%d〜%d", a.MinArmor, a.MaxArmor)
}

// Armor is a weapon.
type Armor struct {
	Objectable
	Position
	Appliable
}

// Type returns the type of the object.
func (a Armor) Type() ObjectType {
	return "armor"
}
