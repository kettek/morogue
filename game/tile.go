package game

import "github.com/kettek/morogue/id"

// Tile is a structure that is used to provide the behavior
// and appearance of a given cell.
type Tile struct {
	Title string
	ID    id.UUID

	Image string // Image for the tile. It should be requested via HTTP to the resources backend.
	// TODO: Other tile properties.
}
