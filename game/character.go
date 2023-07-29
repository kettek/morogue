package game

import "github.com/kettek/morogue/id"

// Character represents a playable character.
type Character struct {
	Position
	Desire     Desire             `json:"-"` // The current desire of the character. Used server-side.
	LastDesire Desire             `json:"-"` // Last desire processed. Used server-side.
	WID        id.WID             // ID assigned when entering da world.
	Archetype  id.UUID            `json:"a,omitempty"`
	Name       string             `json:"n,omitempty"`
	Level      int                `json:"l,omitempty"`
	Skills     map[string]float64 `json:"s,omitempty"`
	Inventory  []Object           `json:"i,omitempty"`
}

// Type returns "character"
func (c Character) Type() ObjectType {
	return "character"
}

// GetWID returns the WID.
func (c Character) GetWID() id.WID {
	return c.WID
}
