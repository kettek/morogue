package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kettek/morogue/client/ifs"
)

type Grid struct {
	offsetX, offsetY      int
	width, height         int
	cellWidth, cellHeight int
	color                 color.Color
	image                 *ebiten.Image
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

func (grid *Grid) Size() (int, int) {
	return grid.width, grid.height
}

func (grid *Grid) SetSize(w, h int) {
	grid.width = w
	grid.height = h
	grid.makeImage()
}

func (grid *Grid) CellSize() (int, int) {
	return grid.cellWidth, grid.cellHeight
}

func (grid *Grid) SetCellSize(w, h int) {
	grid.cellWidth = w
	grid.cellHeight = h
	grid.makeImage()
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
	if grid.image == nil {
		return
	}
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(grid.offsetX), float64(grid.offsetY))
	ctx.Screen.DrawImage(grid.image, &opts)
}
