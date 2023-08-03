package states

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/id"
	"github.com/nfnt/resize"
)

type Data struct {
	archetypes      map[id.UUID]game.Archetype
	archetypeImages map[id.UUID]*ebiten.Image
	tiles           map[id.UUID]game.TileArchetype
	tileImages      map[id.UUID]*ebiten.Image
}

func NewData() *Data {
	return &Data{
		archetypes:      make(map[id.UUID]game.Archetype),
		archetypeImages: make(map[id.UUID]*ebiten.Image),
		tiles:           make(map[id.UUID]game.TileArchetype),
		tileImages:      make(map[id.UUID]*ebiten.Image),
	}
}

func (d *Data) Archetype(id id.UUID) game.Archetype {
	return d.archetypes[id]
}

func (d *Data) ArchetypeImage(id id.UUID) *ebiten.Image {
	return d.archetypeImages[id]
}

func (d *Data) EnsureImage(archetype game.Archetype, zoom float64) (*ebiten.Image, error) {
	if img, ok := d.archetypeImages[archetype.GetID()]; ok {
		return img, nil
	}

	switch a := archetype.(type) {
	case game.CharacterArchetype:
		if img, err := d.LoadImage("archetypes/"+a.Image, zoom); err == nil {
			d.archetypeImages[a.GetID()] = img
			return img, nil
		} else {
			return nil, err
		}
	case game.ItemArchetype:
		if img, err := d.LoadImage("archetypes/"+a.Image, zoom); err == nil {
			d.archetypeImages[a.GetID()] = img
			return img, nil
		} else {
			return nil, err
		}
	case game.WeaponArchetype:
		if img, err := d.LoadImage("archetypes/"+a.Image, zoom); err == nil {
			d.archetypeImages[a.GetID()] = img
			return img, nil
		} else {
			return nil, err
		}
	case game.ArmorArchetype:
		if img, err := d.LoadImage("archetypes/"+a.Image, zoom); err == nil {
			d.archetypeImages[a.GetID()] = img
			return img, nil
		} else {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown archetype type: %T", archetype)
	}
}

func (d *Data) LoadImage(src string, scale float64) (*ebiten.Image, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/"+src, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(resBody))
	if err != nil {
		return nil, err
	}

	// Resize the image to 2x until ebitenui has scaling built-in.
	img = resize.Resize(uint(float64(img.Bounds().Dx())*scale), uint(float64(img.Bounds().Dy())*scale), img, resize.NearestNeighbor)

	return ebiten.NewImageFromImage(img), nil
}
