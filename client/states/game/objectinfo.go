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

func makeDescription(ctx ifs.RunContext, txt string) *widget.TextArea {
	return widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position:  widget.RowLayoutPositionCenter,
					MaxWidth:  400,
					MaxHeight: 300,
				}),
				widget.WidgetOpts.MinSize(300, 30),
			),
		),
		widget.TextAreaOpts.ControlWidgetSpacing(2),
		widget.TextAreaOpts.ProcessBBCode(true),
		widget.TextAreaOpts.FontColor(color.NRGBA{128, 128, 128, 255}),
		widget.TextAreaOpts.FontFace(ctx.UI.BodyCopyFace),
		widget.TextAreaOpts.Text(txt),
		//This sets the background images for the scroll container
		widget.TextAreaOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle: ctx.UI.ItemInfoBackgroundImage,
				Mask: ctx.UI.ItemInfoBackgroundImage,
			}),
		),
		//This sets the eimages to use for the sliders
		widget.TextAreaOpts.SliderOpts(
			widget.SliderOpts.Images(
				// Set the track eimages
				&widget.SliderTrackImage{
					Idle:  eimage.NewNineSliceColor(color.NRGBA{200, 200, 200, 255}),
					Hover: eimage.NewNineSliceColor(color.NRGBA{200, 200, 200, 255}),
				},
				// Set the handle eimages
				&widget.ButtonImage{
					Idle:    eimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
					Hover:   eimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
					Pressed: eimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				},
			),
		),
	)
}

func addObjectInfo(ctx ifs.RunContext, object game.Object, arch game.Archetype, container *widget.Container) {
	switch a := arch.(type) {
	case game.WeaponArchetype:

		title := widget.NewText(widget.TextOpts.ProcessBBCode(true), widget.TextOpts.Text(fmt.Sprintf("%s", a.Title), ctx.UI.BodyCopyFace, color.White))
		values := widget.NewText(widget.TextOpts.ProcessBBCode(true), widget.TextOpts.Text(fmt.Sprintf("%s %s", a.RangeString(), a.WeaponType), ctx.UI.BodyCopyFace, a.WeaponType.Color()))
		slots := widget.NewText(widget.TextOpts.ProcessBBCode(true), widget.TextOpts.Text(fmt.Sprintf("%s", a.Slots.String()), ctx.UI.BodyCopyFace, color.NRGBA{R: 200, G: 200, B: 200, A: 255}))
		desc := makeDescription(ctx, a.Description)

		weaponLine := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Direction(widget.DirectionHorizontal))),
		)

		var icon *ebiten.Image

		switch a.WeaponType {
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
		weaponLine.AddChild(values)
		weaponLine.AddChild(graphic)

		container.AddChild(title)
		container.AddChild(slots)
		container.AddChild(weaponLine)
		container.AddChild(desc)
	case game.ArmorArchetype:
		title := widget.NewText(widget.TextOpts.ProcessBBCode(true), widget.TextOpts.Text(fmt.Sprintf("%s", a.Title), ctx.UI.BodyCopyFace, color.White))
		values := widget.NewText(widget.TextOpts.ProcessBBCode(true), widget.TextOpts.Text(fmt.Sprintf("%s %s", a.RangeString(), a.ArmorType), ctx.UI.BodyCopyFace, a.ArmorType.Color()))
		slots := widget.NewText(widget.TextOpts.ProcessBBCode(true), widget.TextOpts.Text(fmt.Sprintf("%s", a.Slots.String()), ctx.UI.BodyCopyFace, color.NRGBA{R: 200, G: 200, B: 200, A: 255}))
		desc := makeDescription(ctx, a.Description)

		armorLine := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Direction(widget.DirectionHorizontal))),
		)
		graphic := widget.NewGraphic(
			widget.GraphicOpts.Image(embed.IconDefense),
			widget.GraphicOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
			),
		)
		armorLine.AddChild(values)
		armorLine.AddChild(graphic)

		container.AddChild(title)
		container.AddChild(slots)
		container.AddChild(armorLine)
		container.AddChild(desc)
	case game.ItemArchetype:
		// TODO
	}

}
