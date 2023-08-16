package game

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

type Trait interface {
	String() string
	CanApply(o Object) bool
	AdjustDamages(damages []Damage) []Damage
}

type TraitList []Trait

func (tl *TraitList) ToStrings() []string {
	var traits []string
	for _, trait := range *tl {
		traits = append(traits, trait.String())
	}
	return traits
}

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

func (tl TraitList) MarshalMsgpack() ([]byte, error) {
	return msgpack.Marshal(tl.ToStrings())
}

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

type TraitKungFu struct {
}

func (t TraitKungFu) String() string {
	return "kung fu"
}

func (t TraitKungFu) CanApply(o Object) bool {
	return true
}

func (t TraitKungFu) AdjustDamages(damages []Damage) []Damage {
	for i, damage := range damages {
		if damage.Weapon == WeaponTypeUnarmed {
			damages[i].Min *= 3
			damages[i].Max *= 3
		}
	}
	return damages
}

type TraitClubber struct {
}

func (t TraitClubber) String() string {
	return "only clubs"
}

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

func (t TraitClubber) AdjustDamages(damages []Damage) []Damage {
	return damages
}
