package game

import "github.com/kettek/morogue/id"

// Character represents a playable character.
type Character struct {
	Position
	Events     []Event  `json:"-"` // Events that have happened to the character. These are only sent to the owning client.
	Desire     Desire   `json:"-"` // The current desire of the character. Used server-side.
	LastDesire Desire   `json:"-"` // Last desire processed. Used server-side.
	WID        id.WID   // ID assigned when entering da world.
	Archetype  id.UUID  `json:"a,omitempty"`
	Name       string   `json:"n,omitempty"`
	Level      int      `json:"l,omitempty"`
	Skills     Skills   `json:"-"`
	Inventory  []Object `json:"-"`
}

// Type returns "character"
func (c Character) Type() ObjectType {
	return "character"
}

// GetWID returns the WID.
func (c Character) GetWID() id.WID {
	return c.WID
}

// GetPosition returns the position.
func (c Character) GetPosition() Position {
	return c.Position
}

// SetPosition sets the position.
func (c *Character) SetPosition(p Position) {
	c.Position = p
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

	return EventDrop{
		Dropper:  c.WID,
		WID:      o.GetWID(),
		Position: c.GetPosition(),
	}
}
