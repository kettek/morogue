package game

import "math"

// Edible is an embed that provides logic for being eaten.
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

// NextCalories returns the amount of calories to eat next, limited to 300 or CurrentCalories.
func (e *Edible) NextCalories() int {
	max := 300
	return int(math.Min(float64(max), float64(e.CurrentCalories)))
}

// Throw throws the edible. TODO
func (e *Edible) Throw() {
}
