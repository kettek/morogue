package game

import (
	"math"

	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kettek/morogue/client/ifs"
)

type Scroller struct {
	lastX, lastY   int
	x, y           int
	limitX, limitY int
	held           bool
	moved          float64
	maxW, maxH     int
	handler        func(x, y int)
}

func (scroller *Scroller) Init() {
	scroller.lastX, scroller.lastY = scroller.GetConstrainedMousePosition()
}

func (scroller *Scroller) Update(ctx ifs.RunContext) error {
	scroller.maxW = ctx.UI.Width
	scroller.maxH = ctx.UI.Height
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		scroller.held = true
		scroller.moved = 0
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if scroller.moved > 2 {
			input.SetCursorShape(input.CURSOR_DEFAULT)
		}
		scroller.held = false
		scroller.moved = 0
	}
	x, y := scroller.GetConstrainedMousePosition()
	if scroller.held {
		if scroller.moved > 2 {
			input.SetCursorShape("move")
			scroller.SetScroll(
				scroller.x+(x-scroller.lastX),
				scroller.y+(y-scroller.lastY),
			)
		} else {
			scroller.moved += math.Abs(float64(x-scroller.lastX)) + math.Abs(float64(y-scroller.lastY))
		}
	}
	scroller.lastX = x
	scroller.lastY = y

	return nil
}

func (scroller *Scroller) SetScroll(x, y int) {
	var minX, minY, maxX, maxY int

	if scroller.limitX > scroller.maxW {
		minX = -(scroller.limitX - scroller.maxW)
		maxX = 0
	} else {
		maxX = scroller.maxW - scroller.limitX
	}
	if scroller.limitY > scroller.maxH {
		minY = -(scroller.limitY - scroller.maxH)
		maxY = 0
	} else {
		maxY = scroller.maxH - scroller.limitY
	}

	if x < minX {
		x = minX
	} else if x > maxX {
		x = maxX
	}
	if y < minY {
		y = minY
	} else if y > maxY {
		y = maxY
	}
	scroller.x = x
	scroller.y = y
	if scroller.handler != nil {
		scroller.handler(scroller.x, scroller.y)
	}
}

func (scroller *Scroller) Scroll() (int, int) {
	return scroller.x, scroller.y
}

func (scroller *Scroller) SetLimit(x, y int) {
	scroller.limitX = x
	scroller.limitY = y
}

func (scroller *Scroller) Limit() (int, int) {
	return scroller.limitX, scroller.limitY
}

func (scroller *Scroller) CenterTo(x, y int) {
	scroller.SetScroll(x, y)
}

func (scoller *Scroller) SetHandler(cb func(x, y int)) {
	scoller.handler = cb
}

func (scroller *Scroller) GetConstrainedMousePosition() (int, int) {
	w, h := ebiten.WindowSize()
	x, y := ebiten.CursorPosition()
	if x < 0 {
		x = 0
	} else if x > w && w != 0 {
		x = w
	}
	if y < 0 {
		y = 0
	} else if y > h && h != 0 {
		y = h
	}

	return x, y
}
