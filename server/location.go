package main

import (
	"fmt"

	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/gen"
	"github.com/kettek/morogue/id"
)

type location struct {
	game.Location
	active bool
}

func (l *location) generate() error {
	// FIXME: Pull from somewhere.
	for x := 0; x < 10; x++ {
		l.Cells = append(l.Cells, make([]game.Cell, 0))
		for y := 0; y < 10; y++ {
			l.Cells[x] = append(l.Cells[x], game.Cell{})
		}
	}

	gen.Generate(gen.Styles[gen.StyleBox], gen.Config{
		Width:  10,
		Height: 10,
		Cell: func(x, y int) gen.Cell {
			return &l.Cells[x][y]
		},
		SetCell: func(x, y int, cell gen.Cell) {
			if cell.Flags().Has("blocked") {
				l.Cells[x][y].Blocks = game.MovementAll
				if tid, err := id.UID(id.Tile, "stone-wall"); err == nil {
					l.Cells[x][y].TileID = &tid
				}
			}
		},
	})
	return nil
}

func (l *location) process() error {
	for _, c := range l.Characters {
		fmt.Println("TODO: Handle character", c)
	}
	for _, m := range l.Mobs {
		fmt.Println("TODO: Handle mob", m)
	}
	for _, o := range l.Objects {
		fmt.Println("TODO: Handle object", o)
	}
	return nil
}

type locationConfig struct {
	ID    id.UUID
	Depth int
}
