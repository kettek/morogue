package gen

import "errors"

type Config interface {
}

var (
	ErrWrongConfig = errors.New("wrong config type provided")
)
