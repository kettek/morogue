package game

import (
	"errors"

	"github.com/kettek/morogue/gen"
	"github.com/kettek/morogue/id"
)

// Cell is a single cell in the game world.
type Cell struct {
	TileID  *id.UUID     `msgpack:"id,omitempty"` // The Tile ID of the cell.
	Blocks  MovementType `msgpack:"b,omitempty"`  // Whether the cell blocks. This should be generated from the TileID and the contained Objects.
	Objects Objects      `msgpack:"o,omitempty"`  // Non-thinking/active objects. These will generally be weapons, armor, gold, food, etc.
	//
	value int       `msgpack:"-"`
	flags gen.Flags `msgpack:"-"`
}

// Value returns the value of the cell.
func (t *Cell) Value() int {
	return t.value
}

// SetValue sets the value of the cell.
func (t *Cell) SetValue(v int) {
	t.value = v
}

// Flags returns the flags of the cell.
func (t *Cell) Flags() gen.Flags {
	return t.flags
}

// SetFlags sets the flags of the cell.
func (t *Cell) SetFlags(v gen.Flags) {
	t.flags = v
}

// Data returns the data of the cell.
func (t *Cell) Data() interface{} {
	return nil
}

// SetData sets the data of the cell.
func (t *Cell) SetData(d interface{}) {}

// Cells is a 2D array of Cells.
type Cells [][]Cell

// NewCells creates a new 2D array of Cells.
func NewCells(w, h int) (cells [][]Cell) {
	for x := 0; x < w; x++ {
		cells = append(cells, make([]Cell, 0))
		for y := 0; y < h; y++ {
			cells[x] = append(cells[x], Cell{})
		}
	}
	return
}

// At returns the cell at a given coordinate.
func (c Cells) At(x, y int) (Cell, error) {
	if x < 0 || y < 0 || x >= len(c) || y >= len(c[0]) {
		return Cell{}, ErrOutOfBoundCell
	}
	return c[x][y], nil
}

// TODO: Replace MovementType with bitflag
// MoveType is the type of movement.
type MovementType uint8

// Our movement types.
const (
	MovementNone MovementType = iota
	MovementAll
	MovementWalk
	MovementSwim
	MovementHover
	MovementFly
)

// Our cell-related errors.
var (
	ErrOutOfBoundCell = errors.New("oob cell")
)
