package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kettek/morogue/client/ifs"
)

// Grid renders a grid over an area.
type Grid struct {
	offsetX, offsetY      int
	width, height         int
	cellWidth, cellHeight int
	color                 color.Color
	image                 *ebiten.Image
	clickHandler          func(x, y int)
	heldHandler           func(x, y int)
}

func (grid *Grid) makeImage() {
	if grid.width == 0 || grid.height == 0 {
		return
	}
	grid.image = ebiten.NewImage(grid.width, grid.height)
	cx := grid.width / grid.cellWidth
	cy := grid.height / grid.cellHeight
	for x := 1; x != cx; x++ {
		vector.StrokeLine(grid.image, float32(x*grid.cellWidth), 0, float32(x*grid.cellWidth), float32(grid.height), 1, grid.color, false)
	}
	for y := 1; y != cy; y++ {
		vector.StrokeLine(grid.image, 0, float32(y*grid.cellHeight), float32(grid.width), float32(y*grid.cellHeight), 1, grid.color, false)
	}
}

// Size returns the grid's size.
func (grid *Grid) Size() (int, int) {
	return grid.width, grid.height
}

// SetSize sets the grid's size.
func (grid *Grid) SetSize(w, h int) {
	grid.width = w
	grid.height = h
	grid.makeImage()
}

// CellSize returns the current cell size.
func (grid *Grid) CellSize() (int, int) {
	return grid.cellWidth, grid.cellHeight
}

// SetCellSize sets the cell size.
func (grid *Grid) SetCellSize(w, h int) {
	grid.cellWidth = w
	grid.cellHeight = h
	grid.makeImage()
}

// Offset returns the grid's placement offset.
func (grid *Grid) Offset() (x, y int) {
	return grid.offsetX, grid.offsetY
}

// SetOffset sets the grid's placemeent offset.
func (grid *Grid) SetOffset(x, y int) {
	grid.offsetX = x
	grid.offsetY = y
}

// Color returns the grid's line color.
func (grid *Grid) Color() color.Color {
	return grid.color
}

// SetColor sets the grid's line color.
func (grid *Grid) SetColor(c color.Color) {
	grid.color = c
}

// SetClickHandler sets the grid's click handler.
func (grid *Grid) SetClickHandler(cb func(x, y int)) {
	grid.clickHandler = cb
}

// SetHeldHandler sets the grid's held handler.
func (grid *Grid) SetHeldHandler(cb func(x, y int)) {
	grid.heldHandler = cb
}

func (grid *Grid) Update(ctx ifs.RunContext) {
	if grid.clickHandler == nil && grid.heldHandler == nil {
		return
	}
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		return
	}
	x, y := ebiten.CursorPosition()
	x -= grid.offsetX
	y -= grid.offsetY
	if x < 0 || y < 0 {
		return
	}
	x /= grid.cellWidth
	y /= grid.cellHeight
	if grid.clickHandler != nil && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		grid.clickHandler(x, y)
	}
	if grid.heldHandler != nil {
		grid.heldHandler(x, y)
	}
}

// Draw draws the grid to the screen.
func (grid *Grid) Draw(ctx ifs.DrawContext) {
	if grid.image == nil {
		return
	}
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(grid.offsetX), float64(grid.offsetY))
	ctx.Screen.DrawImage(grid.image, &opts)
}
