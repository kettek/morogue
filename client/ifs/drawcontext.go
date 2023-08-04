package ifs

import (
	"image/color"

	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// DrawContext contains various structures that are useful during Draw calls.
type DrawContext struct {
	Txt    *TextRenderer
	Screen *ebiten.Image
	UI     *DrawContextUI
	Game   *GameContext
}

// DrawContextUI is a structure containing the UI styling and data.
type DrawContextUI struct {
	Width, Height int
	//
	HeadlineFace     font.Face
	BodyCopyFace     font.Face
	ButtonImage      *widget.ButtonImage
	ButtonTextColor  *widget.ButtonTextColor
	ButtonPadding    widget.Insets
	TextInputColor   *widget.TextInputColor
	TextInputImage   *widget.TextInputImage
	TextInputPadding widget.Insets
	//
	ItemInfoBackgroundImage *eimage.NineSlice
}

// Init sets up the necessary data structuers, such as fonts, styling, images, etc.
func (ui *DrawContextUI) Init(txt *TextRenderer) {
	ui.HeadlineFace, _ = opentype.NewFace(txt.Renderer.GetFont(), &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	ui.BodyCopyFace, _ = opentype.NewFace(txt.Renderer.GetFont(), &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	ui.ButtonImage = &widget.ButtonImage{
		Idle:    eimage.NewNineSliceColor(color.NRGBA{R: 40, G: 30, B: 40, A: 255}),
		Hover:   eimage.NewNineSliceColor(color.NRGBA{R: 50, G: 40, B: 50, A: 255}),
		Pressed: eimage.NewNineSliceColor(color.NRGBA{R: 60, G: 50, B: 60, A: 255}),
	}

	ui.ButtonPadding = widget.Insets{
		Left:   30,
		Right:  30,
		Top:    5,
		Bottom: 5,
	}

	ui.ButtonTextColor = &widget.ButtonTextColor{
		Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
	}

	ui.TextInputImage = &widget.TextInputImage{
		Idle:     eimage.NewNineSliceColor(color.NRGBA{R: 60, G: 50, B: 60, A: 255}),
		Disabled: eimage.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
	}

	ui.TextInputPadding = widget.NewInsetsSimple(5)

	ui.TextInputColor = &widget.TextInputColor{
		Idle:          color.NRGBA{254, 255, 255, 255},
		Disabled:      color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		Caret:         color.NRGBA{254, 255, 255, 255},
		DisabledCaret: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
	}

	ui.ItemInfoBackgroundImage = eimage.NewNineSliceColor(color.NRGBA{R: 20, G: 20, B: 20, A: 255})

}
