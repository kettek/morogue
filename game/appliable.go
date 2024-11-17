package game

// Appliable is an embed that provides logic for being applied or unapplied.
type Appliable struct {
	Applied bool `msgpack:"a,omitempty"`
}

// Apply sets the appliable state to true.
func (a *Appliable) Apply() {
	a.Applied = true
}

// Unapply sets the appliable state to false.
func (a *Appliable) Unapply() {
	a.Applied = false
}

// IsApplied returns the appliable state.
func (a *Appliable) IsApplied() bool {
	return a.Applied
}
