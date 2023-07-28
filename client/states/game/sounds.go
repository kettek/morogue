package game

import (
	"image/color"

	"github.com/kettek/morogue/client/ifs"
)

type Sounds struct {
	sounds           []*sound
	offsetX, offsetY int
}

func (sounds *Sounds) Offset() (x, y int) {
	return sounds.offsetX, sounds.offsetY
}

func (sounds *Sounds) SetOffset(x, y int) {
	sounds.offsetX = x
	sounds.offsetY = y
}

type sound struct {
	offsetX, offsetY int
	x, y             int
	lifetime         int
	message          string
}

func (sounds *Sounds) Add(x, y int, message string) {
	// Replace sounds at same position. TODO: Maybe vertical stack sounds in same position?
	for _, s := range sounds.sounds {
		if s.x == x && s.y == y {
			s.message = message
			s.lifetime = 10 * len(message)
			return
		}
	}

	sounds.sounds = append(sounds.sounds, &sound{
		x:        x,
		y:        y,
		lifetime: 10 * len(message),
		message:  message,
	})
}

func (sounds *Sounds) Update() {
	i := 0
	for _, sound := range sounds.sounds {
		if sound.lifetime > 0 {
			sound.lifetime--
			sounds.sounds[i] = sound
			i++
		}
	}
	for j := i; j < len(sounds.sounds); j++ {
		sounds.sounds[j] = nil
	}
	sounds.sounds = sounds.sounds[:i]
}

func (sounds *Sounds) Draw(ctx ifs.DrawContext) {

	prevSize := ctx.Txt.Renderer.GetSize()
	prevClr := ctx.Txt.Renderer.GetColor()
	ctx.Txt.Renderer.SetSize(16)
	for _, sound := range sounds.sounds {
		clr := color.NRGBA{225, 225, 225, 255}
		if sound.lifetime < 10 {
			clr.A = uint8(float64(sound.lifetime) / 10 * 255)
		}
		x := sound.x*16*2 + sounds.offsetX + (16 * 2 / 2)
		y := sound.y*16*2 + sounds.offsetY // + (16 * 2 / 2)
		ctx.Txt.Renderer.SetColor(clr)
		ctx.Txt.Renderer.Draw(ctx.Screen, sound.message, x, y)
	}
	ctx.Txt.Renderer.SetColor(prevClr)
	ctx.Txt.Renderer.SetSize(prevSize)
}
