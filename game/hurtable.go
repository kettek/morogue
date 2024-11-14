package game

import (
	"fmt"
)

type Hurtable struct {
	Health      int `msgpack:"h,omitempty"`
	MaxHealth   int `msgpack:"H,omitempty"`
	HealthRegen int `msgpack:"r,omitempty"`
	Downs       int `msgpack:"d,omitempty"`
	MaxDowns    int `msgpack:"D,omitempty"`
	MinArmor    int `msgpack:"a,omitempty"`
	MaxArmor    int `msgpack:"A,omitempty"`
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
	health += int(c.Swole()) * 2
	t := (c.Swole() + c.Zooms() + c.Brains() + c.Funk()) / 4

	h.MaxHealth = health + int(t)
	h.MaxDowns = 1 + int(c.Funk()/4)

	h.HealthRegen = 1 + int(c.Zooms()+c.Funk()/2+c.Swole()/4)/3
}

func (h *Hurtable) CalculateArmorFromCharacter(c *Character) {
	h.MinArmor = int(c.Swole()) / 2
	h.MaxArmor = int(c.Zooms()) / 2

	for _, a := range c.Inventory {
		if a, ok := a.(*Armor); ok {
			if !a.Applied || a.Archetype == nil {
				continue
			}
			h.MinArmor += a.Archetype.(ArmorArchetype).MinArmor
			h.MaxArmor += a.Archetype.(ArmorArchetype).MaxArmor
		}
	}
}

// String returns a string representation of the health.
func (h Hurtable) String() string {
	return fmt.Sprintf("%d/%d", h.Health, h.MaxHealth)
}

func (h Hurtable) ArmorRangeString() string {
	if h.MinArmor == 0 {
		return fmt.Sprintf("〜%d", h.MaxArmor)
	}
	if h.MinArmor == h.MaxArmor {
		return fmt.Sprintf("%d", h.MaxArmor)
	}
	return fmt.Sprintf("%d〜%d", h.MinArmor, h.MaxArmor)
}

func (h *Hurtable) TakeHeal(heal int) bool {
	if h.Health == h.MaxHealth {
		return false
	}
	h.Health += heal
	if h.Health > h.MaxHealth {
		h.Health = h.MaxHealth
	}
	return true
}

func (h *Hurtable) TakeDamages(damages []DamageResult) {
	for _, damage := range damages {
		h.Health -= damage.Damage
	}
	if h.Health < 0 {
		h.Downs++
	}
}

func (h *Hurtable) IsDead() bool {
	return h.Health <= 0 && h.Downs >= h.MaxDowns
}
