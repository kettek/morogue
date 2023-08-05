package client

import "github.com/kettek/morogue/client/ifs"

type stateMachine struct {
}

func (a *app) Top() ifs.State {
	if len(a.states) == 0 {
		return nil
	}
	return a.states[len(a.states)-1]
}

func (a *app) Pop() (ifs.State, error) {
	s := a.Top()
	if s == nil {
		return nil, nil
	}
	r, err := s.End()
	if err != nil {
		return nil, err
	}
	a.states = a.states[:len(a.states)-1]
	s = a.Top()
	if s == nil {
		return nil, nil
	}
	err = s.Return(r)
	return s, err
}

func (a *app) Push(s ifs.State) error {
	t := a.Top()
	if t != nil {
		err := t.Leave()
		if err != nil {
			return err
		}
	}
	a.states = append(a.states, s)
	err := s.Begin(a.runContext)
	return err
}
