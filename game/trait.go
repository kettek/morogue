package game

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

// Trait is an interface that all traits must implement.
type Trait interface {
	String() string
	CanApply(o Object) bool
	AdjustDamages(damages []Damage) []Damage
}

// TraitList is a list of traits.
type TraitList []Trait

// ToStrings converts a TraitList to a list of strings.
func (tl *TraitList) ToStrings() []string {
	var traits []string
	for _, trait := range *tl {
		traits = append(traits, trait.String())
	}
	return traits
}

// UnmarshalJSON unmarshals a list of traits to their actual concrete types, if available.
func (tl *TraitList) UnmarshalJSON(b []byte) error {
	var traits []string
	if err := json.Unmarshal(b, &traits); err != nil {
		return err
	}
	for _, trait := range traits {
		switch trait {
		case TraitKungFu{}.String():
			*tl = append(*tl, TraitKungFu{})
		case TraitClubber{}.String():
			*tl = append(*tl, TraitClubber{})
		default:
			fmt.Println("unknown trait:", trait)
		}
	}
	return nil
}

// MarshalMsgpack marshals a traitlist to a msgpack.
func (tl TraitList) MarshalMsgpack() ([]byte, error) {
	return msgpack.Marshal(tl.ToStrings())
}

// UnmarshalMsgpack unmarshals a traitlist.
func (tl *TraitList) UnmarshalMsgpack(b []byte) error {
	var traits []string
	if err := msgpack.Unmarshal(b, &traits); err != nil {
		return err
	}
	for _, trait := range traits {
		switch trait {
		case TraitKungFu{}.String():
			*tl = append(*tl, TraitKungFu{})
		case TraitClubber{}.String():
			*tl = append(*tl, TraitClubber{})
		}
	}
	return nil
}

// TraitKungFu is the kung fu!
type TraitKungFu struct {
}

// String returns "kung fu".
func (t TraitKungFu) String() string {
	return "kung fu"
}

// CanApply returns true if the trait can be applied to the object.
func (t TraitKungFu) CanApply(o Object) bool {
	return true
}

// AdjustDamages adjusts the provided damages.
func (t TraitKungFu) AdjustDamages(damages []Damage) []Damage {
	for i, damage := range damages {
		if damage.Weapon == WeaponTypeUnarmed {
			damages[i].Min *= 3
			damages[i].Max *= 3
		}
	}
	return damages
}

// TraitClubber is the clubber.
type TraitClubber struct {
}

// String returns "only clubs".
func (t TraitClubber) String() string {
	return "only clubs"
}

// CanApply returns true if the trait can be applied to the object.
func (t TraitClubber) CanApply(o Object) bool {
	switch o := o.(type) {
	case *Weapon:
		title := strings.ToLower(o.Archetype.(WeaponArchetype).Title)
		// FIXME: Use a weapon kind system.
		if !strings.Contains(title, "club") && !strings.Contains(title, "stick") && !strings.Contains(title, "cudgel") && !strings.Contains(title, "cane") {
			return false
		}
	default:
		fmt.Println("unhandled", o)
	}
	return true
}

// AdjustDamages adjusts the provided damages.
func (t TraitClubber) AdjustDamages(damages []Damage) []Damage {
	return damages
}
