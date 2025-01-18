package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/gen"
	"github.com/kettek/morogue/id"
)

// Places is a slice of our generate-able places.
type Places []gen.Place

// ByID returns a place by its UUID.
func (p Places) ByID(uid id.UUID) (gen.Place, error) {
	for _, place := range p {
		if place.ID == uid {
			return place, nil
		}
	}
	return gen.Place{}, errors.New("no such place")
}

// Fixtures is a slice of our generate-able fixtures.
type Fixtures []gen.Fixture

// ByID returns a fixture by its UUID.
func (f Fixtures) ByID(uid id.UUID) (gen.Fixture, error) {
	for _, fixture := range f {
		if fixture.ID == uid {
			return fixture, nil
		}
	}
	return gen.Fixture{}, ErrNoSuchFixture
}

// Data contains our archetypes, places, and fixtures.
type Data struct {
	Archetypes []game.Archetype
	Places     Places
	Fixtures   Fixtures
}

func (d *Data) hasArchetype(uuid id.UUID) bool {
	for _, a := range d.Archetypes {
		if a.GetID() == uuid {
			return true
		}
	}
	return false
}

// Archetype returns an archetype by its UUID.
func (d *Data) Archetype(uuid id.UUID) game.Archetype {
	for _, a := range d.Archetypes {
		if a.GetID() == uuid {
			return a
		}
	}
	return nil
}

// Tile returns a TileArchetype by its UUID.
func (d *Data) Tile(uuid id.UUID) (game.TileArchetype, error) {
	for _, t := range d.TileArchetypes() {
		if t.ID == uuid {
			return t, nil
		}
	}
	return game.TileArchetype{}, ErrNoSuchTile
}

// TileArchetypes returns a slice of all TileArchetypes.
func (d *Data) TileArchetypes() []game.TileArchetype {
	var archetypes []game.TileArchetype
	for _, a := range d.Archetypes {
		if t, ok := a.(game.TileArchetype); ok {
			archetypes = append(archetypes, t)
		}
	}
	return archetypes
}

// CharacterArchetypes returns a slice of all CharacterArchetypes.
func (d *Data) CharacterArchetypes() []game.CharacterArchetype {
	var archetypes []game.CharacterArchetype
	for _, a := range d.Archetypes {
		if c, ok := a.(game.CharacterArchetype); ok {
			archetypes = append(archetypes, c)
		}
	}
	return archetypes
}

// ItemArchetypes returns a slice of all ItemArchetypes.
func (d *Data) ItemArchetypes() []game.ItemArchetype {
	var archetypes []game.ItemArchetype
	for _, a := range d.Archetypes {
		if i, ok := a.(game.ItemArchetype); ok {
			archetypes = append(archetypes, i)
		}
	}
	return archetypes
}

// WeaponArchetypes returns a slice of all WeaponArchetypes.
func (d *Data) WeaponArchetypes() []game.WeaponArchetype {
	var archetypes []game.WeaponArchetype
	for _, a := range d.Archetypes {
		if w, ok := a.(game.WeaponArchetype); ok {
			archetypes = append(archetypes, w)
		}
	}
	return archetypes
}

// ArmorArchetypes returns a slice of all ArmorArchetypes.
func (d *Data) ArmorArchetypes() []game.ArmorArchetype {
	var archetypes []game.ArmorArchetype
	for _, a := range d.Archetypes {
		if a, ok := a.(game.ArmorArchetype); ok {
			archetypes = append(archetypes, a)
		}
	}
	return archetypes
}

// FoodArchetypes returns a slice of all FoodArchetypes.
func (d *Data) FoodArchetypes() []game.FoodArchetype {
	var archetypes []game.FoodArchetype
	for _, a := range d.Archetypes {
		if f, ok := a.(game.FoodArchetype); ok {
			archetypes = append(archetypes, f)
		}
	}
	return archetypes
}

// LoadArchetypes loads all archetypes from the archetypes directory.
func (d *Data) LoadArchetypes() error {
	var iterate func(string, string) error

	iterate = func(fulldir string, partialdir string) error {
		entries, err := os.ReadDir(fulldir)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				if err := iterate(filepath.Join(fulldir, entry.Name()), filepath.Join(partialdir, entry.Name())); err != nil {
					log.Println(err)
				}
			} else {
				fullpath := filepath.Join(fulldir, entry.Name())
				if strings.HasSuffix(entry.Name(), ".json") {
					bytes, err := os.ReadFile(fullpath)
					if err != nil {
						log.Println(err)
						continue
					}
					if a, err := game.DecodeArchetype(bytes, partialdir); err != nil {
						log.Println(errors.Join(fmt.Errorf("failed to decode archetype %s", filepath.Join(fullpath, entry.Name())), err))
					} else {
						d.Archetypes = append(d.Archetypes, a)
					}
				}
			}
		}
		return nil
	}

	iterate("archetypes", "")

	return nil
}

// LoadPlaces loads all places from the places directory.
func (d *Data) LoadPlaces() error {
	var iterate func(string, string) error

	iterate = func(fulldir string, partialdir string) error {
		entries, err := os.ReadDir(fulldir)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				if err := iterate(filepath.Join(fulldir, entry.Name()), filepath.Join(partialdir, entry.Name())); err != nil {
					log.Println(err)
				}
			} else {
				fullpath := filepath.Join(fulldir, entry.Name())
				if strings.HasSuffix(entry.Name(), ".json") {
					bytes, err := os.ReadFile(fullpath)
					if err != nil {
						log.Println(err)
						continue
					}
					var p gen.Place
					if err := json.Unmarshal(bytes, &p); err != nil {
						log.Println(errors.Join(fmt.Errorf("failed to decode place %s", fullpath), err))
					} else {
						d.Places = append(d.Places, p)
					}
				}
			}
		}
		return nil
	}

	iterate("places", "")

	return nil
}

// LoadFixtures loads all fixtures from the fixtures directory.
func (d *Data) LoadFixtures() error {
	var iterate func(string, string) error

	iterate = func(fulldir string, partialdir string) error {
		entries, err := os.ReadDir(fulldir)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				if err := iterate(filepath.Join(fulldir, entry.Name()), filepath.Join(partialdir, entry.Name())); err != nil {
					log.Println(err)
				}
			} else {
				fullpath := filepath.Join(fulldir, entry.Name())
				if strings.HasSuffix(entry.Name(), ".json") {
					bytes, err := os.ReadFile(fullpath)
					if err != nil {
						log.Println(err)
						continue
					}
					var f gen.Fixture
					if err := json.Unmarshal(bytes, &f); err != nil {
						log.Println(errors.Join(fmt.Errorf("failed to decode place %s", fullpath), err))
					} else {
						d.Fixtures = append(d.Fixtures, f)
					}
				}
			}
		}
		return nil
	}

	iterate("fixtures", "")

	return nil
}

// Error types, yo.
var (
	ErrNoSuchFixture = errors.New(lc.T("no such fixture"))
	ErrNoSuchTile    = errors.New(lc.T("no such tile"))
)
