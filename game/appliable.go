package game

type Appliable struct {
	Applied bool `msgpack:"a,omitempty"`
}

func (a *Appliable) Apply() {
	a.Applied = true
}

func (a *Appliable) Unapply() {
	a.Applied = false
}

func (a *Appliable) IsApplied() bool {
	return a.Applied
}
