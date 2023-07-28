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
	fromX, fromY     int
	lifetime         int
	message          string
}

func (sounds *Sounds) Add(message string, x, y int, fromX, fromY int) {
	// Replace sounds at same position. TODO: Maybe vertical stack sounds in same position?
	for _, s := range sounds.sounds {
		if s.x == x && s.y == y && s.fromX == fromX && s.fromY == fromY {
			s.message = message
			s.lifetime = 10 * len(message)
			return
		}
	}

	sounds.sounds = append(sounds.sounds, &sound{
		x:        x,
		y:        y,
		fromX:    fromX,
		fromY:    fromY,
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
	ctx.Txt.Save()

	ctx.Txt.SetSize(16)
	for _, sound := range sounds.sounds {
		clr := color.NRGBA{225, 225, 225, 255}
		oclr := color.NRGBA{0, 0, 0, 100}
		if sound.lifetime < 10 {
			clr.A = uint8(float64(sound.lifetime) / 10 * 255)
			oclr.A = uint8(float64(sound.lifetime) / 10 * 100)
		}

		// FIXME: 16*2 and 16*2/2 are placeholders for cellSize * zoom and halfCellSize * zoom
		x := sound.x*16*2 + sounds.offsetX + (16 * 2 / 2)
		y := sound.y*16*2 + sounds.offsetY + (16 * 2 / 2)

		// Adjust the sound in the direction of where it came from, if available.
		if sound.fromX > sound.x {
			x += 16 * 2 / 2
		} else if sound.fromX < sound.x {
			x -= 16 * 2 / 2
		}
		if sound.fromY > sound.y {
			y += 16 * 2 / 2
		} else if sound.fromY < sound.y {
			y -= 16 * 2 / 2
		}

		ctx.Txt.SetColor(clr)
		ctx.Txt.SetOutlineColor(oclr)
		ctx.Txt.DrawWithOutline(ctx.Screen, sound.message, x, y)
	}

	ctx.Txt.Restore()
}
