package id

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"
)

// UUID is a UUID. In general, most UUIDs in morogue are UUIDv5s used to identify
// static types of objects, such as archetypes, items, and weapons. These are
// generally acquired through the UID func.
type UUID uuid.UUID

// UnmarshalJSON unmarshals a UUID from data in either UUID form or in a morogue valid "namespace:name" format, such as "morogue:tile:stone-wall".
func (u *UUID) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		var v [16]byte
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		*u = UUID(v)
		return nil
	}
	if len(v) == 0 {
		return nil
	}
	uuid, err := uuid.FromString(v)
	if err == nil {
		*u = UUID(uuid)
	} else {
		strs := strings.SplitN(v, ":", 3)
		if len(strs) < 3 {
			return errors.New("bad namespace length")
		}
		if ns, ok := KeyToNamespace[strs[0]+":"+strs[1]]; !ok {
			return errors.New("no such namespace: " + strs[0] + ":" + strs[1])
		} else {
			uid, err := UID(ns, strs[2])
			if err != nil {
				return errors.Join(err, fmt.Errorf("%s", strs))
			}
			*u = uid
		}
	}
	return nil
}

// IsNil returns true if the UUID is nil.
func (u UUID) IsNil() bool {
	return uuid.UUID(u).IsNil()
}

// Bytes returns a byte slice representation of the UUID.
func (u UUID) Bytes() []byte {
	return u[:]
}

// String returns a string representation of the UUID.
func (u UUID) String() string {
	return uuid.UUID(u).String()
}

// UID generates a unique identifier for the given name in the given morogue namespace. The namespace must be one this is defined in namespaces.
func UID(ns UUID, name string) (UUID, error) {
	if ns != Character && ns != Tile && ns != Mob && ns != Item && ns != Weapon && ns != Armor && ns != Place && ns != Fixture {
		return UUID{}, errors.New("namespace not morogue")
	}
	return UUID(uuid.NewV5(uuid.UUID(ns), name)), nil
}
