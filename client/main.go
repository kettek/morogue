package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	a := newApp()

	if err := ebiten.RunGame(a); err != nil {
		panic(err)
	}
}
