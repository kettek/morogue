package game

import "github.com/kettek/morogue/id"

// ArmorType is the type the armor is considered as.
type ArmorType uint8

const (
	ArmorTypeNone ArmorType = iota
	ArmorTypeLight
	ArmorTypeMedium
	ArmorTypeHeavy
)

// Armor is a weapon.
type Armor struct {
	WID      id.WID
	MinArmor int  `json:"m,omitempty"` // Character proficiency with a weapon increases min up to max.
	MaxArmor int  `json:"M,omitempty"`
	Applied  bool `json:"a,omitempty"`
}

func (a Armor) Type() ObjectType {
	return "armor"
}

func (a Armor) GetWID() id.WID {
	return a.WID
}
