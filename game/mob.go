package game

import "github.com/kettek/morogue/id"

// Mob represents an NPC.
type Mob struct {
	Position
	WID    id.WID // ID assigned when entering da world.
	HP     int
	MaxHP  int
	Swole  AttributeLevel
	Zooms  AttributeLevel
	Brains AttributeLevel
	Funk   AttributeLevel
}

// Type returns "mob"
func (m Mob) Type() ObjectType {
	return "mob"
}

// GetWID returns the WID.
func (m *Mob) GetWID() id.WID {
	return m.WID
}
