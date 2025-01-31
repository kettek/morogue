package game

import (
	"encoding/json"
	"errors"

	"github.com/kettek/morogue/id"
	"github.com/vmihailenco/msgpack/v5"
)

// ObjectType is a key that represents an object's type. This is used
// for safely marshalling/unmarshalling the Object interface.
type ObjectType string

// Object is an interface intended for location objects. This includes
// players, items, and enemies.
type Object interface {
	Type() ObjectType
	SetWID(id.WID) // SetWID should only be called by the world.
	GetWID() id.WID
	SetPosition(Position)
	GetPosition() Position
	GetArchetypeID() id.UUID
	SetArchetype(a Archetype)
	GetArchetype() Archetype
	GetContainerWID() id.WID
	SetContainerWID(id.WID)
}

// CreateObjectFromArchetype creates an Object interface from an Archetype interface.
func CreateObjectFromArchetype(a Archetype) Object {
	switch a := a.(type) {
	case CharacterArchetype:
		return &Character{
			Objectable: Objectable{
				ArchetypeID: a.GetID(),
				Archetype:   a,
			},
			Slots: a.Slots.ToMap(),
		}
	case WeaponArchetype:
		return &Weapon{
			Objectable: Objectable{
				ArchetypeID: a.GetID(),
				Archetype:   a,
			},
		}
	case ArmorArchetype:
		return &Armor{
			Objectable: Objectable{
				ArchetypeID: a.GetID(),
				Archetype:   a,
			},
		}
	case DoorArchetype:
		return &Door{
			Objectable: Objectable{
				ArchetypeID: a.GetID(),
				Archetype:   a,
			},
			Blockable: Blockable{
				BlockType: a.BlockType,
			},
		}
	case FoodArchetype:
		return &Food{
			Objectable: Objectable{
				ArchetypeID: a.GetID(),
				Archetype:   a,
			},
			Edible: Edible{
				Calories:        a.Calories,
				CurrentCalories: a.Calories,
			},
		}
	}
	return nil
}

// RawMessage is used to wrap JSON to be later unmarshalled.
type RawMessage []byte

// MarshalMsgpack marshals the RawMessage to msgpack.
func (m RawMessage) MarshalMsgpack() ([]byte, error) {
	return msgpack.Marshal((msgpack.RawMessage)(m))
}

// UnmarshalMsgpack unmarshals the RawMessage from msgpack.
func (m *RawMessage) UnmarshalMsgpack(b []byte) error {
	return msgpack.Unmarshal(b, (*msgpack.RawMessage)(m))
}

// MarshalJSON marshals the RawMessage to JSON.
func (m RawMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal((json.RawMessage)(m))
}

// UnmarshalJSON unmarshals the RawMessage from JSON.
func (m *RawMessage) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, (*json.RawMessage)(m))
}

// ObjectWrapper wraps an Object interface for msgpack marshal and unmarshal.
type ObjectWrapper struct {
	Type ObjectType `msgpack:"t"`
	Data RawMessage `msgpack:"d"`
}

// Object returns the wrapped Object. To add additional types, they must be
// handled here and ObjectJSON.
func (ow ObjectWrapper) Object() (Object, error) {
	switch ow.Type {
	case (Item{}).Type():
		var o *Item
		if err := msgpack.Unmarshal(ow.Data, &o); err != nil {
			return nil, err
		}
		return o, nil
	case (Character{}).Type():
		var c *Character
		if err := msgpack.Unmarshal(ow.Data, &c); err != nil {
			return nil, err
		}
		return c, nil
	case (Weapon{}).Type():
		var w *Weapon
		if err := msgpack.Unmarshal(ow.Data, &w); err != nil {
			return nil, err
		}
		return w, nil
	case (Armor{}).Type():
		var a *Armor
		if err := msgpack.Unmarshal(ow.Data, &a); err != nil {
			return nil, err
		}
		return a, nil
	case (Food{}).Type():
		var f *Food
		if err := msgpack.Unmarshal(ow.Data, &f); err != nil {
			return nil, err
		}
		return f, nil
	}
	return nil, errors.New("unknown object type: " + string(ow.Type))
}

// ObjectJSON returns the wrapped Object. To add additional types, they must be
// handled here and Object.
func (ow ObjectWrapper) ObjectJSON() (Object, error) {
	switch ow.Type {
	case (Item{}).Type():
		var o *Item
		if err := json.Unmarshal(ow.Data, &o); err != nil {
			return nil, err
		}
		return o, nil
	case (Character{}).Type():
		var c *Character
		if err := json.Unmarshal(ow.Data, &c); err != nil {
			return nil, err
		}
		return c, nil
	case (Weapon{}).Type():
		var w *Weapon
		if err := json.Unmarshal(ow.Data, &w); err != nil {
			return nil, err
		}
		return w, nil
	case (Armor{}).Type():
		var a *Armor
		if err := json.Unmarshal(ow.Data, &a); err != nil {
			return nil, err
		}
		return a, nil
	case (Food{}).Type():
		var f *Food
		if err := json.Unmarshal(ow.Data, &f); err != nil {
			return nil, err
		}
		return f, nil
	}
	return nil, errors.New("unknown object type: " + string(ow.Type))
}
