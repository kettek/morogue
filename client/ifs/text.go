package ifs

import (
	"github.com/tinne26/etxt"
)

type TextRenderer struct {
	etxt.Renderer
}

func NewTextRenderer(r *etxt.Renderer) *TextRenderer {
	return &TextRenderer{*r}
}
