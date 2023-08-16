package game

import (
	"fmt"
	"math/rand"

	"github.com/kettek/morogue/id"
)

// Damage represents the damage range of an attack.
type Damage struct {
	Source   id.WID
	Min, Max int
	Extra    int
	Reduced  bool
	Weapon   WeaponType
}

// RangeString returns a string representation of the damage range.
func (d Damage) RangeString() string {
	var s string
	if d.Min == 0 {
		s = fmt.Sprintf("〜%d", d.Max)
	} else if d.Min == d.Max {
		s = fmt.Sprintf("%d", d.Min)
	} else {
		s = fmt.Sprintf("%d〜%d", d.Min, d.Max)
	}
	if d.Extra > 0 {
		s += fmt.Sprintf(" +%d", d.Extra)
	}
	return s
}

// Roll rolls the damage range and returns the result.
func (d Damage) Roll() int {
	return rand.Intn(d.Max-d.Min+1) + d.Min + d.Extra
}

// DamageResult represents the result of a damage roll.
type DamageResult struct {
	Damage int
}
