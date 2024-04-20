package game

import (
	"fmt"
	"image/color"

	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/morogue/client/embed"
	"github.com/kettek/morogue/client/ifs"
	"github.com/kettek/morogue/game"
)

type Statbar struct {
	container        *widget.Container
	innerContainer   *widget.Container
	armorsContainer  *widget.Container
	damagesContainer *widget.Container
	healthContainer  *widget.Container
	hungerContainer  *widget.Container
}

func (hb *Statbar) Init(container *widget.Container, ctx ifs.RunContext) {
	hb.container = container

	hb.innerContainer = widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{0, 0, 0, 200})),
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

	hb.armorsContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal:  false,
				HorizontalPosition: widget.AnchorLayoutPositionStart,
			}),
		),
	)
	hb.innerContainer.AddChild(hb.armorsContainer)

	hb.damagesContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal:  false,
				HorizontalPosition: widget.AnchorLayoutPositionStart,
			}),
		),
	)
	hb.innerContainer.AddChild(hb.damagesContainer)

	hb.healthContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal:  false,
				HorizontalPosition: widget.AnchorLayoutPositionStart,
			}),
		),
	)
	hb.innerContainer.AddChild(hb.healthContainer)

	hb.hungerContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal:  false,
				HorizontalPosition: widget.AnchorLayoutPositionStart,
			}),
		),
	)
	hb.innerContainer.AddChild(hb.hungerContainer)

	hb.container.AddChild(hb.innerContainer)
}

func (hb *Statbar) Refresh(ctx ifs.RunContext, c *game.Character, a game.CharacterArchetype) {
	if c.Archetype != nil {
		{
			hb.armorsContainer.RemoveChildren()
			c.Hurtable.CalculateArmorFromCharacter(c)

			container := widget.NewContainer(
				widget.ContainerOpts.Layout(widget.NewRowLayout(
					widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
					widget.RowLayoutOpts.Padding(widget.Insets{Left: 8, Right: 8}),
				)),
			)

			value := widget.NewText(widget.TextOpts.ProcessBBCode(true), widget.TextOpts.Text(c.Hurtable.ArmorRangeString(), ctx.UI.BodyCopyFace, game.ArmorTypeNone.Color()))

			var icon *ebiten.Image

			icon = embed.IconDefense
			// TODO: Properly show armor type...

			graphic := widget.NewGraphic(
				widget.GraphicOpts.Image(icon),
				widget.GraphicOpts.WidgetOpts(
					widget.WidgetOpts.LayoutData(widget.RowLayoutData{
						Position: widget.RowLayoutPositionCenter,
					}),
				),
			)

			container.AddChild(value)
			container.AddChild(graphic)
			hb.armorsContainer.AddChild(container)
		}

		hb.damagesContainer.RemoveChildren()
		c.Damager.CalculateFromCharacter(c)
		for _, d := range c.Damages {
			weaponColor := game.WeaponTypeNone.Color()

			container := widget.NewContainer(
				widget.ContainerOpts.Layout(widget.NewRowLayout(
					widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
					widget.RowLayoutOpts.Padding(widget.Insets{Left: 8, Right: 8}),
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

			var icon *ebiten.Image

			switch d.Weapon {
			case game.WeaponTypeMelee:
				icon = embed.IconOffenseMelee
			case game.WeaponTypeRange:
				icon = embed.IconOffenseRanged
			case game.WeaponTypeThrown:
				icon = embed.IconOffenseThrown
			case game.WeaponTypeUnarmed:
				icon = embed.IconOffenseUnarmed
			}

			graphic := widget.NewGraphic(
				widget.GraphicOpts.Image(icon),
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

		{
			hb.healthContainer.RemoveChildren()
			c.Hurtable.CalculateFromCharacter(c)
			container := widget.NewContainer(
				widget.ContainerOpts.Layout(widget.NewRowLayout(
					widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
					widget.RowLayoutOpts.Padding(widget.Insets{Left: 8, Right: 8}),
				)),
			)
			value := widget.NewText(widget.TextOpts.ProcessBBCode(true), widget.TextOpts.Text(fmt.Sprintf("%d/%d", c.Hurtable.Health, c.Hurtable.MaxHealth), ctx.UI.BodyCopyFace, color.RGBA{255, 32, 32, 255}))

			graphic := widget.NewGraphic(
				widget.GraphicOpts.Image(embed.IconHealth),
				widget.GraphicOpts.WidgetOpts(
					widget.WidgetOpts.LayoutData(widget.RowLayoutData{
						Position: widget.RowLayoutPositionCenter,
					}),
				),
			)

			container.AddChild(value)
			container.AddChild(graphic)

			hb.healthContainer.AddChild(container)
		}

		{
			hb.hungerContainer.RemoveChildren()
			// TODO: Implement hunger
			container := widget.NewContainer(
				widget.ContainerOpts.Layout(widget.NewRowLayout(
					widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
					widget.RowLayoutOpts.Padding(widget.Insets{Left: 8, Right: 8}),
				)),
			)
			value := widget.NewText(widget.TextOpts.ProcessBBCode(true), widget.TextOpts.Text("1000", ctx.UI.BodyCopyFace, color.RGBA{255, 255, 32, 255}))

			graphic := widget.NewGraphic(
				widget.GraphicOpts.Image(embed.IconHunger),
				widget.GraphicOpts.WidgetOpts(
					widget.WidgetOpts.LayoutData(widget.RowLayoutData{
						Position: widget.RowLayoutPositionCenter,
					}),
				),
			)

			container.AddChild(value)
			container.AddChild(graphic)

			hb.hungerContainer.AddChild(container)
		}
	}
}
