package game

import "math"

type Edible struct {
	Calories        int `msgpack:"C,omitempty"`
	CurrentCalories int `msgpack:"c,omitempty"`
}

// Eat consumes a portion of the edible and returns the number of calories consumed.
func (e *Edible) Eat() int {
	max := 300
	e.CurrentCalories -= max
	if e.CurrentCalories < 0 {
		max += e.CurrentCalories
		e.CurrentCalories = 0
	}
	return max
}

func (e *Edible) NextCalories() int {
	max := 300
	return int(math.Min(float64(max), float64(e.CurrentCalories)))
}

func (e *Edible) Throw() {
	// TODO: ???
}
