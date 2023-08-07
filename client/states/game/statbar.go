package game

import (
	"fmt"
	"image/color"

	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/game"
)

type Statbar struct {
	container      *widget.Container
	innerContainer *widget.Container
}

func (hb *Statbar) Init(container *widget.Container, ctx ifs.RunContext) {
	hb.container = container

	hb.innerContainer = widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{255, 0x1a, 0x22, 255})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal:  true,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(30, 30),
			widget.WidgetOpts.MouseButtonPressedHandler(func(args *widget.WidgetMouseButtonPressedEventArgs) {
				ctx.Game.PreventMapInput = true
			}),
			widget.WidgetOpts.MouseButtonReleasedHandler(func(args *widget.WidgetMouseButtonReleasedEventArgs) {
				ctx.Game.PreventMapInput = false
			}),
		),
	)

	hb.container.AddChild(hb.innerContainer)
}

func (hb *Statbar) Refresh(ctx ifs.RunContext, c *game.Character, a game.CharacterArchetype) {
	if c.Archetype != nil {
		c.CacheDamages()
		for _, d := range c.Damages {
			fmt.Printf("%+v\n", d)
		}
	}
}
