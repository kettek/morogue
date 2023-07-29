package game

import (
	"encoding/json"
	"errors"

	"github.com/kettek/morogue/id"
)

// ObjectType is a key that represents an object's type. This is used
// for safely marshalling/unmarshalling the Object interface.
type ObjectType string

// Object is an interface intended for location objects. This includes
// players, items, and enemies.
type Object interface {
	Type() ObjectType
	GetWID() id.WID
}

// ObjectWrapper wraps an Object interface for json marshal and unmarshal.
type ObjectWrapper struct {
	Type ObjectType      `json:"t"`
	Data json.RawMessage `json:"d"`
}

// Object returns the wrapped Object. To add additional types, they must be
// handled here.
func (ow ObjectWrapper) Object() (Object, error) {
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
	}
	return nil, errors.New("unknown object type: " + string(ow.Type))
}

// ObjectsWrapper is a slice of ObjectWrappers.
type ObjectsWrapper []ObjectWrapper
