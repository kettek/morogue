package game

// Hungerable is an embed that provides logic for being hungry.
type Hungerable struct {
	Hunger    int `msgpack:"h,omitempty"`
	MaxHunger int `msgpack:"H,omitempty"`
}

// CalculateFromCharacter calculates the hunger from a character.
func (h *Hungerable) CalculateFromCharacter(c *Character) {
	maxHunger := 1700 + int(c.Swole()*300)
	hunger := int(float64(maxHunger) * (float64(h.Hunger) / float64(h.MaxHunger)))

	if hunger < 0 {
		h.Hunger = 0
	} else {
		h.Hunger = hunger
	}
	h.MaxHunger = maxHunger
}

// UseEnergy uses energy and returns true if the character is still hungry.
func (h *Hungerable) UseEnergy(amount int) bool {
	h.Hunger -= amount
	if h.Hunger < 0 {
		h.Hunger = 0
	}
	return h.Hunger != 0
}
