package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/morogue/client"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/config"
)

func main() {
	cfg := config.Init("morogue", &ifs.Configuration{
		Name: "morogue",
	})

	a := client.NewApp(cfg)

	if err := ebiten.RunGame(a); err != nil {
		panic(err)
	}
}
