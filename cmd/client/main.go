package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/morogue/client"
)

func main() {
	a := client.NewApp()

	if err := ebiten.RunGame(a); err != nil {
		panic(err)
	}
}
