package server

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/id"
)

type Data struct {
	Archetypes []game.Archetype
}

func (d *Data) hasArchetype(uuid id.UUID) bool {
	for _, a := range d.Archetypes {
		if a.GetID() == uuid {
			return true
		}
	}
	return false
}

func (d *Data) Archetype(uuid id.UUID) game.Archetype {
	for _, a := range d.Archetypes {
		if a.GetID() == uuid {
			return a
		}
	}
	return nil
}

func (d *Data) Tile(uuid id.UUID) (game.TileArchetype, error) {
	for _, t := range d.TileArchetypes() {
		if t.ID == uuid {
			return t, nil
		}
	}
	return game.TileArchetype{}, errors.New("no such tile")
}

func (d *Data) TileArchetypes() []game.TileArchetype {
	var archetypes []game.TileArchetype
	for _, a := range d.Archetypes {
		if t, ok := a.(game.TileArchetype); ok {
			archetypes = append(archetypes, t)
		}
	}
	return archetypes
}

func (d *Data) CharacterArchetypes() []game.CharacterArchetype {
	var archetypes []game.CharacterArchetype
	for _, a := range d.Archetypes {
		if c, ok := a.(game.CharacterArchetype); ok {
			archetypes = append(archetypes, c)
		}
	}
	return archetypes
}

func (d *Data) ItemArchetypes() []game.ItemArchetype {
	var archetypes []game.ItemArchetype
	for _, a := range d.Archetypes {
		if i, ok := a.(game.ItemArchetype); ok {
			archetypes = append(archetypes, i)
		}
	}
	return archetypes
}

func (d *Data) WeaponArchetypes() []game.WeaponArchetype {
	var archetypes []game.WeaponArchetype
	for _, a := range d.Archetypes {
		if w, ok := a.(game.WeaponArchetype); ok {
			archetypes = append(archetypes, w)
		}
	}
	return archetypes
}

func (d *Data) ArmorArchetypes() []game.ArmorArchetype {
	var archetypes []game.ArmorArchetype
	for _, a := range d.Archetypes {
		if a, ok := a.(game.ArmorArchetype); ok {
			archetypes = append(archetypes, a)
		}
	}
	return archetypes
}

func (d *Data) LoadArchetypes() error {
	entries, err := os.ReadDir("archetypes")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			var kind string
			switch entry.Name() {
			case "characters":
				kind = "characters"
			case "weapons":
				kind = "weapons"
			case "armors":
				kind = "armors"
			case "items":
				kind = "items"
			case "tiles":
				kind = "tiles"
			case "doors":
				kind = "doors"
			}
			if kind == "" {
				continue
			}

			fullpath := filepath.Join("archetypes", entry.Name())
			entries, err := os.ReadDir(fullpath)
			if err == nil {
				for _, entry := range entries {
					if strings.HasSuffix(entry.Name(), ".json") {
						bytes, err := os.ReadFile(filepath.Join(fullpath, entry.Name()))
						if err != nil {
							log.Println(err)
							continue
						}

						if kind == "characters" {
							var c game.CharacterArchetype
							err = json.Unmarshal(bytes, &c)
							if err != nil {
								log.Println(err)
								continue
							}
							c.PlayerOnly = true
							c.Image = "characters/" + c.Image
							d.Archetypes = append(d.Archetypes, c)
						} else if kind == "weapons" {
							var w game.WeaponArchetype
							err = json.Unmarshal(bytes, &w)
							if err != nil {
								log.Println(err)
								continue
							}
							w.Image = "weapons/" + w.Image
							d.Archetypes = append(d.Archetypes, w)
						} else if kind == "armors" {
							var a game.ArmorArchetype
							err = json.Unmarshal(bytes, &a)
							if err != nil {
								log.Println(err)
								continue
							}
							a.Image = "armors/" + a.Image
							d.Archetypes = append(d.Archetypes, a)
						} else if kind == "items" {
							var i game.ItemArchetype
							err = json.Unmarshal(bytes, &i)
							if err != nil {
								log.Println(err)
								continue
							}
							i.Image = "items/" + i.Image
							d.Archetypes = append(d.Archetypes, i)
						} else if kind == "tiles" {
							var t game.TileArchetype
							err = json.Unmarshal(bytes, &t)
							if err != nil {
								log.Println(err)
								continue
							}
							t.Image = "tiles/" + t.Image
							d.Archetypes = append(d.Archetypes, t)
						} else if kind == "doors" {
							var a game.DoorArchetype
							err = json.Unmarshal(bytes, &a)
							if err != nil {
								log.Println(err)
								continue
							}
							a.Image = "doors/" + a.Image
							d.Archetypes = append(d.Archetypes, a)
						}
					}
				}
			}
		}
	}
	return nil
}
