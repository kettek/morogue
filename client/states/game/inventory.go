package game

import (
	"image"
	"image/color"
	"time"

	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/morogue/client/embed"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/id"
)

type Inventory struct {
	Data           Data
	container      *widget.Container
	innerContainer *widget.Container
	cells          []*inventoryCell
	ApplyItem      func(wid id.WID, apply bool)
	DropItem       func(wid id.WID)
	PickupItem     func(wid id.WID)
}

type inventoryCell struct {
	cell           *widget.Container
	tooltip        *widget.ToolTip
	tooltipContent *widget.Container
	graphic        *widget.Graphic
	indicator      *widget.Graphic
	WID            id.WID
	//
	applied bool
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
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MouseButtonPressedHandler(func(args *widget.WidgetMouseButtonPressedEventArgs) {
				ctx.Game.PreventMapInput = true
			}),
			widget.WidgetOpts.MouseButtonReleasedHandler(func(args *widget.WidgetMouseButtonReleasedEventArgs) {
				ctx.Game.PreventMapInput = false
			}),
		),
	)

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			clickCount := 0
			lastTime := time.Now()

			tooltipContent := widget.NewContainer(
				widget.ContainerOpts.BackgroundImage(ctx.UI.ItemInfoBackgroundImage),
				widget.ContainerOpts.AutoDisableChildren(),
				widget.ContainerOpts.Layout(widget.NewRowLayout(
					widget.RowLayoutOpts.Direction(widget.DirectionVertical),
					widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(8)),
				)),
			)

			tool := widget.NewToolTip(
				widget.ToolTipOpts.Content(tooltipContent),
				widget.ToolTipOpts.Delay(0),
				widget.ToolTipOpts.Offset(image.Point{-1000, -1000}),
				widget.ToolTipOpts.ContentOriginHorizontal(widget.TOOLTIP_ANCHOR_START),
				widget.ToolTipOpts.ContentOriginVertical(widget.TOOLTIP_ANCHOR_START),
			)
			tool.Position = widget.TOOLTIP_POS_CURSOR_STICKY

			graphic := widget.NewGraphic(
				widget.GraphicOpts.Image(nil),
				widget.GraphicOpts.WidgetOpts(
					widget.WidgetOpts.LayoutData(
						widget.RowLayoutPositionCenter,
					),
				),
			)

			indicatorOuter := widget.NewContainer(
				widget.ContainerOpts.Layout(widget.NewAnchorLayout(
					widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(2)),
				)),
			)
			indicator := widget.NewGraphic(
				widget.GraphicOpts.Image(nil),
				widget.GraphicOpts.WidgetOpts(
					widget.WidgetOpts.LayoutData(
						widget.AnchorLayoutData{
							HorizontalPosition: widget.AnchorLayoutPositionEnd,
							VerticalPosition:   widget.AnchorLayoutPositionStart,
						},
					),
				),
			)
			indicatorOuter.AddChild(indicator)

			invCell := &inventoryCell{}

			cell := widget.NewContainer(
				widget.ContainerOpts.BackgroundImage(cellBackgroundImage),
				widget.ContainerOpts.Layout(widget.NewStackedLayout()),
				widget.ContainerOpts.WidgetOpts(
					widget.WidgetOpts.LayoutData(widget.GridLayoutData{
						MaxWidth:           34,
						MaxHeight:          34,
						HorizontalPosition: widget.GridLayoutPositionEnd,
						VerticalPosition:   widget.GridLayoutPositionStart,
					}),
					widget.WidgetOpts.ToolTip(tool),
					widget.WidgetOpts.EnableDragAndDrop(
						widget.NewDragAndDrop(
							widget.DragAndDropOpts.ContentsCreater(makeDragWidget(ctx, dragContainer{cell: invCell, container: inv})),
							widget.DragAndDropOpts.MinDragStartDistance(8),
							widget.DragAndDropOpts.ContentsOriginVertical(widget.DND_ANCHOR_END),
							widget.DragAndDropOpts.ContentsOriginHorizontal(widget.DND_ANCHOR_END),
							widget.DragAndDropOpts.Offset(image.Point{16, 16}),
						),
					),
					widget.WidgetOpts.CanDrop(func(args *widget.DragAndDropDroppedEventArgs) bool {
						switch args.Data.(dragContainer).cell.(type) {
						case *inventoryCell:
							return true
						case *belowCell:
							return true
						}
						return false
					}),
					widget.WidgetOpts.Dropped(func(args *widget.DragAndDropDroppedEventArgs) {
						switch cell := args.Data.(dragContainer).cell.(type) {
						case *inventoryCell:
							// TODO: Reorganize inventory?
						case *belowCell:
							inv.PickupItem(cell.WID)
						}
					}),
					widget.WidgetOpts.MouseButtonReleasedHandler(func(args *widget.WidgetMouseButtonReleasedEventArgs) {
						if args.Inside {
							if args.Button == ebiten.MouseButtonLeft {
								if time.Since(lastTime) > 500*time.Millisecond {
									clickCount = 0
									lastTime = time.Now()
								}
								clickCount++
								if clickCount == 2 {
									inv.ApplyItem(invCell.WID, !invCell.applied)
									clickCount = 0
									return
								}
							} else if args.Button == ebiten.MouseButtonRight {
								if time.Since(lastTime) > 500*time.Millisecond {
									clickCount = 0
									lastTime = time.Now()
								}
								clickCount++
								if clickCount == 2 {
									inv.DropItem(invCell.WID)
									clickCount = 0
									return
								}
							}
						}

						if args.Inside && args.Button == ebiten.MouseButtonLeft && ebiten.IsKeyPressed(ebiten.KeyControl) {
							args.Widget.DragAndDrop.StartDrag()
						}
						if args.Button == ebiten.MouseButtonRight {
							args.Widget.DragAndDrop.StopDrag()
						}
					}),
				),
			)

			cell.AddChild(graphic)
			cell.AddChild(indicatorOuter)

			invCell.cell = cell
			invCell.tooltip = tool
			invCell.tooltipContent = tooltipContent
			invCell.graphic = graphic
			invCell.indicator = indicator

			inv.cells = append(inv.cells, invCell)

			inv.innerContainer.AddChild(cell)
		}
	}

	inv.container.AddChild(inv.innerContainer)
}

func (inv *Inventory) Refresh(ctx ifs.RunContext, objects game.Objects) {
	// TODO: Don't clear cells that have remained the same.
	// Clear old cells.
	for _, cell := range inv.cells {
		if cell.WID == 0 {
			continue
		}
		cell.tooltip.Offset = image.Pt(-1000, -1000)
		cell.graphic.Image = nil
		cell.indicator.Image = nil
		cell.tooltipContent.RemoveChildren()
		cell.applied = false
		cell.WID = 0
	}

	// Refresh it.
	for i, o := range objects {
		img := inv.Data.ArchetypeImage(o.GetArchetypeID())
		inv.cells[i].WID = o.GetWID()
		inv.cells[i].graphic.Image = img
		inv.cells[i].tooltip.Offset = image.Pt(2, 2)

		arch := inv.Data.Archetype(o.GetArchetypeID())
		inv.cells[i].tooltipContent.RemoveChildren()
		addObjectInfo(ctx, o, arch, inv.cells[i].tooltipContent)

		inv.cells[i].indicator.Image = nil
		switch o := o.(type) {
		case *game.Weapon:
			if o.Applied {
				inv.cells[i].indicator.Image = embed.IndicatorApplied
				inv.cells[i].applied = true
			}
		case *game.Armor:
			if o.Applied {
				inv.cells[i].indicator.Image = embed.IndicatorApplied
				inv.cells[i].applied = true
			}
		}
	}
}
