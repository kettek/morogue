package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/id"
)

type Data interface {
	Archetype(id.UUID) game.Archetype
	ArchetypeImage(id.UUID) *ebiten.Image
	EnsureImage(archetype game.Archetype, zoom float64) (*ebiten.Image, error)
	LoadImage(path string, zoom float64) (*ebiten.Image, error)
}
