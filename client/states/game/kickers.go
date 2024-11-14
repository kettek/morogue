package game

import (
	"image/color"
	"math/rand"

	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/game"
	"github.com/tinne26/etxt"
)

// Kickers provides a visual rendering of kickers in the game world.
// It manages their timeout, etc.
type Kickers struct {
	kickers          []*Kicker
	offsetX, offsetY int
}

// Offset returns the current visual offset of the kickers.
func (kickers *Kickers) Offset() (x, y int) {
	return kickers.offsetX, kickers.offsetY
}

// SetOffset sets the current visual offset of the kickers.
func (kickers *Kickers) SetOffset(x, y int) {
	kickers.offsetX = x
	kickers.offsetY = y
}

// Kicker is a given kicker instance.
type Kicker struct {
	offsetX, offsetY int
	Position         game.Position
	Lifetime         int
	Message          string
	Color            color.NRGBA
}

// Add creates and adds a kicker to the world.
func (kickers *Kickers) Add(kicker Kicker) {
	// Lightly randomize x and y offset
	kicker.offsetX = rand.Intn(5) - 2
	kicker.offsetY = rand.Intn(5) - 2
	kickers.kickers = append(kickers.kickers, &kicker)
}

// Update manages the lifetime of kickers.
func (kickers *Kickers) Update() {
	i := 0
	for _, kicker := range kickers.kickers {
		if kicker.Lifetime > 0 {
			if kicker.Lifetime%2 == 0 {
				kicker.offsetY--
			}
			kicker.Lifetime--
			kickers.kickers[i] = kicker
			i++
		}
	}
	for j := i; j < len(kickers.kickers); j++ {
		kickers.kickers[j] = nil
	}
	kickers.kickers = kickers.kickers[:i]
}

// Draw draws the kickers to the provided context's screen.
func (kickers *Kickers) Draw(ctx ifs.DrawContext) {
	cw := int(float64(ctx.Game.CellWidth) * ctx.Game.Zoom)
	ch := int(float64(ctx.Game.CellHeight) * ctx.Game.Zoom)

	ctx.Txt.Save()
	ctx.Txt.SetAlign(etxt.VertCenter | etxt.HorzCenter)
	for _, kicker := range kickers.kickers {
		clr := kicker.Color
		oclr := color.NRGBA{0, 0, 0, 100}
		if kicker.Lifetime < 10 {
			clr.A = uint8(float64(kicker.Lifetime) / 10 * 255)
			oclr.A = uint8(float64(kicker.Lifetime) / 10 * 100)
		}

		x := kicker.Position.X*cw + kickers.offsetX + (cw / 2) + kicker.offsetX
		y := kicker.Position.Y*ch + kickers.offsetY + kicker.offsetY

		ctx.Txt.SetColor(clr)
		ctx.Txt.SetOutlineColor(oclr)
		ctx.Txt.DrawWithOutline(ctx.Screen, kicker.Message, x, y)
	}

	ctx.Txt.Restore()
}
