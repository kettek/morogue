package game

// Mob represents an NPC.
type Mob struct {
	Position
	HP     int
	MaxHP  int
	Swole  AttributeLevel
	Zooms  AttributeLevel
	Brains AttributeLevel
	Funk   AttributeLevel
}
