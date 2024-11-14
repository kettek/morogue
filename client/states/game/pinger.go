package game

import (
	"image/color"
	"math"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kettek/morogue/client/embed"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/game"
)

type ping struct {
	offsetX, offsetY int
	x, y             int
	lifetime         int
	kind             string
}

var pingKinds = [4]string{"look", "target", "defend", "danger"}

type Pinger struct {
	offsetX, offsetY int
	clickX, clickY   int
	endX, endY       int
	container        *widget.Container
	open             bool
	PingLocation     func(x, y int, kind string)
	pings            []*ping
	activeKind       string
}

func (pinger *Pinger) Init(container *widget.Container, ctx ifs.RunContext) {
	pinger.container = container
	// TODO: Create ping radial menu.
}

func (pinger *Pinger) SetOffset(x, y int) {
	pinger.offsetX = x
	pinger.offsetY = y
}

func (pinger *Pinger) Add(to game.Position, kind string) {
	pinger.pings = append(pinger.pings, &ping{
		x:        to.X,
		y:        to.Y,
		kind:     kind,
		lifetime: 200,
	})
}

func (pinger *Pinger) ActiveKind() string {
	distance := math.Sqrt(math.Pow(float64(pinger.endX-pinger.clickX), 2) + math.Pow(float64(pinger.endY-pinger.clickY), 2))
	if distance < float64(embed.PingBasic.Bounds().Dx()) {
		return "basic"
	}
	rads := math.Atan2(float64(pinger.endY-pinger.clickY), float64(pinger.endX-pinger.clickX))
	rads += math.Pi / 4
	if rads < 0 {
		rads += 2 * math.Pi
	}
	index := int(rads / (2 * math.Pi) * float64(len(pingKinds)))
	kind := pingKinds[index]

	return kind
}

func (pinger *Pinger) Update(ctx ifs.RunContext) {
	if !pinger.open {
		if ebiten.IsKeyPressed(ebiten.KeyAlt) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			pinger.open = true
			pinger.clickX, pinger.clickY = ebiten.CursorPosition()
			ctx.Game.PreventMapInput = true
		}
	} else {
		pinger.endX, pinger.endY = ebiten.CursorPosition()
		pinger.activeKind = pinger.ActiveKind()
		if !(ebiten.IsKeyPressed(ebiten.KeyAlt) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)) {
			x, y := pinger.getPosition()
			pinger.PingLocation(x, y, pinger.activeKind)
			pinger.open = false
			ctx.Game.PreventMapInput = false
		}
	}
	// Process pings.
	i := 0
	for _, ping := range pinger.pings {
		if ping.lifetime > 0 {
			ping.lifetime--
			pinger.pings[i] = ping
			i++
		}
	}
	for j := i; j < len(pinger.pings); j++ {
		pinger.pings[j] = nil
	}
	pinger.pings = pinger.pings[:i]
}

func (pinger *Pinger) Draw(ctx ifs.DrawContext) {
	// Draw our pings.
	for _, ping := range pinger.pings {
		img := embed.PingBasic
		if ping.kind == "danger" {
			img = embed.PingDanger
		} else if ping.kind == "look" {
			img = embed.PingLook
		} else if ping.kind == "target" {
			img = embed.PingTarget
		} else if ping.kind == "defend" {
			img = embed.PingDefend
		}

		x := ping.x + pinger.offsetX - img.Bounds().Dx()/2
		y := ping.y + pinger.offsetY - img.Bounds().Dy()/2
		opts := ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x), float64(y))

		alpha := float32(ping.lifetime) / 200
		opts.ColorScale.ScaleAlpha(alpha)

		ctx.Screen.DrawImage(img, &opts)
	}
	if pinger.open {
		vector.DrawFilledCircle(ctx.Screen, float32(pinger.clickX), float32(pinger.clickY), float32((len(pingKinds)-2)*embed.PingBasic.Bounds().Dx()), color.NRGBA{0, 0, 0, 100}, false)

		opts := ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(pinger.clickX)-float64(embed.PingBasic.Bounds().Dx()/2), float64(pinger.clickY)-float64(embed.PingBasic.Bounds().Dy()/2))
		ctx.Screen.DrawImage(embed.PingBasic, &opts)

		for i, kind := range pingKinds {
			img := embed.PingDefend
			if kind == "danger" {
				img = embed.PingDanger
			} else if kind == "look" {
				img = embed.PingLook
			} else if kind == "target" {
				img = embed.PingTarget
			}

			angle := float64(i) * (math.Pi / 2)
			x := pinger.clickX + int(math.Cos(angle)*float64(len(pingKinds)*img.Bounds().Dx()/2))
			y := pinger.clickY + int(math.Sin(angle)*float64(len(pingKinds)*img.Bounds().Dy()/2))
			opts := ebiten.DrawImageOptions{}
			if kind != pinger.activeKind {
				opts.ColorScale.ScaleAlpha(0.5)
			}
			opts.GeoM.Translate(float64(x-img.Bounds().Dx()/2), float64(y-img.Bounds().Dy()/2))
			ctx.Screen.DrawImage(img, &opts)
		}
	}
}

func (pinger *Pinger) getPosition() (x, y int) {
	x, y = pinger.clickX, pinger.clickY
	x -= pinger.offsetX
	y -= pinger.offsetY
	return
}
