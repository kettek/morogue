package main

import (
	"image"
	"image/color"
	"math"

	"github.com/carlmjohnson/versioninfo"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/morogue/client/embed"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/client/states"
	"github.com/tinne26/etxt"
	"github.com/tinne26/fonts/liberation/lbrtserif"
)

// app is our ebiten "game" interface. It contains and sets up our necessary
// data structures for managing the world. Most importantly, it contains
// a state machine, of which the top-most will have Update and Draw
// calls that pass runContext and drawContext respectively. Each state
// is where to look for actually managing program logic.
type app struct {
	states       []ifs.State
	connectState states.Connect
	drawContext  ifs.DrawContext
	runContext   ifs.RunContext
}

func newApp() *app {
	a := &app{}
	a.runContext.Sm = a

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("morogue")

	embed.Setup()

	ebiten.SetWindowIcon([]image.Image{embed.Icon})

	// copied and pasted from tinne example -- thx tinne! :)
	// create text renderer, set the font and cache
	renderer := etxt.NewRenderer()
	renderer.SetFont(lbrtserif.Font())
	renderer.Utils().SetCache8MiB()

	// adjust main text style properties
	renderer.SetColor(color.RGBA{239, 91, 91, 255})
	renderer.SetAlign(etxt.Center)
	renderer.SetSize(32)

	a.runContext.Txt = ifs.NewTextRenderer(renderer)
	a.drawContext.Txt = a.runContext.Txt

	a.drawContext.UI = &ifs.DrawContextUI{}
	a.runContext.UI = a.drawContext.UI
	a.drawContext.UI.Init(a.drawContext.Txt)

	a.drawContext.Game = &ifs.GameContext{}
	a.runContext.Game = a.drawContext.Game

	a.drawContext.Game.Zoom = 2.0
	a.drawContext.Game.CellWidth = 16
	a.drawContext.Game.CellHeight = 16

	a.connectState.Begin(a.runContext)

	return a
}

// Update is called per tick and updates the connect state and whatever the
// top-most state is.
func (a *app) Update() error {
	if err := a.connectState.Update(a.runContext); err != nil {
		return err
	}

	if t := a.Top(); t != nil {
		t.Update(a.runContext)
	}

	return nil
}

// Draw updates contexts, draws the connect state, the top-most state, and
// also renders version and debugh info.
func (a *app) Draw(screen *ebiten.Image) {
	a.drawContext.Screen = screen
	a.runContext.UI.Width = screen.Bounds().Dx()
	a.runContext.UI.Height = screen.Bounds().Dy()
	if t := a.Top(); t != nil {
		t.Draw(a.drawContext)
	}

	a.connectState.Draw(a.drawContext)

	{
		b := screen.Bounds()
		fs := a.drawContext.Txt.GetSize()
		al := a.drawContext.Txt.GetAlign()
		a.drawContext.Txt.Renderer.SetAlign(etxt.TopBaseline | etxt.Right)
		a.drawContext.Txt.Renderer.SetSize(16)
		a.drawContext.Txt.Renderer.Draw(screen, versioninfo.Short(), b.Dx(), b.Dy()-4)
		a.drawContext.Txt.Renderer.SetSize(fs)
		a.drawContext.Txt.Renderer.SetAlign(al)
	}
}

// Layout does sizing in accordance to HiDPI scaling.
func (a *app) Layout(winWidth, winHeight int) (int, int) {
	scale := ebiten.DeviceScaleFactor()
	a.runContext.Txt.SetScale(scale) // relevant for HiDPI
	canvasWidth := int(math.Ceil(float64(winWidth) * scale))
	canvasHeight := int(math.Ceil(float64(winHeight) * scale))
	return canvasWidth, canvasHeight
}
