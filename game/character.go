package game

import (
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
	Swole           AttributeLevel     // Raw Strength + Health
	Zooms           AttributeLevel     // Dex, basically
	Brains          AttributeLevel     // Thinkin' and spell-related
	Funk            AttributeLevel     // Charm and god-related
	Traits          []string           // Traits
	StartingObjects []id.UUID          // Starting objects
	StartingSkills  map[string]float64 // Starting skills
}

func (c CharacterArchetype) Type() string {
	return "character"
}

func (c CharacterArchetype) GetID() id.UUID {
	return c.ID
}

// Character represents a character. This can be a player or an NPC.
type Character struct {
	Position
	Events     []Event `msgpack:"-" json:"-"` // Events that have happened to the character. These are only sent to the owning client.
	Desire     Desire  `msgpack:"-" json:"-"` // The current desire of the character. Used server-side.
	LastDesire Desire  `msgpack:"-" json:"-"` // Last desire processed. Used server-side.
	WID        id.WID  // ID assigned when entering da world.
	Archetype  id.UUID `msgpack:"A,omitempty"`
	Name       string  `msgpack:"n,omitempty"`
	Level      int     `msgpack:"l,omitempty"`
	Skills     Skills  `msgpack:"-"`
	Inventory  Objects `msgpack:"-"`
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

// GetArchetype returns the archetype.
func (c *Character) GetArchetype() id.UUID {
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
func (c *Character) Apply(o Object) Event {
	if !c.InInventory(o.GetWID()) {
		return EventNotice{
			Message: "You don't have that item.",
		}
	}
	switch o := o.(type) {
	case *Weapon:
		return c.applyWeapon(o)
	case *Armor:
		return c.applyArmor(o)
	case *Item:
		// TODO
	}
	return nil
}

func (c *Character) applyWeapon(w *Weapon) Event {
	w.Applied = true
	return EventApply{
		Applier: c.WID,
		WID:     w.WID,
		Applied: true,
	}
}

func (c *Character) applyArmor(a *Armor) Event {
	a.Applied = true
	return EventApply{
		Applier: c.WID,
		WID:     a.WID,
		Applied: true,
	}
}

// Unapply unapplies an object from the character's inventory.
func (c *Character) Unapply(o Object) Event {
	if !c.InInventory(o.GetWID()) {
		return EventNotice{
			Message: "You don't have that item.",
		}
	}
	switch o := o.(type) {
	case *Weapon:
		return c.unapplyWeapon(o)
	case *Armor:
		return c.unapplyArmor(o)
	case *Item:
		// TODO
	}
	return nil
}

func (c *Character) unapplyWeapon(w *Weapon) Event {
	w.Applied = false
	return EventApply{
		Applier: c.WID,
		WID:     w.WID,
		Applied: false,
	}
}

func (c *Character) unapplyArmor(a *Armor) Event {
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
		c.unapplyWeapon(o)
	case *Armor:
		c.unapplyArmor(o)
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
		WID:      o.GetWID(),
		Position: c.GetPosition(),
	}
}
