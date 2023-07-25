package main

import (
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

	embed.Setup()

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

	a.connectState.Begin(a.runContext)

	return a
}

func (a *app) Update() error {
	if err := a.connectState.Update(a.runContext); err != nil {
		return err
	}

	if t := a.Top(); t != nil {
		t.Update(a.runContext)
	}

	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	a.drawContext.Screen = screen
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

func (a *app) Layout(winWidth, winHeight int) (int, int) {
	scale := ebiten.DeviceScaleFactor()
	a.runContext.Txt.SetScale(scale) // relevant for HiDPI
	canvasWidth := int(math.Ceil(float64(winWidth) * scale))
	canvasHeight := int(math.Ceil(float64(winHeight) * scale))
	return canvasWidth, canvasHeight
}
