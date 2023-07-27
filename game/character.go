package game

import "github.com/kettek/morogue/id"

// Character represents a playable character.
type Character struct {
	Position
	WID       id.WID // ID assigned when entering da world.
	Archetype id.UUID
	Name      string
	Level     int
	Skills    map[string]float64
	Inventory []Item
}

// Item represents an item in the world
type Item struct {
	Name string
	ID   int
}
