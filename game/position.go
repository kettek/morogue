package game

// Position represents a position in the world.
type Position struct {
	X int `msgpack:"x,omitempty"`
	Y int `msgpack:"y,omitempty"`
}

// GetPosition returns the position.
func (o Position) GetPosition() Position {
	return o
}

// SetPosition sets the position.
func (o *Position) SetPosition(p Position) {
	o.X = p.X
	o.Y = p.Y
}
