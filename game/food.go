package game

import "github.com/kettek/morogue/id"

// FoodArchetype is effectively a blueprint for food.
type FoodArchetype struct {
	ID          id.UUID
	Title       string `msgpack:"T,omitempty"`
	Description string `msgpack:"d,omitempty"`
	Image       string `msgpack:"i,omitempty"`
	Calories    int    `msgpack:"c,omitempty"`
}

func (a FoodArchetype) Type() string {
	return "food"
}

func (a FoodArchetype) GetID() id.UUID {
	return a.ID
}

// Food represents a food object in the world.
type Food struct {
	Objectable
	Position
	Edible
	Name string `msgpack:"n,omitempty"`
}

// Type returns "food"
func (o Food) Type() ObjectType {
	return "food"
}
