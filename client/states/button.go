package states

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

func buttonImages() (*widget.ButtonImage, error) {
	idle := image.NewNineSliceColor(color.NRGBA{R: 180, G: 170, B: 180, A: 255})

	hover := image.NewNineSliceColor(color.NRGBA{R: 150, G: 130, B: 150, A: 255})

	pressed := image.NewNineSliceColor(color.NRGBA{R: 120, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}
