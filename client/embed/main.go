package embed

import (
	"embed"
	"image"
	_ "image/png"

	einput "github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font/sfnt"
)

//go:embed images fonts
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
var FontDefault *sfnt.Font

var IndicatorImage *ebiten.Image
var IndicatorApplied *ebiten.Image
var IndicatorCursed *ebiten.Image

var IconOffenseMelee *ebiten.Image
var IconOffenseRanged *ebiten.Image
var IconOffenseThrown *ebiten.Image
var IconOffenseUnarmed *ebiten.Image
var IconDefense *ebiten.Image
var IconHealth *ebiten.Image
var IconHunger *ebiten.Image

var PingLook *ebiten.Image
var PingTarget *ebiten.Image
var PingDanger *ebiten.Image
var PingBasic *ebiten.Image
var PingDefend *ebiten.Image

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

	//
	f, err = Assets.Open("images/indicators.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	IndicatorImage, _, _ = ebitenutil.NewImageFromReader(f)
	IndicatorApplied = ebiten.NewImageFromImage(IndicatorImage.SubImage(image.Rect(0, 0, 8, 8)))
	IndicatorCursed = ebiten.NewImageFromImage(IndicatorImage.SubImage(image.Rect(8, 0, 16, 8)))

	//
	f, err = Assets.Open("images/defense.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	IconDefense, _, _ = ebitenutil.NewImageFromReader(f)

	f, err = Assets.Open("images/offense-melee.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	IconOffenseMelee, _, _ = ebitenutil.NewImageFromReader(f)

	f, err = Assets.Open("images/offense-thrown.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	IconOffenseThrown, _, _ = ebitenutil.NewImageFromReader(f)

	f, err = Assets.Open("images/offense-ranged.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	IconOffenseRanged, _, _ = ebitenutil.NewImageFromReader(f)

	f, err = Assets.Open("images/offense-unarmed.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	IconOffenseUnarmed, _, _ = ebitenutil.NewImageFromReader(f)

	f, err = Assets.Open("images/health.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	IconHealth, _, _ = ebitenutil.NewImageFromReader(f)

	f, err = Assets.Open("images/hunger.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	IconHunger, _, _ = ebitenutil.NewImageFromReader(f)

	f, err = Assets.Open("images/ping-look.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	PingLook, _, _ = ebitenutil.NewImageFromReader(f)

	f, err = Assets.Open("images/ping-target.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	PingTarget, _, _ = ebitenutil.NewImageFromReader(f)

	f, err = Assets.Open("images/ping-danger.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	PingDanger, _, _ = ebitenutil.NewImageFromReader(f)

	f, err = Assets.Open("images/ping-basic.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	PingBasic, _, _ = ebitenutil.NewImageFromReader(f)

	f, err = Assets.Open("images/ping-defend.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	PingDefend, _, _ = ebitenutil.NewImageFromReader(f)

	//
	b, err := Assets.ReadFile("fonts/x12y16pxMaruMonica.ttf")
	if err != nil {
		panic(err)
	}
	FontDefault, err = sfnt.Parse(b)
	if err != nil {
		panic(err)
	}
}
