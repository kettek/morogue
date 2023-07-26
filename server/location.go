package main

import (
	"fmt"

	"github.com/kettek/morogue/game"
)

type location struct {
	game.Location
	active bool
}

func (l *location) process() error {
	for _, c := range l.Characters {
		fmt.Println("TODO: Handle character", c)
	}
	for _, m := range l.Mobs {
		fmt.Println("TODO: Handle mob", m)
	}
	for _, o := range l.Objects {
		fmt.Println("TODO: Handle object", o)
	}
	return nil
}
