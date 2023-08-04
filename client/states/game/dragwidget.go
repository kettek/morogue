package game

import (
	"fmt"
	"image/color"

	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kettek/morogue/client/ifs"
)

type dragWidget struct {
	container      *widget.Container
	graphic        *widget.Graphic
	text           *widget.Text
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
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
			widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{0, 200, 100, 255})),
		)

		w.text = widget.NewText(widget.TextOpts.Text("Cannot Drop", w.ctx.UI.HeadlineFace, color.White), widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})))

		w.container.AddChild(w.text)
	}
	// return the container to be dragged and any arbitrary data associated with this operation
	return w.container, w.data
}

func (w *dragWidget) Update(canDrop bool, targetWidget widget.HasWidget, dragData interface{}) {
	if canDrop {
		w.text.Label = "okay"
		if targetWidget != nil {
			targetWidget.(*widget.Container).BackgroundImage = eimage.NewNineSliceColor(color.NRGBA{128, 128, 255, 128})
			w.targetedWidget = targetWidget
		}
	} else {
		w.text.Label = "nokay"
		if w.targetedWidget != nil {
			w.targetedWidget.(*widget.Container).BackgroundImage = eimage.NewNineSliceColor(color.NRGBA{128, 128, 128, 128})
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
		w.targetedWidget.(*widget.Container).BackgroundImage = eimage.NewNineSliceColor(color.NRGBA{128, 128, 128, 128})
	}
}

type dragContainer struct {
	cell      interface{}
	container interface{}
}
