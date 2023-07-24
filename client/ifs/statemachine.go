package ifs

type StateMachine interface {
	Top() State
	Pop() (State, error)
	Push(State) error
}
