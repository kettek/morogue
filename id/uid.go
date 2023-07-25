package id

import (
	"errors"

	"github.com/gofrs/uuid/v5"
)

type UUID = uuid.UUID

// UID generates a unique identifier for the given name in the given morogue namespace. The namespace must be NSArchetype, NSMob, or NSItem.
func UID(ns uuid.UUID, name string) (uuid.UUID, error) {
	if ns != Archetype && ns != Mob && ns != Item {
		return uuid.UUID{}, errors.New("namespace not morogue")
	}
	return uuid.NewV5(ns, name), nil
}
