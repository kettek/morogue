package game

import "errors"

// Openable is an embed that provides logic for being opened or closed.
type Openable struct {
	Opened bool
}

// IsOpened returns if the Openable is opened.
func (o *Openable) IsOpened() bool {
	return o.Opened
}

// Open opens the Openable.
func (o *Openable) Open() error {
	if o.Opened {
		return ErrAlreadyOpen
	}
	o.Opened = true
	return nil
}

// Close closes the Openable.
func (o *Openable) Close() error {
	if !o.Opened {
		return ErrAlreadyClosed
	}
	o.Opened = false
	return nil
}

// Our openable errors.
var (
	ErrAlreadyOpen   = errors.New(lc.T("it is already open"))
	ErrAlreadyClosed = errors.New(lc.T("it is already closed"))
)
