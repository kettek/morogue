package game

import "image/color"

type Attribute uint8

const (
	AttributeSwole Attribute = iota
	AttributeZooms
	AttributeBrains
	AttributeFunk
)

const (
	AttributeSwoleDescription  = "Swole determines damage and health"
	AttributeZoomsDescription  = "Zooms determines speed and dodge"
	AttributeBrainsDescription = "Brains determines spell damage and ability to by-pass traps"
	AttributeFunkDescription   = "Funk determines luck and area of effect bonuses"
)

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
