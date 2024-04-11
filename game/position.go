package game

type Position struct {
	X int `msgpack:"x,omitempty"`
	Y int `msgpack:"y,omitempty"`
}

func (o Position) GetPosition() Position {
	return o
}

func (o *Position) SetPosition(p Position) {
	o.X = p.X
	o.Y = p.Y
}
