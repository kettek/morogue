package main

import "github.com/kettek/morogue/client/ifs"

type stateMachine struct {
	states []ifs.State
}

func (sm *stateMachine) Top() ifs.State {
	if len(sm.states) == 0 {
		return nil
	}
	return sm.states[len(sm.states)-1]
}

func (sm *stateMachine) Pop() (ifs.State, error) {
	s := sm.Top()
	if s == nil {
		return nil, nil
	}
	r, err := s.End()
	if err != nil {
		return nil, err
	}
	sm.states = sm.states[:len(sm.states)-1]
	s = sm.Top()
	if s == nil {
		return nil, nil
	}
	err = s.Return(r)
	return s, err
}

func (sm *stateMachine) Push(s ifs.State) error {
	t := sm.Top()
	if t != nil {
		err := t.Leave()
		if err != nil {
			return err
		}
	}
	sm.states = append(sm.states, s)
	err := s.Begin()
	return err
}
