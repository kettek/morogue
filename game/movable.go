package game

type Movable struct {
	Actions int `msgpack:"s"` // Amount of actions taken in a turn. This will generally be 1 for each player.
}

func (m *Movable) CalculateFromCharacter(c *Character) {
	value := c.Archetype.(CharacterArchetype).Zooms + c.Archetype.(CharacterArchetype).Funk/4
	value += c.Attributes.Zooms + c.Attributes.Funk/4
	m.Actions = 1 + int(value)/4
}
