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

var Icon image.Image
var CursorDefault *ebiten.Image
var CursorDefaultTooltip *ebiten.Image
var CursorPointer *ebiten.Image
var CursorPointerTooltip *ebiten.Image
var CursorText *ebiten.Image
var CursorTextTooltip *ebiten.Image
var CursorDelete *ebiten.Image
var CursorDeleteTooltip *ebiten.Image
var CursorMove *ebiten.Image

func Setup() {
	f, err := Assets.Open("images/icon.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	Icon, _, err = image.Decode(f)
	if err != nil {
		panic(err)
	}

	f, err = Assets.Open("images/cursors.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	i, _, _ := ebitenutil.NewImageFromReader(f)
	CursorDefault = ebiten.NewImageFromImage(i.SubImage(image.Rect(0, 0, 16, 16)))
	CursorDefaultTooltip = ebiten.NewImageFromImage(i.SubImage(image.Rect(0, 16, 16, 32)))
	CursorPointer = ebiten.NewImageFromImage(i.SubImage(image.Rect(16, 0, 32, 16)))
	CursorPointerTooltip = ebiten.NewImageFromImage(i.SubImage(image.Rect(16, 16, 32, 32)))
	CursorText = ebiten.NewImageFromImage(i.SubImage(image.Rect(32, 0, 48, 16)))
	CursorTextTooltip = ebiten.NewImageFromImage(i.SubImage(image.Rect(32, 16, 48, 32)))
	CursorDelete = ebiten.NewImageFromImage(i.SubImage(image.Rect(48, 0, 64, 16)))
	CursorDeleteTooltip = ebiten.NewImageFromImage(i.SubImage(image.Rect(48, 16, 64, 32)))
	CursorMove = ebiten.NewImageFromImage(i.SubImage(image.Rect(64, 0, 80, 16)))

	//
	einput.SetCursorImage(einput.CURSOR_DEFAULT, CursorDefault)
	einput.SetCursorImage("default-tooltip", CursorDefaultTooltip)
	einput.SetCursorImage("interactive", CursorPointer)
	einput.SetCursorImage("interactive-tooltip", CursorPointerTooltip)
	einput.SetCursorImage("text", CursorText)
	einput.SetCursorImage("text-tooltip", CursorTextTooltip)
	einput.SetCursorImage("delete", CursorDelete)
	einput.SetCursorImage("delete-tooltip", CursorDeleteTooltip)
	einput.SetCursorImage("move", CursorMove)
}
