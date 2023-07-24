package game

type Archetype struct {
	Title string
	Image string // Image for the archetype. Should be requested via HTTP to the resources backend.
	//
	Swole  AttributeLevel // Raw Strength + Health
	Zooms  AttributeLevel // Dex, basically
	Brains AttributeLevel // Thinkin' and spell-related
	Funk   AttributeLevel // Charm and god-related
}
