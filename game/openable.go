package game

import "errors"

type Openable struct {
}

func (o *Openable) Open(b *Blockable) error {
	if b.BlockType == BlockTypeSolid {
		b.BlockType = BlockTypeNone
		return nil
	}
	return ErrAlreadyOpen
}

func (o *Openable) Close(b *Blockable) error {
	if b.BlockType == BlockTypeNone {
		b.BlockType = BlockTypeSolid
		return nil
	}
	return ErrAlreadyClosed
}

var (
	ErrAlreadyOpen   = errors.New("it is already open")
	ErrAlreadyClosed = errors.New("it is already closed")
)
