package game

import (
	"errors"
	"fmt"
)

// Slot is a slot used for equipment.
type Slot string

const (
	SlotNone     Slot = ""
	SlotHead     Slot = "head"
	SlotTorso    Slot = "torso"
	SlotArms     Slot = "arms"
	SlotHands    Slot = "hands"
	SlotOther    Slot = "other"
	SlotLegs     Slot = "legs"
	SlotFeet     Slot = "feet"
	SlotMainHand Slot = "main-hand"
	SlotOffHand  Slot = "off-hand"
)

// Slots is a slice of slots used for equipment.
type Slots []Slot

// ToMap converts the slice of slots to a map of slots.
func (s Slots) ToMap() SlotMap {
	m := make(SlotMap)
	for _, slot := range s {
		m[slot] = false
	}
	return m
}

// String returns a string representation of the slice of slots, separated by " & ".
func (s Slots) String() string {
	var str string
	for i, slot := range s {
		str += string(slot)
		if i != len(s)-1 {
			str += " & "
		}
	}
	return str
}

// HasSlot returns true if the slice of slots has the given slot.
func (s Slots) HasSlot(slot Slot) bool {
	for _, s := range s {
		if s == slot {
			return true
		}
	}
	return false
}

// SlotMap is a map of slots used by characters.
type SlotMap map[Slot]bool

// HasSlot returns true if the slot map has the given slot.
func (s SlotMap) HasSlot(slot Slot) bool {
	if _, ok := s[slot]; ok {
		return true
	}
	return false
}

// HasSlots returns nil if the slot map has all of the given slots. Returns ErrMissingSlots wrapping a list of missing slots otherwise.
func (s SlotMap) HasSlots(slots Slots) error {
	var err []error
	for _, name := range slots {
		if !s.HasSlot(name) {
			err = append(err, fmt.Errorf("%s", name))
		}
	}
	if len(err) > 0 {
		return errors.Join(append([]error{ErrMissingSlots}, err...)...)
	}
	return nil
}

// AreSlotsOpen returns nil if all of the slots in the slot map are open. Returns ErrUsedSlots wrapping a list of used slots otherwise.
func (s SlotMap) AreSlotsOpen(slots Slots) error {
	var err []error
	for _, name := range slots {
		if slot, ok := s[name]; ok && slot {
			err = append(err, fmt.Errorf("%s", name))
		}
	}
	if len(err) > 0 {
		return errors.Join(append([]error{ErrUsedSlots}, err...)...)
	}
	return nil
}

// Apply adds the given slots to the slot map. If any slots are missing, an error will instead be returned.
func (s SlotMap) Apply(slots Slots) error {
	if err := s.HasSlots(slots); err != nil {
		return err
	}
	if err := s.AreSlotsOpen(slots); err != nil {
		return err
	}
	for _, slot := range slots {
		s[slot] = true
	}
	return nil
}

// Unapply removes the given slots from the slot map. Even if slots are missing, slots that are not missing will be removed.
func (s SlotMap) Unapply(slots Slots) error {
	var errs []error
	for _, slot := range slots {
		if _, ok := s[slot]; !ok {
			errs = append(errs, ErrMissingSlot)
		} else {
			s[slot] = false
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

var (
	ErrMissingSlot  = errors.New("missing slot")
	ErrMissingSlots = errors.New("missing slots")
	ErrUsedSlots    = errors.New("used slots")
)
