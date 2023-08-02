package game

import (
	"math"

	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kettek/morogue/client/ifs"
)

// Scroller allows an area within the limits of the window screen.
// It will allow areas smaller than the screen to be moved within
// the screen's bounds and will allow areas beyond the screen
// to be scrolled into view, up to the screen's bounds.
type Scroller struct {
	lastX, lastY   int
	x, y           int
	limitX, limitY int
	held           bool
	moved          float64
	maxW, maxH     int
	handler        func(x, y int)
}

// Init sets some initial state.
func (scroller *Scroller) Init() {
	scroller.lastX, scroller.lastY = scroller.GetConstrainedMousePosition()
}

// Update does the bulk of things.
func (scroller *Scroller) Update(ctx ifs.RunContext) error {
	if ctx.Game.PreventMapInput {
		scroller.held = false
		return nil
	}

	scroller.maxW = ctx.UI.Width
	scroller.maxH = ctx.UI.Height

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !scroller.held {
		scroller.held = true
		scroller.moved = 0
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && scroller.held {
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

// SetScroll sets the current scroll position. This will call the handler if one is set.
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

// Scroll returns the current scroll positions.
func (scroller *Scroller) Scroll() (int, int) {
	return scroller.x, scroller.y
}

// SetLimit sets the scroller's bounds.
func (scroller *Scroller) SetLimit(x, y int) {
	scroller.limitX = x
	scroller.limitY = y
}

// Limit returns the scroller's bounds.
func (scroller *Scroller) Limit() (int, int) {
	return scroller.limitX, scroller.limitY
}

// CenterTo is supposed to center the scroller on the given coordinate, but
// currently just called SetScroll.
func (scroller *Scroller) CenterTo(x, y int) {
	scroller.SetScroll(x, y)
}

// SetHandler sets the scroller's handler.
func (scoller *Scroller) SetHandler(cb func(x, y int)) {
	scoller.handler = cb
}

// GetConstrainedMousePosition gets the mouse position fully constrained to
// the current window's dimensions.
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
