package game

import (
	"errors"
)

// BlockType is the type of blocking an object is.
type BlockType int

// Our block types.
const (
	BlockTypeNone BlockType = iota
	BlockTypeSolid
	BlockTypeHole
	BlockTypeLiquid
)

// Blockable is an embed that provides logic for being blocked.
type Blockable struct {
	BlockType BlockType `msgpack:"b"`
}

// UnmarshalJSON unmarhsals a string into our BlockType.
func (b *BlockType) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"none"`:
		*b = BlockTypeNone
	case `"solid"`:
		*b = BlockTypeSolid
	case `"hole"`:
		*b = BlockTypeHole
	case `"liquid"`:
		*b = BlockTypeLiquid
	default:
		*b = BlockTypeSolid
	}
	return nil
}

// IsBlocked returns if the Blockable is BlockTypeNone.
func (b *Blockable) IsBlocked() bool {
	return b.BlockType != BlockTypeNone
}

// Our block type errors.
var (
	ErrInvalidBlockType = errors.New("invalid block type")
)
