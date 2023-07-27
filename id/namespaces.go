package id

import (
	"crypto/sha1"

	"github.com/gofrs/uuid/v5"
)

const (
	KeyArchetype = "morogue:archetype"
	KeyTile      = "morogue:tile"
	KeyMob       = "morogue:mob"
	KeyItem      = "morogue:item"
)

var (
	Archetype UUID
	Tile      UUID
	Mob       UUID
	Item      UUID
)

// NamespaceToKey provides a mapping of morogue's UUIDv5s to their string keys.
var NamespaceToKey map[UUID]string
var KeyToNamespace map[string]UUID

func init() {
	NamespaceToKey = make(map[UUID]string)
	KeyToNamespace = make(map[string]UUID)
	{
		hasher := sha1.New()
		hasher.Write([]byte(KeyArchetype))
		sha := hasher.Sum(nil)

		Archetype = UUID(uuid.Must(uuid.FromBytes(sha[:16])))
		NamespaceToKey[Archetype] = KeyArchetype
		KeyToNamespace[KeyArchetype] = Archetype
	}
	{
		hasher := sha1.New()
		hasher.Write([]byte(KeyTile))
		sha := hasher.Sum(nil)

		Tile = UUID(uuid.Must(uuid.FromBytes(sha[:16])))
		NamespaceToKey[Tile] = KeyTile
		KeyToNamespace[KeyTile] = Tile
	}
	{
		hasher := sha1.New()
		hasher.Write([]byte(KeyMob))
		sha := hasher.Sum(nil)

		Mob = UUID(uuid.Must(uuid.FromBytes(sha[:16])))
		NamespaceToKey[Mob] = KeyMob
		KeyToNamespace[KeyMob] = Mob
	}
	{
		hasher := sha1.New()
		hasher.Write([]byte(KeyItem))
		sha := hasher.Sum(nil)

		Item = UUID(uuid.Must(uuid.FromBytes(sha[:16])))
		NamespaceToKey[Item] = KeyItem
		KeyToNamespace[KeyItem] = Item
	}
}
