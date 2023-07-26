package id

import (
	"crypto/sha1"

	"github.com/gofrs/uuid/v5"
)

const (
	KeyArchetype = "morogue:archetype"
	KeyMob       = "morogue:mob"
	KeyItem      = "morogue:item"
)

var (
	Archetype uuid.UUID
	Mob       uuid.UUID
	Item      uuid.UUID
)

// NamespaceToKey provides a mapping of morogue's UUIDv5s to their string keys.
var NamespaceToKey map[uuid.UUID]string

func init() {
	NamespaceToKey = make(map[uuid.UUID]string)
	{
		hasher := sha1.New()
		hasher.Write([]byte(KeyArchetype))
		sha := hasher.Sum(nil)

		Archetype, _ = uuid.FromBytes(sha[:16])
		NamespaceToKey[Archetype] = KeyArchetype
	}
	{
		hasher := sha1.New()
		hasher.Write([]byte(KeyMob))
		sha := hasher.Sum(nil)

		Mob, _ = uuid.FromBytes(sha[:16])
		NamespaceToKey[Mob] = KeyMob
	}
	{
		hasher := sha1.New()
		hasher.Write([]byte(KeyItem))
		sha := hasher.Sum(nil)

		Item, _ = uuid.FromBytes(sha[:16])
		NamespaceToKey[Item] = KeyItem
	}
}
