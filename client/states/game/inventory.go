package game

import (
	"image"
	"image/color"
	"time"

	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/game"
)

type Inventory struct {
	inventory      *[]game.Object
	container      *widget.Container
	innerContainer *widget.Container
	cells          []inventoryCell
}

type inventoryCell struct {
	cell    *widget.Container
	tooltip *widget.ToolTip
	graphic *widget.Graphic
}

func (inv *Inventory) Init(container *widget.Container, ctx ifs.RunContext) {
	inv.container = container

	inv.innerContainer = widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 32})),
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			//Define number of columns in the grid
			widget.GridLayoutOpts.Columns(5),
			//Define how far apart the rows and columns should be
			widget.GridLayoutOpts.Spacing(2, 2),
			//Define how to stretch the rows and columns. Note it is required to
			//specify the Stretch for each row and column.
			widget.GridLayoutOpts.Stretch([]bool{true, true, true, true, true}, []bool{true, true, true, true, true}),
		)),
	)

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			tool := widget.NewTextToolTip("", ctx.UI.BodyCopyFace, color.White, eimage.NewNineSliceColor(color.NRGBA{R: 50, G: 50, B: 50, A: 255}))
			tool.Position = widget.TOOLTIP_POS_CURSOR_STICKY
			tool.Delay = time.Duration(time.Millisecond * 200)
			// Set to -1000 to hide it for the time being.
			tool.Offset = image.Pt(-1000, -1000)

			graphic := widget.NewGraphic(
				widget.GraphicOpts.Image(nil),
				widget.GraphicOpts.WidgetOpts(
					widget.WidgetOpts.LayoutData(
						widget.RowLayoutPositionCenter,
					),
				),
			)

			cell := widget.NewContainer(
				widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{128, 128, 128, 128})),
				widget.ContainerOpts.WidgetOpts(
					widget.WidgetOpts.LayoutData(widget.GridLayoutData{
						MaxWidth:  34,
						MaxHeight: 34,
					}),
					widget.WidgetOpts.ToolTip(tool),
				),
			)

			cell.AddChild(graphic)

			inv.cells = append(inv.cells, inventoryCell{
				cell:    cell,
				tooltip: tool,
				graphic: graphic,
			})

			inv.innerContainer.AddChild(cell)
		}
	}

	inv.container.AddChild(inv.innerContainer)
}

func (inv *Inventory) SetInventory(inventory *[]game.Object) {
	inv.inventory = inventory
}

func (inv *Inventory) Refresh() {
	// TODO: Refresh inventory
}

func (inv *Inventory) Draw(ctx ifs.DrawContext) {
	// TODO: Draw
}
