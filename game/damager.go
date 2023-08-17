package game

type Damager struct {
	Damages []Damage
}

func (d *Damager) RollDamages() (results []DamageResult) {
	for _, damage := range d.Damages {
		dmg := damage.Roll()
		results = append(results, DamageResult{
			Damage: dmg,
		})
	}
	return
}

func (d *Damager) CalculateFromCharacter(c *Character) {
	d.Damages = []Damage{}
	var mainHand, offHand *Weapon
	var mainType, offType WeaponType
	for _, w := range c.Inventory {
		if w, ok := w.(*Weapon); ok {
			if !w.Applied || w.Archetype == nil {
				continue
			}
			if w.Archetype.(WeaponArchetype).Slots.HasSlot(SlotMainHand) {
				mainHand = w
				mainType = w.Archetype.(WeaponArchetype).WeaponType
			} else if w.Archetype.(WeaponArchetype).Slots.HasSlot(SlotOffHand) {
				offHand = w
				offType = w.Archetype.(WeaponArchetype).WeaponType
			}
		}
	}
	// FIXME: Don't assume Swole, use the weapon's preferred attribute.
	var mainMin, mainMax, mainExtra, offMin, offMax, offExtra int
	if mainHand != nil {
		mainMin = mainHand.Archetype.(WeaponArchetype).MinDamage
		mainMax = mainHand.Archetype.(WeaponArchetype).MaxDamage
		if c.Archetype.(CharacterArchetype).Swole > AttributeLevel(mainMin) {
			if c.Archetype.(CharacterArchetype).Swole > AttributeLevel(mainMax) {
				mainMin = mainMax
				mainExtra = (int(c.Archetype.(CharacterArchetype).Swole) - mainMax) / 2
			} else {
				mainMin = int(c.Archetype.(CharacterArchetype).Swole)
			}
		}
		d.Damages = append(d.Damages, Damage{
			Source:  mainHand.WID,
			Min:     mainMin,
			Max:     mainMax,
			Extra:   mainExtra,
			Reduced: false,
			Weapon:  mainType,
		})
	}
	if offHand != nil {
		offMin = offHand.Archetype.(WeaponArchetype).MinDamage
		offMax = offHand.Archetype.(WeaponArchetype).MaxDamage
		if c.Archetype.(CharacterArchetype).Swole > AttributeLevel(offMin) {
			if c.Archetype.(CharacterArchetype).Swole > AttributeLevel(offMax) {
				offMin = offMax
				offExtra = (int(c.Archetype.(CharacterArchetype).Swole) - offMax) / 2
			} else {
				offMin = int(c.Archetype.(CharacterArchetype).Swole)
			}
		}
		d.Damages = append(d.Damages, Damage{
			Source:  offHand.WID,
			Min:     offMin,
			Max:     offMax,
			Extra:   offExtra,
			Reduced: true,
			Weapon:  offType,
		})
	}
	if mainHand == nil && offHand == nil {
		d.Damages = append(d.Damages, Damage{
			Source:  c.WID,
			Min:     0,
			Max:     int(c.Archetype.(CharacterArchetype).Swole) / 2,
			Extra:   0,
			Reduced: true,
			Weapon:  WeaponTypeUnarmed,
		})
	}

	for _, t := range c.Archetype.(CharacterArchetype).Traits {
		d.Damages = t.AdjustDamages(d.Damages)
	}
}
