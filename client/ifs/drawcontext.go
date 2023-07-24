package ifs

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type DrawContext struct {
	Txt    *TextRenderer
	Screen *ebiten.Image
}
