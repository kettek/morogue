package game

import "errors"

// Slot is a slot used for equipment.
type Slot string

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

// SlotMap is a map of slots used by characters.
type SlotMap map[Slot]bool

// HasSlot returns true if the slot map has the given slot.
func (s SlotMap) HasSlot(slot Slot) bool {
	if _, ok := s[slot]; ok {
		return true
	}
	return false
}

// HasSlots returns true if the slot map has all of the given slots.
func (s SlotMap) HasSlots(slots Slots) bool {
	for _, slot := range slots {
		if !s.HasSlot(slot) {
			return false
		}
	}
	return true
}

// Apply adds the given slots to the slot map. If any slots are missing, an error will instead be returned.
func (s SlotMap) Apply(slots Slots) error {
	if !s.HasSlots(slots) {
		return ErrMissingSlots
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
)
