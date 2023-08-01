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

// GetPosition returns the position.
func (m *Mob) GetPosition() Position {
	return m.Position
}

// SetPosition sets the position.
func (m *Mob) SetPosition(p Position) {
	m.Position = p
}
