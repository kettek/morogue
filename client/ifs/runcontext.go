package ifs

// RunContext provides various structures useful during Update calls.
type RunContext struct {
	Sm   StateMachine
	Txt  *TextRenderer
	UI   *DrawContextUI
	Game *GameContext
}
