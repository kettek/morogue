package game

import (
	"fmt"
)

type Hurtable struct {
	Health    int `msgpack:"h,omitempty"`
	MaxHealth int `msgpack:"H,omitempty"`
	Downs     int `msgpack:"d,omitempty"`
	MaxDowns  int `msgpack:"D,omitempty"`
}

func (h *Hurtable) CalculateFromObject(o Object) {
	switch o := o.(type) {
	case *Character:
		h.CalculateFromCharacter(o)
	case *Door:
		h.Health = 10
		h.MaxHealth = 10
	}
}

func (h *Hurtable) CalculateFromCharacter(c *Character) {
	health := 5
	health += int(c.Archetype.(CharacterArchetype).Swole) * 2
	t := (c.Archetype.(CharacterArchetype).Swole + c.Archetype.(CharacterArchetype).Zooms + c.Archetype.(CharacterArchetype).Brains + c.Archetype.(CharacterArchetype).Funk) / 4

	h.MaxHealth = health + int(t)
	h.MaxDowns = 1 + int(c.Archetype.(CharacterArchetype).Funk/4)
}

// String returns a string representation of the health.
func (h Hurtable) String() string {
	return fmt.Sprintf("%d/%d", h.Health, h.MaxHealth)
}

func (h *Hurtable) Damage(damages ...Damage) (results []DamageResult) {
	for _, damage := range damages {
		dmg := damage.Roll()
		h.Health -= dmg
		results = append(results, DamageResult{
			Damage: dmg,
		})
	}
	return results
}

func (h *Hurtable) Heal(heal int) {
	h.Health += heal
}

func (h *Hurtable) IsDead() bool {
	return h.Health <= 0 && h.Downs >= h.MaxDowns
}
