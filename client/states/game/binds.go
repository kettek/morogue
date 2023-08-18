package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/morogue/client/ifs"
)

// Binds provides a way to map input to game actions. Keys can be bound to
// actions that then allow the game state to determine if a particular
// action is active or not. Additionally, "multi actions" can be triggered
// when one or more actions are fulfilled.
type Binds struct {
	Keys          map[ebiten.Key]Action
	Actions       map[Action][]ebiten.Key
	heldActions   map[Action]int
	multiActions  []MultiAction
	activeActions []Action
}

// Init initializes the binds structures and sets up some default binds.
func (b *Binds) Init() {
	b.Keys = make(map[ebiten.Key]Action)
	b.Actions = make(map[Action][]ebiten.Key)
	b.heldActions = make(map[Action]int)

	// FIXME: Read this in from a binds json.
	b.SetActionKeys("move-left", []ebiten.Key{ebiten.KeyArrowLeft, ebiten.KeyA, ebiten.KeyH})
	b.SetActionKeys("move-right", []ebiten.Key{ebiten.KeyArrowRight, ebiten.KeyD, ebiten.KeyL})
	b.SetActionKeys("move-up", []ebiten.Key{ebiten.KeyArrowUp, ebiten.KeyW, ebiten.KeyK})
	b.SetActionKeys("move-down", []ebiten.Key{ebiten.KeyArrowDown, ebiten.KeyS, ebiten.KeyJ})
	b.SetActionKeys("bash", []ebiten.Key{ebiten.KeyB})
	b.SetActionKeys("pickup", []ebiten.Key{ebiten.KeyComma})
	b.SetActionKeys("lock-camera", []ebiten.Key{ebiten.KeyC})
	b.SetActionKeys("snap-camera", []ebiten.Key{ebiten.KeySpace})
	b.SetMultiAction("move-upleft", []Action{"move-left", "move-up"})
	b.SetMultiAction("move-upright", []Action{"move-right", "move-up"})
	b.SetMultiAction("move-downleft", []Action{"move-left", "move-down"})
	b.SetMultiAction("move-downright", []Action{"move-right", "move-down"})
}

// Update is called per tick and synchronizes inputs to actions.
func (b *Binds) Update(ctx ifs.RunContext) {
	for action := range b.heldActions {
		held := false
		for _, act := range b.activeActions {
			if act == action {
				held = true
				break
			}
		}
		if !held {
			delete(b.heldActions, action)
		}
	}

	b.activeActions = b.activeActions[:0]
	for key, action := range b.Keys {
		if ebiten.IsKeyPressed(key) {
			b.activeActions = append(b.activeActions, action)
		}
	}

	for _, ma := range b.multiActions {
		if ma.HasActions(b.activeActions) {
			b.activeActions = append(b.activeActions, ma.Action)
		}
	}

	for _, action := range b.activeActions {
		if _, ok := b.heldActions[action]; ok {
			b.heldActions[action]++
		} else {
			b.heldActions[action] = 0
		}
	}
}

// SetActionKey sets an action associated with a key.
func (b *Binds) SetActionKeys(a Action, k []ebiten.Key) {
	b.Actions[a] = k
	for _, key := range k {
		b.Keys[key] = a
	}
}

// SetMultiAction sets an action associated with multiple actions.
func (b *Binds) SetMultiAction(action Action, actions []Action) {
	b.multiActions = append(b.multiActions, MultiAction{Action: action, Actions: actions})
}

// MultiAction maps one or more actions to represent another action. For example,
// if the two actions "move-left" and "move-up" are active, a multi action could
// convert that into a "move-up-left" action.
type MultiAction struct {
	Actions []Action
	Action  Action
}

// HasActions returns if the passed in actions fulfill the multi action's actions.
func (ma *MultiAction) HasActions(actions []Action) bool {
	for _, a1 := range ma.Actions {
		has := false
		for _, a2 := range actions {
			if a1 == a2 {
				has = true
				break
			}
		}
		if !has {
			return false
		}
	}
	return true
}

// Action represents a given action, such as "move-left".
type Action string

// ActiveActions returns the current active actions.
func (b *Binds) ActiveActions() []Action {
	return b.activeActions
}

// IsActionActive returns if the given action is active.
func (b *Binds) IsActionActive(a Action) bool {
	for _, act := range b.activeActions {
		if act == a {
			return true
		}
	}
	return false
}

// IsActionHeld returns how many ticks an action has been held. It returns -1
// if the action is not held. It may return 0 on the first tick of being held.
func (b *Binds) IsActionHeld(a Action) int {
	if count, ok := b.heldActions[a]; ok {
		return count
	}
	return -1
}
