package game

type WeaponType uint8

const (
	WeaponTypeRange WeaponType = iota
	WeaponTypeMelee
	WeaponTypeUnarmed
)

type Weapon struct {
	PrimaryAttribute   Attribute // Primary attribute to draw damage from.
	SecondaryAttribute Attribute // Secondary attribute to draw 50% damage from.
	MinDamage          int       // Character proficiency with a weapon increases min up to max.
	MaxDamage          int
}
