package game

import "github.com/kettek/morogue/id"

type Location struct {
	ID      id.UUID `msgpack:"id,omitempty"`
	Cells   Cells   `msgpack:"c,omitempty"`
	Objects Objects `msgpack:"o,omitempty"`
}

func (l *Location) Character(wid id.WID) *Character {
	for _, c := range l.Objects {
		if c, ok := c.(*Character); ok && c.WID == wid {
			return c
		}
	}
	return nil
}

func (l *Location) Characters() (chars []*Character) {
	for _, c := range l.Objects {
		if c, ok := c.(*Character); ok {
			chars = append(chars, c)
		}
	}
	return
}

func (l *Location) ObjectByWID(wid id.WID) Object {
	for _, o := range l.Objects {
		if o.GetWID() == wid {
			return o
		}
	}
	return nil
}
