package server

import "github.com/kettek/morogue/game"

type Hurtable interface {
	Damage(damages ...game.Damage) (results []game.DamageResult)
	Heal(heal int)
	IsDead() bool
}
