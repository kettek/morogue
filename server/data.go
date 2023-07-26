package main

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
	Tiles      []game.Tile
}

func (d *Data) hasArchetype(uuid id.UUID) bool {
	for _, a := range d.Archetypes {
		if a.UUID == uuid {
			return true
		}
	}
	return false
}

func (d *Data) Tile(uuid id.UUID) (game.Tile, error) {
	for _, t := range d.Tiles {
		if t.UUID == uuid {
			return t, nil
		}
	}
	return game.Tile{}, errors.New("no such tile")
}

func (d *Data) loadArchetypes() error {
	entries, err := os.ReadDir("archetypes")
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".json") {
			bytes, err := os.ReadFile(filepath.Join("archetypes", entry.Name()))
			if err != nil {
				log.Println(err)
				continue
			}
			var a game.Archetype
			err = json.Unmarshal(bytes, &a)
			if err != nil {
				log.Println(err)
				continue
			}
			d.Archetypes = append(d.Archetypes, a)
		}
	}
	return nil
}

func (d *Data) loadTiles() error {
	entries, err := os.ReadDir("tiles")
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".json") {
			bytes, err := os.ReadFile(filepath.Join("tiles", entry.Name()))
			if err != nil {
				log.Println(err)
				continue
			}
			var t game.Tile
			err = json.Unmarshal(bytes, &t)
			if err != nil {
				log.Println(err)
				continue
			}
			d.Tiles = append(d.Tiles, t)
		}
	}
	return nil
}
