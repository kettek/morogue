package game

type Edible struct {
	Calories int `msgpack:"c,omitempty"`
}

func (e *Edible) Eat() {
	// TODO: ???
}

func (e *Edible) Throw() {
	// TODO: ???
}
