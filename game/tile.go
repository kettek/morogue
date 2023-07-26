package game

import "github.com/kettek/morogue/id"

type Tile struct {
	ID      id.UUID      `json:"id,omitempty"`
	Blocks  MovementType `json:"b,omitempty"`
	Objects []Object     `json:"o,omitempty"` // Non-thinking/active objects. These will generally be weapons, armor, gold, food, etc.
}

type Tiles [][]Tile

func NewTiles(w, h int) (tiles [][]Tile) {
	for x := 0; x < w; x++ {
		tiles = append(tiles, make([]Tile, 0))
		for y := 0; y < h; y++ {
			tiles[x] = append(tiles[x], Tile{})
		}
	}
	return
}

// TODO: Replace MovementType with bitflag
type MovementType uint8

const (
	MovementWalk MovementType = iota
	MovementSwim
	MovementHover
	MovementFly
)
