package game

import (
	"github.com/kettek/morogue/id"
)

// Location represents a location in the game world. It contains cells and objects.
type Location struct {
	ID      id.UUID `msgpack:"id,omitempty"`
	Cells   Cells   `msgpack:"c,omitempty"`
	Objects Objects `msgpack:"o,omitempty"`
}

// Character returns a Character associated with wid.
func (l *Location) Character(wid id.WID) *Character {
	for _, c := range l.Objects {
		if c, ok := c.(*Character); ok && c.WID == wid {
			return c
		}
	}
	return nil
}

// Characters returns all Characters in the location.
func (l *Location) Characters() (chars []*Character) {
	for _, c := range l.Objects {
		if c, ok := c.(*Character); ok {
			chars = append(chars, c)
		}
	}
	return
}

// ObjectByWID returns an object by its WID.
func (l *Location) ObjectByWID(wid id.WID) Object {
	for _, o := range l.Objects {
		if o.GetWID() == wid {
			return o
		}
	}
	return nil
}
