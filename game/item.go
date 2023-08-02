package game

import "github.com/kettek/morogue/id"

// ItemArchetype is effectively a blueprint for an item.
type ItemArchetype struct {
	ID    id.UUID
	Title string `json:"T,omitempty"`
	Image string `json:"i,omitempty"`
}

func (a ItemArchetype) Type() string {
	return "item"
}

func (a ItemArchetype) GetID() id.UUID {
	return a.ID
}

// Item represents a generic item in the world.
type Item struct {
	Position
	Archetype id.UUID `json:"A,omitempty"`
	WID       id.WID  // ID assigned when entering da world.
	Container id.WID  `json:"c,omitempty"` // The container of the item, if any.
	ID        id.UUID `json:"id,omitempty"`
	Name      string  `json:"n,omitempty"`
}

// Type returns "item"
func (o Item) Type() ObjectType {
	return "item"
}

// GetWID returns the world ID of the item.
func (o Item) GetWID() id.WID {
	return o.WID
}

// GetPosition returns the position of the item.
func (o Item) GetPosition() Position {
	return o.Position
}

// SetPosition sets the position of the item.
func (o *Item) SetPosition(p Position) {
	o.Position = p
}
