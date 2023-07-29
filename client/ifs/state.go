package ifs

// State is what a state machine expects.
type State interface {
	// Leave is called on a State when another State is pushed into the state machine.
	Leave() error
	// Return is called on the previous state from the top-level state when Pop is called.
	Return(interface{}) error
	// Begin is called on a State when it is pushed into the state machine.
	Begin(RunContext) error
	// Update is called once per tick.
	Update(RunContext) error
	// Draw is called when drawing occurs.
	Draw(DrawContext)
	// End is called on a State when it is popped.
	End() (interface{}, error)
}
