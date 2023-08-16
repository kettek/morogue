package game

// Lockable is a feature that enables locking and unlocking
type Lockable struct {
	Locked bool `msgpack:"l"`
}

// Lock locks the lockable
func (l *Lockable) Lock() {
	l.Locked = true
}

// Unlock unlocks the lockable
func (l *Lockable) Unlock() {
	l.Locked = false
}

// IsLocked returns true if the lockable is locked
func (l *Lockable) IsLocked() bool {
	return l.Locked
}
