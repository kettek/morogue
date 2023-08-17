package server

import "github.com/kettek/morogue/game"

type Hurtable interface {
	CalculateFromObject(o game.Object)
	TakeDamages(damages []game.DamageResult)
	TakeHeal(heal int)
	IsDead() bool
}

type Damager interface {
	CalculateFromCharacter(c *game.Character)
	RollDamages() []game.DamageResult
}

type Blockable interface {
	IsBlocked() bool
}

type Openable interface {
	IsOpened() bool
	Open() error
	Close() error
}

type Lockable interface {
	Lock()
	Unlock()
	IsLocked() bool
}
