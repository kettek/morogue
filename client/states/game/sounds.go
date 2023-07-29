package game

import (
	"image/color"

	"github.com/kettek/morogue/client/ifs"
)

// Sounds provides a visual rendering of sounds in the game world.
// It manages their timeout, etc.
type Sounds struct {
	sounds           []*sound
	offsetX, offsetY int
}

// Offset returns the current visual offset of the sounds.
func (sounds *Sounds) Offset() (x, y int) {
	return sounds.offsetX, sounds.offsetY
}

// SetOffset sets the current visual offset of the sounds.
func (sounds *Sounds) SetOffset(x, y int) {
	sounds.offsetX = x
	sounds.offsetY = y
}

// sound is a given sound instance.
type sound struct {
	offsetX, offsetY int
	x, y             int
	fromX, fromY     int
	lifetime         int
	message          string
}

// Add creates and adds a sound to the world.
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

// Update manages the lifetime of sounds.
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

// Draw draws the sounds to the provided context's screen.
func (sounds *Sounds) Draw(ctx ifs.DrawContext) {
	cw := int(float64(ctx.Game.CellWidth) * ctx.Game.Zoom)
	ch := int(float64(ctx.Game.CellHeight) * ctx.Game.Zoom)

	ctx.Txt.Save()
	ctx.Txt.SetSize(16)
	for _, sound := range sounds.sounds {
		clr := color.NRGBA{225, 225, 225, 255}
		oclr := color.NRGBA{0, 0, 0, 100}
		if sound.lifetime < 10 {
			clr.A = uint8(float64(sound.lifetime) / 10 * 255)
			oclr.A = uint8(float64(sound.lifetime) / 10 * 100)
		}

		x := sound.x*cw + sounds.offsetX + (cw / 2)
		y := sound.y*ch + sounds.offsetY + (ch / 2)

		// Adjust the sound in the direction of where it came from, if available.
		if sound.fromX > sound.x {
			x += cw / 2
		} else if sound.fromX < sound.x {
			x -= cw / 2
		}
		if sound.fromY > sound.y {
			y += ch / 2
		} else if sound.fromY < sound.y {
			y -= ch / 2
		}

		ctx.Txt.SetColor(clr)
		ctx.Txt.SetOutlineColor(oclr)
		ctx.Txt.DrawWithOutline(ctx.Screen, sound.message, x, y)
	}

	ctx.Txt.Restore()
}
