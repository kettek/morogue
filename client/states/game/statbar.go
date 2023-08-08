package game

import (
	"image/color"

	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kettek/morogue/client/embed"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/game"
)

type Statbar struct {
	container        *widget.Container
	innerContainer   *widget.Container
	damagesContainer *widget.Container
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
			widget.WidgetOpts.MinSize(400, 20),
			widget.WidgetOpts.MouseButtonPressedHandler(func(args *widget.WidgetMouseButtonPressedEventArgs) {
				ctx.Game.PreventMapInput = true
			}),
			widget.WidgetOpts.MouseButtonReleasedHandler(func(args *widget.WidgetMouseButtonReleasedEventArgs) {
				ctx.Game.PreventMapInput = false
			}),
		),
	)

	hb.damagesContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal:  true,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
			}),
		),
	)
	hb.innerContainer.AddChild(hb.damagesContainer)

	hb.container.AddChild(hb.innerContainer)
}

func (hb *Statbar) Refresh(ctx ifs.RunContext, c *game.Character, a game.CharacterArchetype) {
	if c.Archetype != nil {
		hb.damagesContainer.RemoveChildren()
		c.CacheDamages()
		for _, d := range c.Damages {
			weaponColor := game.WeaponTypeNone.Color()

			container := widget.NewContainer(
				widget.ContainerOpts.Layout(widget.NewRowLayout(
					widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
				)),
			)

			if obj := c.Inventory.ObjectByWID(d.Source); obj != nil {
				if obj.GetArchetype() == nil {
					continue
				}
				if w, ok := obj.(*game.Weapon); ok {
					weaponColor = w.Archetype.(game.WeaponArchetype).WeaponType.Color()
				}
			}

			values := widget.NewText(widget.TextOpts.ProcessBBCode(true), widget.TextOpts.Text(d.RangeString(), ctx.UI.BodyCopyFace, weaponColor))

			graphic := widget.NewGraphic(
				widget.GraphicOpts.Image(embed.IconOffense),
				widget.GraphicOpts.WidgetOpts(
					widget.WidgetOpts.LayoutData(widget.RowLayoutData{
						Position: widget.RowLayoutPositionCenter,
					}),
				),
			)

			container.AddChild(values)
			container.AddChild(graphic)
			hb.damagesContainer.AddChild(container)
		}
	}
}
