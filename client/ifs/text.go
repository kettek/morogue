package ifs

import (
	"image/color"

	"github.com/tinne26/etxt"
	"golang.org/x/image/font/sfnt"
)

type TextRenderer struct {
	etxt.Renderer

	outlineColor color.Color
	//
	savedFont         *sfnt.Font
	savedColor        color.Color
	savedSize         float64
	savedScale        float64
	savedBlendMode    etxt.BlendMode
	savedAlign        etxt.Align
	savedOutlineColor color.Color
}

func NewTextRenderer(r *etxt.Renderer) *TextRenderer {
	return &TextRenderer{
		Renderer: *r,
	}
}

func (t *TextRenderer) DrawWithOutline(target etxt.TargetImage, text string, x, y int) {
	c := t.GetColor()
	t.SetColor(t.outlineColor)
	{
		t.Draw(target, text, x-1, y)
		t.Draw(target, text, x+1, y)
		t.Draw(target, text, x, y-1)
		t.Draw(target, text, x, y+1)

		t.Draw(target, text, x-1, y+1)
		t.Draw(target, text, x+1, y-1)
		t.Draw(target, text, x-1, y-1)
		t.Draw(target, text, x+1, y+1)
	}
	t.SetColor(c)
	t.Draw(target, text, x, y)
}

func (t *TextRenderer) SetOutlineColor(c color.Color) {
	t.outlineColor = c
}

func (t *TextRenderer) GetOutlineColor() color.Color {
	return t.outlineColor
}

func (t *TextRenderer) Save() {
	t.savedFont = t.GetFont()
	t.savedColor = t.GetColor()
	t.savedSize = t.GetSize()
	t.savedScale = t.GetScale()
	t.savedBlendMode = t.GetBlendMode()
	t.savedAlign = t.GetAlign()
	t.savedOutlineColor = t.GetOutlineColor()
}

func (t *TextRenderer) Restore() {
	t.SetFont(t.savedFont)
	t.SetColor(t.savedColor)
	t.SetSize(t.savedSize)
	t.SetScale(t.savedScale)
	t.SetBlendMode(t.savedBlendMode)
	t.SetAlign(t.savedAlign)
	t.SetOutlineColor(t.outlineColor)
}
