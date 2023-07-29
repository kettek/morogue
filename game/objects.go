package game

import (
	"encoding/json"
	"errors"

	"github.com/kettek/morogue/id"
)

// Objects is a slice of Object interfaces. This type wrapper provides
// various convenience functions for accessing and modifying the slice.
type Objects []Object

// Add adds an Object to the slice.
func (o *Objects) Add(obj Object) {
	*o = append(*o, obj)
}

// RemoveByWID removes a stored object by its WID.
func (o *Objects) RemoveByWID(wid id.WID) Object {
	for i, obj := range *o {
		if obj.GetWID() == wid {
			*o = append((*o)[:i], (*o)[i+1:]...)
			return obj
		}
	}
	return nil
}

// MarshalJSON returns bytes as an ObjectsWrapper.
func (o Objects) MarshalJSON() ([]byte, error) {
	var ow ObjectsWrapper

	for _, obj := range o {
		objBytes, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}
		ow = append(ow, ObjectWrapper{
			Type: obj.Type(),
			Data: objBytes,
		})
	}
	return json.Marshal(ow)
}

// UnmarshalJSON unmarshals the given bytes as an ObjectsWrapper and
// appends objects into the slice.
func (o *Objects) UnmarshalJSON(b []byte) error {
	var osw ObjectsWrapper
	if err := json.Unmarshal(b, &osw); err != nil {
		return err
	}

	var errs []error
	for _, ow := range osw {
		if obj, err := ow.Object(); err != nil {
			errs = append(errs, err)
		} else {
			*o = append(*o, obj)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
