package game

import "github.com/kettek/morogue/id"

type Containerable struct {
	Inventory    Objects `msgpack:"o"`
	ContainerWID id.WID  `json:"-"` // ID of the container -- this should be the same as the Objectable's.
	Capacity     int     `msgpack:"c"`
	Limit        int     `msgpack:"l"`
}

// InInventory returns true if the containerable has the object in its inventory.
func (c *Containerable) InInventory(wid id.WID) bool {
	for _, i := range c.Inventory {
		if i.GetWID() == wid {
			return true
		}
	}
	return false
}

// Pickup adds an object to the containerable's inventory.
func (c *Containerable) Pickup(o Object) bool {
	c.Inventory = append(c.Inventory, o)

	// Set container to the container.
	o.SetContainerWID(c.ContainerWID)

	// It's a bit cheesy, but we use -1/-1 to signify off the map.
	o.SetPosition(Position{-1, -1})

	return true
}

// Drop removes an object from the containerable's inventory.
func (c *Containerable) Drop(o Object) bool {
	if !c.InInventory(o.GetWID()) {
		return false
	}

	// Remove from the inventory.
	for i, item := range c.Inventory {
		if item.GetWID() == o.GetWID() {
			c.Inventory = append(c.Inventory[:i], c.Inventory[i+1:]...)
			break
		}
	}

	// Clear the container.
	o.SetContainerWID(0)

	return true
}
