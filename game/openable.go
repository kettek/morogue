package game

import "errors"

type Openable struct {
	Opened bool
}

func (o *Openable) IsOpened() bool {
	return o.Opened
}

func (o *Openable) Open() error {
	if o.Opened {
		return ErrAlreadyOpen
	}
	o.Opened = true
	return nil
}

func (o *Openable) Close() error {
	if !o.Opened {
		return ErrAlreadyClosed
	}
	o.Opened = false
	return nil
}

var (
	ErrAlreadyOpen   = errors.New("it is already open")
	ErrAlreadyClosed = errors.New("it is already closed")
)
