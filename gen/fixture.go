package gen

import "github.com/kettek/morogue/id"

type Fixture struct {
	ID   id.UUID
	Keys map[string]id.UUID
	Rows []string
}

func (f Fixture) Width() (w int) {
	for _, row := range f.Rows {
		if len(row) > w {
			w = len(row)
		}
	}
	return
}

func (f Fixture) Height() int {
	return len(f.Rows)
}
