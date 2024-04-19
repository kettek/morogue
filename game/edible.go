package game

type Edible struct {
	Calories        int `msgpack:"C,omitempty"`
	CurrentCalories int `msgpack:"c,omitempty"`
}

func (e *Edible) Eat() int {
	if e.CurrentCalories < 200 {
		return e.CurrentCalories
	}
	e.CurrentCalories -= 200
	return 200
}

func (e *Edible) Throw() {
	// TODO: ???
}
