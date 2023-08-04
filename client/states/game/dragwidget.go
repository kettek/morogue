package game

import (
	"fmt"
	"image/color"

	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kettek/morogue/client/ifs"
)

// This is annoying
var cellBackgroundImage *eimage.NineSlice
var cellBackgroundHoverImage *eimage.NineSlice

func init() {
	cellBackgroundImage = eimage.NewNineSliceColor(color.NRGBA{64, 64, 64, 128})
	cellBackgroundHoverImage = eimage.NewNineSliceColor(color.NRGBA{128, 128, 128, 128})
}

type dragWidget struct {
	container      *widget.Container
	graphic        *widget.Graphic
	targetedWidget widget.HasWidget
	ctx            ifs.RunContext
	data           interface{}
}

func makeDragWidget(ctx ifs.RunContext, data interface{}) *dragWidget {
	return &dragWidget{
		ctx:  ctx,
		data: data,
	}
}

func (w *dragWidget) Create(parent widget.HasWidget) (*widget.Container, interface{}) {
	// For this example we do not need to recreate the Dragged element each time. We can re-use it.
	if w.container == nil {
		w.container = widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			)),
			//widget.ContainerOpts.BackgroundImage(cellBackgroundHoverImage),
		)
	}
	if w.graphic == nil {
		w.graphic = widget.NewGraphic(widget.GraphicOpts.Image(nil), widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})))
		w.container.AddChild(w.graphic)
	}

	if d, ok := w.data.(dragContainer); ok {
		switch cell := d.cell.(type) {
		case *inventoryCell:
			w.graphic.Image = cell.graphic.Image
		case *belowCell:
			w.graphic.Image = cell.graphic.Image
		}
	}

	// return the container to be dragged and any arbitrary data associated with this operation
	return w.container, w.data
}

func (w *dragWidget) Update(canDrop bool, targetWidget widget.HasWidget, dragData interface{}) {
	if canDrop {
		if targetWidget != nil {
			targetWidget.(*widget.Container).BackgroundImage = cellBackgroundHoverImage
			w.targetedWidget = targetWidget
		}
	} else {
		if w.targetedWidget != nil {
			w.targetedWidget.(*widget.Container).BackgroundImage = cellBackgroundImage
			w.targetedWidget = nil
		}
	}
}

func (w *dragWidget) EndDrag(dropped bool, sourceWidget widget.HasWidget, dragData interface{}) {
	if !dropped {
		if dragData, ok := dragData.(dragContainer); ok {
			switch cell := dragData.cell.(type) {
			case *inventoryCell:
				if container, ok := dragData.container.(*Inventory); ok {
					container.DropItem(cell.WID)
				}
			case *belowCell:
				if container, ok := dragData.container.(*Below); ok {
					container.PickupItem(cell.WID)
				}
			case *hotbarCell:
				fmt.Println("dropped hotbar cell nowhere, probably clear it", cell)
			}
		}
	}

	if w.targetedWidget != nil {
		w.targetedWidget.(*widget.Container).BackgroundImage = cellBackgroundImage
	}
}

type dragContainer struct {
	cell      interface{}
	container interface{}
}
