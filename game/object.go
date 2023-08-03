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
}

func CreateObjectFromArchetype(a Archetype) Object {
	switch a.(type) {
	case CharacterArchetype:
		return &Character{
			Archetype: a.GetID(),
		}
	case WeaponArchetype:
		return &Weapon{
			Archetype: a.GetID(),
		}
	case ArmorArchetype:
		return &Armor{
			Archetype: a.GetID(),
		}
	}
	return nil
}

type RawMessage []byte

func (m RawMessage) MarshalMsgpack() ([]byte, error) {
	return msgpack.Marshal((msgpack.RawMessage)(m))
}

func (m *RawMessage) UnmarshalMsgpack(b []byte) error {
	return msgpack.Unmarshal(b, (*msgpack.RawMessage)(m))
}

func (m RawMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal((json.RawMessage)(m))
}

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
	}
	return nil, errors.New("unknown object type: " + string(ow.Type))
}
