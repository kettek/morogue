package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kettek/morogue/client/ifs"
)

type Grid struct {
	offsetX, offsetY      int
	width, height         int
	cellWidth, cellHeight int
	color                 color.Color
}

func (grid *Grid) Size() (int, int) {
	return grid.width, grid.height
}

func (grid *Grid) SetSize(w, h int) {
	grid.width = w
	grid.height = h
}

func (grid *Grid) CellSize() (int, int) {
	return grid.cellWidth, grid.cellHeight
}

func (grid *Grid) SetCellSize(w, h int) {
	grid.cellWidth = w
	grid.cellHeight = h
}

func (grid *Grid) Offset() (x, y int) {
	return grid.offsetX, grid.offsetY
}

func (grid *Grid) SetOffset(x, y int) {
	grid.offsetX = x
	grid.offsetY = y
}

func (grid *Grid) Color() color.Color {
	return grid.color
}

func (grid *Grid) SetColor(c color.Color) {
	grid.color = c
}

func (grid *Grid) Draw(ctx ifs.DrawContext) {
	cx := grid.width / grid.cellWidth
	cy := grid.height / grid.cellHeight
	for x := 1; x != cx; x++ {
		vector.StrokeLine(ctx.Screen, float32(x*grid.cellWidth+grid.offsetX), float32(grid.offsetY), float32(x*grid.cellWidth+grid.offsetX), float32(grid.height+grid.offsetY), 1, grid.color, false)
	}
	for y := 1; y != cy; y++ {
		vector.StrokeLine(ctx.Screen, float32(grid.offsetX), float32(y*grid.cellHeight+grid.offsetY), float32(grid.width+grid.offsetX), float32(y*grid.cellHeight+grid.offsetY), 1, grid.color, false)
	}
}
