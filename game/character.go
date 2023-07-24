package game

// Character represents a playable character.
type Character struct {
	Name      string
	Level     int
	Skills    map[string]float64
	Inventory []Item
}

type Item struct {
	Name string
	ID   int
}
