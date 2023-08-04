package game

import (
	"fmt"
	"image"
	"image/color"
	"time"

	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/game"
)

type Hotbar struct {
	hbentory       *[]game.Object
	container      *widget.Container
	innerContainer *widget.Container
	cells          []*hotbarCell
	binds          *Binds
	prefix         string
}

type hotbarCell struct {
	cell    *widget.Container
	tooltip *widget.ToolTip
	graphic *widget.Graphic
}

func (hb *Hotbar) Init(container *widget.Container, ctx ifs.RunContext, binds *Binds) {
	hb.container = container
	hb.binds = binds
	hb.prefix = "default"

	for i := 1; i <= 10; i++ {
		if i == 10 {
			hb.binds.SetActionKeys(Action(fmt.Sprintf("hotbar-%s-%d", hb.prefix, 0)), []ebiten.Key{ebiten.Key0})
		} else {
			hb.binds.SetActionKeys(Action(fmt.Sprintf("hotbar-%s-%d", hb.prefix, i)), []ebiten.Key{ebiten.Key0 + ebiten.Key(i)})
		}
	}

	hb.innerContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(10),
			widget.GridLayoutOpts.Spacing(2, 2),
			widget.GridLayoutOpts.Stretch([]bool{true, true, true, true, true, true, true, true, true, true}, nil),
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

	for i := 0; i < 10; i++ {
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

		hbCell := &hotbarCell{}

		cell := widget.NewContainer(
			widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{64, 64, 64, 128})),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.GridLayoutData{
					MaxWidth:  68,
					MaxHeight: 68,
				}),
				widget.WidgetOpts.ToolTip(tool),
				widget.WidgetOpts.EnableDragAndDrop(
					widget.NewDragAndDrop(
						widget.DragAndDropOpts.ContentsCreater(makeDragWidget(ctx, dragContainer{cell: hbCell, container: hb})),
						widget.DragAndDropOpts.MinDragStartDistance(8),
						widget.DragAndDropOpts.ContentsOriginVertical(widget.DND_ANCHOR_END),
						widget.DragAndDropOpts.ContentsOriginHorizontal(widget.DND_ANCHOR_END),
						widget.DragAndDropOpts.Offset(image.Point{-5, -5}),
					),
				),
				widget.WidgetOpts.CanDrop(func(args *widget.DragAndDropDroppedEventArgs) bool {
					switch args.Data.(dragContainer).cell.(type) {
					case *inventoryCell:
						return true
					case *hotbarCell:
						return true
					}
					return false
				}),
				widget.WidgetOpts.Dropped(func(args *widget.DragAndDropDroppedEventArgs) {
					switch args.Data.(dragContainer).cell.(type) {
					case *inventoryCell:
						// TODO: Assign item to slot.
					case *hotbarCell:
						// TODO: Swap hotbar entries.
					}
				}),
				widget.WidgetOpts.MouseButtonReleasedHandler(func(args *widget.WidgetMouseButtonReleasedEventArgs) {
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

		hbCell.cell = cell
		hbCell.tooltip = tool
		hbCell.graphic = graphic

		hb.cells = append(hb.cells, hbCell)

		hb.innerContainer.AddChild(cell)
	}

	hb.container.AddChild(hb.innerContainer)
}

func (hb *Hotbar) Update(ctx ifs.RunContext) {
	for i := 1; i <= 10; i++ {
		var action Action
		if i == 10 {
			action = Action(fmt.Sprintf("hotbar-%s-%d", hb.prefix, 0))
		} else {
			action = Action(fmt.Sprintf("hotbar-%s-%d", hb.prefix, i))
		}

		if hb.binds.IsActionHeld(action) == 0 {
			hb.cells[i-1].cell.BackgroundImage = eimage.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 128})
		} else {
			hb.cells[i-1].cell.BackgroundImage = eimage.NewNineSliceColor(color.NRGBA{128, 128, 128, 128})
		}
	}
}

func (hb *Hotbar) Refresh() {
	// TODO: ???
}
