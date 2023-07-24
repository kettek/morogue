package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/carlmjohnson/versioninfo"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/client/states"
	"github.com/tinne26/etxt"
	"github.com/tinne26/fonts/liberation/lbrtserif"
)

type app struct {
	stateMachine
	drawContext ifs.DrawContext
	runContext  ifs.RunContext
}

func newApp() *app {
	a := &app{}

	ebiten.SetWindowSize(1280, 720)

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
	a.drawContext.Txt = ifs.NewTextRenderer(renderer)

	if err := a.Push(&states.Connect{}); err != nil {
		panic(err)
	}

	return a
}

func (a *app) Update() error {
	if t := a.stateMachine.Top(); t != nil {
		t.Update(a.runContext)
	}

	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	a.drawContext.Screen = screen
	if t := a.stateMachine.Top(); t != nil {
		t.Draw(a.drawContext)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%s %s", versioninfo.Short(), versioninfo.Revision))
}

func (a *app) Layout(winWidth, winHeight int) (int, int) {
	scale := ebiten.DeviceScaleFactor()
	a.runContext.Txt.SetScale(scale) // relevant for HiDPI
	canvasWidth := int(math.Ceil(float64(winWidth) * scale))
	canvasHeight := int(math.Ceil(float64(winHeight) * scale))
	return canvasWidth, canvasHeight
}
