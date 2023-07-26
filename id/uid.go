package id

import (
	"errors"

	"github.com/gofrs/uuid/v5"
)

// UUID is a UUID. In general, most UUIDs in morogue are UUIDv5s used to identify
// static types of objects, such as archetypes, items, and weapons. These are
// generally acquired through the UID func.
type UUID = uuid.UUID

// UID generates a unique identifier for the given name in the given morogue namespace. The namespace must be NSArchetype, NSMob, or NSItem.
func UID(ns uuid.UUID, name string) (uuid.UUID, error) {
	if ns != Archetype && ns != Tile && ns != Mob && ns != Item {
		return uuid.UUID{}, errors.New("namespace not morogue")
	}
	return uuid.NewV5(ns, name), nil
}
