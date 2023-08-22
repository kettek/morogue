package game

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kettek/goro/pathing"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/game"
)

type Pather struct {
	offsetX, offsetY int
	targetX, targetY int
	currX, currY     int
	Steps            []pathing.Step
	Location         *game.Location
	MoveTimer        int
}

func (pather *Pather) Init() {
}

func (path *Pather) SetOffset(x, y int) {
	path.offsetX = x
	path.offsetY = y
}

// Update returns the current move direction from the current paths.
func (pather *Pather) Update(character *game.Character) game.Desire {
	if len(pather.Steps) == 0 {
		return nil
	}

	if pather.MoveTimer > 0 {
		pather.MoveTimer--
		return nil
	}
	// FIXME: Make MoveTimer based upon character speed.
	pather.MoveTimer = 8

	step := pather.Steps[0]
	pather.Steps = pather.Steps[1:]

	// Get the direction of the character to the step.
	dx := step.X() - character.X
	dy := step.Y() - character.Y
	var d game.MoveDirection
	if dx < 0 {
		if dy < 0 {
			d = game.UpLeftMoveDirection
		} else if dy > 0 {
			d = game.DownLeftMoveDirection
		} else {
			d = game.LeftMoveDirection
		}
	} else if dx > 0 {
		if dy < 0 {
			d = game.UpRightMoveDirection
		} else if dy > 0 {
			d = game.DownRightMoveDirection
		} else {
			d = game.RightMoveDirection
		}
	} else {
		if dy < 0 {
			d = game.UpMoveDirection
		} else if dy > 0 {
			d = game.DownMoveDirection
		}
	}

	if d != 0 {
		return game.DesireMove{
			Direction: d,
		}
	}
	return nil
}

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	whiteImage.Fill(color.NRGBA{0xff, 0xff, 0xff, 0x88})
}

func (p *Pather) Draw(ctx ifs.DrawContext, character *game.Character) {
	cw := float64(ctx.Game.CellWidth) * ctx.Game.Zoom
	ch := float64(ctx.Game.CellHeight) * ctx.Game.Zoom
	var path vector.Path
	path.MoveTo(float32(float64(character.X)*cw+float64(p.offsetX)+cw/2), float32(float64(character.Y)*ch+float64(p.offsetY)+ch/2))
	for _, step := range p.Steps {
		x, y := float64(step.X()), float64(step.Y())
		x *= cw
		y *= ch
		x += float64(p.offsetX) + cw/2
		y += float64(p.offsetY) + ch/2
		path.LineTo(float32(x), float32(y))
	}
	op := &vector.StrokeOptions{}
	op.LineCap = vector.LineCapRound
	op.LineJoin = vector.LineJoinRound
	op.Width = 2
	vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, op)
	ctx.Screen.DrawTriangles(vs, is, whiteSubImage, &ebiten.DrawTrianglesOptions{})
}
