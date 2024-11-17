package game

import (
	"fmt"
)

// Hurtable is an embed that provides logic for being hurt. This includes health, regen, downs, and armor.
type Hurtable struct {
	Health      int `msgpack:"h,omitempty"`
	MaxHealth   int `msgpack:"H,omitempty"`
	HealthRegen int `msgpack:"r,omitempty"`
	Downs       int `msgpack:"d,omitempty"`
	MaxDowns    int `msgpack:"D,omitempty"`
	MinArmor    int `msgpack:"a,omitempty"`
	MaxArmor    int `msgpack:"A,omitempty"`
}

// CalculateFromObject calculates hurtable values from an object.
func (h *Hurtable) CalculateFromObject(o Object) {
	switch o := o.(type) {
	case *Character:
		h.CalculateFromCharacter(o)
	case *Door:
		h.Health = 10
		h.MaxHealth = 10
	}
}

// CalculateFromCharacter calculates hurtable values from a character.
func (h *Hurtable) CalculateFromCharacter(c *Character) {
	health := 5
	health += int(c.Swole()) * 2
	t := (c.Swole() + c.Zooms() + c.Brains() + c.Funk()) / 4

	h.MaxHealth = health + int(t)
	h.MaxDowns = 1 + int(c.Funk()/4)

	h.HealthRegen = 1 + int(c.Zooms()+c.Funk()/2+c.Swole()/4)/3
}

// CalculateArmorFromCharacter calculates the armor from a character.
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

// ArmorRangeString returns a string representation of the armor range.
func (h Hurtable) ArmorRangeString() string {
	if h.MinArmor == 0 {
		return fmt.Sprintf("〜%d", h.MaxArmor)
	}
	if h.MinArmor == h.MaxArmor {
		return fmt.Sprintf("%d", h.MaxArmor)
	}
	return fmt.Sprintf("%d〜%d", h.MinArmor, h.MaxArmor)
}

// TakeHeal takes a heal and returns true if the character was healed.
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

// TakeDamages takes damages and applies them to the character. This will increment downs if health is below 0.
func (h *Hurtable) TakeDamages(damages []DamageResult) {
	for _, damage := range damages {
		h.Health -= damage.Damage
	}
	if h.Health < 0 {
		h.Downs++
	}
}

// IsDead returns true if the character is out of downs and their health is 0 or less.
func (h *Hurtable) IsDead() bool {
	return h.Health <= 0 && h.Downs >= h.MaxDowns
}
