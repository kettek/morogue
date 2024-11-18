package server

import "github.com/kettek/morogue/game"

// Hurtable is the interface for objects that can be hurt.
type Hurtable interface {
	CalculateFromObject(o game.Object)
	TakeDamages(damages []game.DamageResult)
	TakeHeal(heal int)
	IsDead() bool
}

// Appliable is the interface for objects that can be applied.
type Appliable interface {
	Apply()
	Unapply()
	IsApplied() bool
}

// Damager is the interface for objects that can deal damage.
type Damager interface {
	CalculateFromCharacter(c *game.Character)
	RollDamages() []game.DamageResult
}

// Blockable is the interface for objects that can block.
type Blockable interface {
	IsBlocked() bool
}

// Openable is the interface for objects that can be opened.
type Openable interface {
	IsOpened() bool
	Open() error
	Close() error
}

// Lockable is the interface for objects that can be locked.
type Lockable interface {
	Lock()
	Unlock()
	IsLocked() bool
}

// Edible is the interface for objects that can be eaten.
type Edible interface {
	Eat() int
	Throw()
}
