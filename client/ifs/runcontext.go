package ifs

type RunContext struct {
	Sm   StateMachine
	Txt  *TextRenderer
	UI   *DrawContextUI
	Game *GameContext
}
