package game

import "github.com/kettek/morogue/id"

// Archetype is a structure that is used to act as a "template" for
// creating playable characters.
type Archetype struct {
	Title string
	ID    id.UUID

	Image string // Image for the archetype. Should be requested via HTTP to the resources backend.
	//
	Swole           AttributeLevel     // Raw Strength + Health
	Zooms           AttributeLevel     // Dex, basically
	Brains          AttributeLevel     // Thinkin' and spell-related
	Funk            AttributeLevel     // Charm and god-related
	Traits          []string           // Traits
	StartingObjects []id.UUID          // Starting objects
	StartingSkills  map[string]float64 // Starting skills
}
