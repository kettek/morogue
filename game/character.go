package game

import (
	"fmt"

	"github.com/kettek/morogue/id"
)

// CharacterArchetype is a structure that is used to act as a "template" for
// creating playable characters.
type CharacterArchetype struct {
	Title string  `msgpack:"t,omitempty"` // Title of the archetype. This is what is displayed to the user.
	ID    id.UUID `msgpack:"id,omitempty"`

	PlayerOnly bool // If the archetype is for players only during character creation.

	Image string `msgpack:"i,omitempty"` // Image for the archetype. Should be requested via HTTP to the resources backend.
	//
	Swole  AttributeLevel // Raw Strength + Health
	Zooms  AttributeLevel // Dex, basically
	Brains AttributeLevel // Thinkin' and spell-related
	Funk   AttributeLevel // Charm and god-related
	//Traits          []string           // Traits
	Traits          TraitList
	Slots           Slots              // Slots
	StartingObjects []id.UUID          // Starting objects
	StartingSkills  map[string]float64 // Starting skills
}

// Type returns "character"
func (c CharacterArchetype) Type() string {
	return "character"
}

// GetID returns the ID.
func (c CharacterArchetype) GetID() id.UUID {
	return c.ID
}

// Character represents a character. This can be a player or an NPC.
type Character struct {
	Position
	Blockable
	Hurtable
	Damager
	Events      []Event   `msgpack:"-" json:"-"` // Events that have happened to the character. These are only sent to the owning client.
	Desire      Desire    `msgpack:"-" json:"-"` // The current desire of the character. Used server-side.
	LastDesire  Desire    `msgpack:"-" json:"-"` // Last desire processed. Used server-side.
	WID         id.WID    // ID assigned when entering da world.
	Archetype   Archetype `msgpack:"-" json:"-"` // Archetype of the character. This is likely a pointer.
	ArchetypeID id.UUID   `msgpack:"A,omitempty"`
	Name        string    `msgpack:"n,omitempty"`
	Level       int       `msgpack:"l,omitempty"`
	Slots       SlotMap   `msgpack:"-"`
	Skills      Skills    `msgpack:"-"`
	Inventory   Objects   `msgpack:"-"`
	//
}

// Type returns "character"
func (c Character) Type() ObjectType {
	return "character"
}

// GetWID returns the WID.
func (c Character) GetWID() id.WID {
	return c.WID
}

// SetWID sets the WID.
func (c *Character) SetWID(wid id.WID) {
	c.WID = wid
}

// GetPosition returns the position.
func (c Character) GetPosition() Position {
	return c.Position
}

// SetPosition sets the position.
func (c *Character) SetPosition(p Position) {
	c.Position = p
}

// GetArchetypeID returns the archetype.
func (c *Character) GetArchetypeID() id.UUID {
	return c.ArchetypeID
}

// SetArchetypeID sets the archetype.
func (c *Character) SetArchetype(a Archetype) {
	c.Archetype = a
}

// GetArchetype returns the archetype.
func (c *Character) GetArchetype() Archetype {
	return c.Archetype
}

// InInventory returns true if the character has the object in their inventory.
func (c *Character) InInventory(wid id.WID) bool {
	for _, i := range c.Inventory {
		if i.GetWID() == wid {
			return true
		}
	}
	return false
}

// Apply applies an object from the character's inventory.
func (c *Character) Apply(o Object, force bool) Event {
	if !c.InInventory(o.GetWID()) {
		return EventNotice{
			Message: "You don't have that item.",
		}
	}
	switch o := o.(type) {
	case *Weapon:
		return c.applyWeapon(o, force)
	case *Armor:
		return c.applyArmor(o, force)
	case *Item:
		// TODO
	}
	return nil
}

// applyWeapon applies a weapon to the character.
func (c *Character) applyWeapon(w *Weapon, force bool) Event {
	if w.Archetype != nil {
		for _, trait := range c.Archetype.(CharacterArchetype).Traits {
			if !trait.CanApply(w) {
				return EventNotice{
					Message: fmt.Sprintf("You can't use %s.", w.Archetype.(WeaponArchetype).Title),
				}
			}
		}

		if err := c.Slots.Apply(w.Archetype.(WeaponArchetype).Slots); err != nil {
			if !force {
				return EventNotice{
					Message: err.Error(),
				}
			}
		}
	}

	w.Applied = true

	c.Damager.CalculateFromCharacter(c)

	return EventApply{
		Applier: c.WID,
		WID:     w.WID,
		Applied: true,
	}
}

// applyArmor applies an armor to the character.
func (c *Character) applyArmor(a *Armor, force bool) Event {
	if a.Archetype != nil {
		for _, trait := range c.Archetype.(CharacterArchetype).Traits {
			if !trait.CanApply(a) {
				return EventNotice{
					Message: fmt.Sprintf("You can't use %s.", a.Archetype.(ArmorArchetype).Title),
				}
			}
		}

		if err := c.Slots.Apply(a.Archetype.(ArmorArchetype).Slots); err != nil {
			if !force {
				return EventNotice{
					Message: err.Error(),
				}
			}
		}
	}

	a.Applied = true
	return EventApply{
		Applier: c.WID,
		WID:     a.WID,
		Applied: true,
	}
}

// Unapply unapplies an object from the character's inventory.
func (c *Character) Unapply(o Object, force bool) Event {
	if !c.InInventory(o.GetWID()) {
		return EventNotice{
			Message: "You don't have that item.",
		}
	}
	switch o := o.(type) {
	case *Weapon:
		return c.unapplyWeapon(o, force)
	case *Armor:
		return c.unapplyArmor(o, force)
	case *Item:
		// TODO
	}
	return nil
}

// unapplyWeapon unapplies a weapon from the character.
func (c *Character) unapplyWeapon(w *Weapon, force bool) Event {
	if w.Archetype != nil {
		if err := c.Slots.Unapply(w.Archetype.(WeaponArchetype).Slots); err != nil {
			if !force {
				return EventNotice{
					Message: err.Error(),
				}
			}
		}
	}

	w.Applied = false

	c.Damager.CalculateFromCharacter(c)

	return EventApply{
		Applier: c.WID,
		WID:     w.WID,
		Applied: false,
	}
}

// unapplyArmor unapplies an armor from the character.
func (c *Character) unapplyArmor(a *Armor, force bool) Event {
	if a.Archetype != nil {
		if err := c.Slots.Unapply(a.Archetype.(ArmorArchetype).Slots); err != nil {
			if !force {
				return EventNotice{
					Message: err.Error(),
				}
			}
		}
	}

	a.Applied = false
	return EventApply{
		Applier: c.WID,
		WID:     a.WID,
		Applied: false,
	}
}

// Pickup adds an object to the character's inventory.
func (c *Character) Pickup(o Object) Event {
	c.Inventory = append(c.Inventory, o)

	// It's a bit cheesy, but we use -1/-1 to signify off the map.
	o.SetPosition(Position{-1, -1})

	return EventPickup{
		Picker: c.WID,
		WID:    o.GetWID(),
	}
}

// Drop removes an object from the character's inventory.
func (c *Character) Drop(o Object) Event {
	if !c.InInventory(o.GetWID()) {
		return EventNotice{
			Message: "You don't have that item.",
		}
	}

	// TODO: Here is where we could check for items that cannot be dropped, such as cursed items.

	// Unapply it, for obvious reasons.
	switch o := o.(type) {
	case *Weapon:
		c.unapplyWeapon(o, true)
	case *Armor:
		c.unapplyArmor(o, true)
	}

	// Remove from the inventory.
	for i, item := range c.Inventory {
		if item.GetWID() == o.GetWID() {
			c.Inventory = append(c.Inventory[:i], c.Inventory[i+1:]...)
			break
		}
	}

	return EventDrop{
		Dropper:  c.WID,
		Object:   o,
		Position: c.GetPosition(),
	}
}

// Health represents a character's health.
type Health struct {
	Current int `webpack:"c,omitempty"`
	Max     int `webpack:"m,omitempty"`
}

// String returns a string representation of the health.
func (h Health) String() string {
	return fmt.Sprintf("%d/%d", h.Current, h.Max)
}
