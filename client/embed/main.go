package embed

import (
	"embed"
	"image"
	_ "image/png"

	einput "github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed images
var Assets embed.FS

var CursorDefault *ebiten.Image
var CursorPointer *ebiten.Image
var CursorText *ebiten.Image

func Setup() {
	f, err := Assets.Open("images/cursors.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	i, _, _ := ebitenutil.NewImageFromReader(f)
	CursorDefault = ebiten.NewImageFromImage(i.SubImage(image.Rect(0, 0, 16, 16)))
	CursorPointer = ebiten.NewImageFromImage(i.SubImage(image.Rect(16, 0, 32, 16)))
	CursorText = ebiten.NewImageFromImage(i.SubImage(image.Rect(32, 0, 48, 16)))

	//
	einput.SetCursorImage(einput.CURSOR_DEFAULT, CursorDefault)
	einput.SetCursorImage("interactive", CursorPointer)
	einput.SetCursorImage("text", CursorText)
}
