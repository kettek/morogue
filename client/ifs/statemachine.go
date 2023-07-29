package ifs

// StateMachine is what you would expect it to be.
type StateMachine interface {
	Top() State
	Pop() (State, error)
	Push(State) error
}
