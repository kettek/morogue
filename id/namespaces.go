package id

import (
	"crypto/sha1"

	"github.com/gofrs/uuid/v5"
)

const (
	KeyCharacter = "morogue:character"
	KeyTile      = "morogue:tile"
	KeyDoor      = "morogue:door"
	KeyMob       = "morogue:mob"
	KeyItem      = "morogue:item"
	KeyWeapon    = "morogue:weapon"
	KeyArmor     = "morogue:armor"
	//
	KeyPlace   = "morogue:place"
	KeyFixture = "morogue:fixture"
)

var (
	Character UUID
	Tile      UUID
	Door      UUID
	Mob       UUID
	Item      UUID
	Weapon    UUID
	Armor     UUID
	//
	Place   UUID
	Fixture UUID
)

// NamespaceToKey provides a mapping of morogue's UUIDv5s to their string keys.
var NamespaceToKey map[UUID]string
var KeyToNamespace map[string]UUID

func init() {
	NamespaceToKey = make(map[UUID]string)
	KeyToNamespace = make(map[string]UUID)
	{
		hasher := sha1.New()
		hasher.Write([]byte(KeyCharacter))
		sha := hasher.Sum(nil)

		Character = UUID(uuid.Must(uuid.FromBytes(sha[:16])))
		NamespaceToKey[Character] = KeyCharacter
		KeyToNamespace[KeyCharacter] = Character
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
		hasher.Write([]byte(KeyDoor))
		sha := hasher.Sum(nil)

		Door = UUID(uuid.Must(uuid.FromBytes(sha[:16])))
		NamespaceToKey[Door] = KeyDoor
		KeyToNamespace[KeyDoor] = Door
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
	{
		hasher := sha1.New()
		hasher.Write([]byte(KeyWeapon))
		sha := hasher.Sum(nil)

		Weapon = UUID(uuid.Must(uuid.FromBytes(sha[:16])))
		NamespaceToKey[Weapon] = KeyWeapon
		KeyToNamespace[KeyWeapon] = Weapon
	}
	{
		hasher := sha1.New()
		hasher.Write([]byte(KeyArmor))
		sha := hasher.Sum(nil)

		Armor = UUID(uuid.Must(uuid.FromBytes(sha[:16])))
		NamespaceToKey[Armor] = KeyArmor
		KeyToNamespace[KeyArmor] = Armor
	}
	//
	{
		hasher := sha1.New()
		hasher.Write([]byte(KeyPlace))
		sha := hasher.Sum(nil)

		Place = UUID(uuid.Must(uuid.FromBytes(sha[:16])))
		NamespaceToKey[Place] = KeyPlace
		KeyToNamespace[KeyPlace] = Place
	}
	{
		hasher := sha1.New()
		hasher.Write([]byte(KeyFixture))
		sha := hasher.Sum(nil)

		Fixture = UUID(uuid.Must(uuid.FromBytes(sha[:16])))
		NamespaceToKey[Fixture] = KeyFixture
		KeyToNamespace[KeyFixture] = Fixture
	}
}
