package game

import (
	"errors"

	"github.com/kettek/morogue/gen"
	"github.com/kettek/morogue/id"
)

type Cell struct {
	TileID  *id.UUID     `json:"id,omitempty"` // The Tile ID of the cell.
	Blocks  MovementType `json:"b,omitempty"`  // Whether the cell blocks. This should be generated from the TileID and the contained Objects.
	Objects []Object     `json:"o,omitempty"`  // Non-thinking/active objects. These will generally be weapons, armor, gold, food, etc.
	//
	value int
	flags gen.Flags
}

func (t *Cell) Value() int {
	return t.value
}
func (t *Cell) SetValue(v int) {
	t.value = v
}
func (t *Cell) Flags() gen.Flags {
	return t.flags
}
func (t *Cell) SetFlags(v gen.Flags) {
	t.flags = v
}
func (t *Cell) Data() interface{} {
	return nil
}
func (t *Cell) SetData(d interface{}) {}

type Cells [][]Cell

func NewCells(w, h int) (cells [][]Cell) {
	for x := 0; x < w; x++ {
		cells = append(cells, make([]Cell, 0))
		for y := 0; y < h; y++ {
			cells[x] = append(cells[x], Cell{})
		}
	}
	return
}

func (c Cells) At(x, y int) (Cell, error) {
	if x < 0 || y < 0 || x >= len(c) || y >= len(c[0]) {
		return Cell{}, ErrOutOfBoundCell
	}
	return c[x][y], nil
}

// TODO: Replace MovementType with bitflag
type MovementType uint8

const (
	MovementNone MovementType = iota
	MovementAll
	MovementWalk
	MovementSwim
	MovementHover
	MovementFly
)

var (
	ErrOutOfBoundCell = errors.New("oob cell")
)
