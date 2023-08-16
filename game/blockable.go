package game

import (
	"errors"
)

type BlockType int

const (
	BlockTypeNone BlockType = iota
	BlockTypeSolid
	BlockTypeHole
	BlockTypeLiquid
)

type Blockable struct {
	BlockType BlockType `msgpack:"b"`
}

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

var (
	ErrInvalidBlockType = errors.New("invalid block type")
)
