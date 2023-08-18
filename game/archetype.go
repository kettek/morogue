package game

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/kettek/morogue/id"
)

// Archetype is an interface that all archetypes must implement.
type Archetype interface {
	GetID() id.UUID
	Type() string
}

type archetypeWrapper struct {
	ID string `json:"id"`
}

func DecodeArchetype(bytes []byte, rootPath string) (Archetype, error) {
	var w archetypeWrapper
	err := json.Unmarshal(bytes, &w)
	if err != nil {
		return nil, err
	}

	parts := strings.SplitAfterN(w.ID, ":", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid archetype id: %s", w.ID)
	}
	key := parts[0] + parts[1][:len(parts[1])-1]
	switch key {
	case id.KeyTile:
		var a TileArchetype
		if err = json.Unmarshal(bytes, &a); err != nil {
			return nil, err
		}
		a.Image = path.Join(rootPath, a.Image)
		return a, nil
	case id.KeyCharacter:
		var a CharacterArchetype
		if err = json.Unmarshal(bytes, &a); err != nil {
			return nil, err
		}
		a.Image = path.Join(rootPath, a.Image)
		a.PlayerOnly = true
		return a, nil
	case id.KeyDoor:
		var a DoorArchetype
		if err = json.Unmarshal(bytes, &a); err != nil {
			return nil, err
		}
		a.Image = path.Join(rootPath, a.Image)
		return a, nil
	case id.KeyMob:
		var a CharacterArchetype
		if err = json.Unmarshal(bytes, &a); err != nil {
			return nil, err
		}
		a.Image = path.Join(rootPath, a.Image)
		return a, nil
	case id.KeyItem:
		var a ItemArchetype
		if err = json.Unmarshal(bytes, &a); err != nil {
			return nil, err
		}
		a.Image = path.Join(rootPath, a.Image)
		return a, nil
	case id.KeyWeapon:
		var a WeaponArchetype
		if err = json.Unmarshal(bytes, &a); err != nil {
			return nil, err
		}
		a.Image = path.Join(rootPath, a.Image)
		return a, nil
	case id.KeyArmor:
		var a ArmorArchetype
		if err = json.Unmarshal(bytes, &a); err != nil {
			return nil, err
		}
		a.Image = path.Join(rootPath, a.Image)
		return a, nil
	default:
		return nil, fmt.Errorf("invalid archetype type: %s", key)
	}
}
