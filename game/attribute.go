package game

import "image/color"

// Attribute represents a particular physical or mental ability of a character.
type Attribute uint8

const (
	// AttributeSwole represents physicality
	AttributeSwole Attribute = iota
	// AttributeZooms represents speed
	AttributeZooms
	// AttributeBrains represents intelligence
	AttributeBrains
	// AttributeFunk represents the funk
	AttributeFunk
)

const (
	AttributeSwoleDescription  = "Swole determines damage and health"
	AttributeZoomsDescription  = "Zooms determines speed and dodge"
	AttributeBrainsDescription = "Brains determines spell damage and ability to by-pass traps"
	AttributeFunkDescription   = "Funk determines luck and area of effect bonuses"
)

// AttributeLevel represents the level of an attribute, with the whole number
// representing the level and the fractional representing the experience until
// the next level.
type AttributeLevel float64

var (
	ColorSwole         = color.NRGBA{128, 32, 32, 255}
	ColorSwoleVibrant  = color.NRGBA{200, 100, 100, 255}
	ColorZooms         = color.NRGBA{128, 128, 32, 255}
	ColorZoomsVibrant  = color.NRGBA{200, 200, 100, 255}
	ColorBrains        = color.NRGBA{32, 32, 128, 255}
	ColorBrainsVibrant = color.NRGBA{100, 100, 200, 255}
	ColorFunk          = color.NRGBA{128, 32, 128, 255}
	ColorFunkVibrant   = color.NRGBA{200, 100, 200, 255}
)
