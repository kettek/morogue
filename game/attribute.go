package game

type Attribute uint8

const (
	AttributeSwole Attribute = iota
	AttributeZooms
	AttributeBrains
	AttributeFunk
)

type AttributeLevel float64
