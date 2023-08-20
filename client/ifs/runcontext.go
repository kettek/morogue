package ifs

// RunContext provides various structures useful during Update calls.
type RunContext struct {
	Cfg  *Configuration
	Sm   StateMachine
	Txt  *TextRenderer
	UI   *DrawContextUI
	Game *GameContext
}
