package ifs

// GameContext provides context specific to the game world.
type GameContext struct {
	PreventMapInput       bool
	Zoom                  float64
	CellWidth, CellHeight int
}
