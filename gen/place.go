package gen

import "github.com/kettek/morogue/id"

type Place struct {
	Title    string
	ID       id.UUID
	Width    MinMax
	Height   MinMax
	Fixtures []FixtureEntry
	WFC      []WFCEntry
}

type FixtureEntry struct {
	Targets []FixtureTarget
	Count   MinMax
	X       MinMax
	Y       MinMax
}

type FixtureTarget struct {
	ID        id.UUID
	Rotate    bool
	Overlap   bool
	Intersect bool
}

type WFCEntry struct {
	ID       id.UUID
	Adjacent []id.UUID
}
