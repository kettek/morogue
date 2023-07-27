package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/morogue/client/ifs"
)

type Binds struct {
	Keys          map[ebiten.Key]Action
	Actions       map[Action][]ebiten.Key
	heldActions   map[Action]int
	multiActions  []MultiAction
	activeActions []Action
}

func (b *Binds) Init() {
	b.Keys = make(map[ebiten.Key]Action)
	b.Actions = make(map[Action][]ebiten.Key)
	b.heldActions = make(map[Action]int)

	// FIXME: Read this in from a binds json.
	b.SetActionKeys("move-left", []ebiten.Key{ebiten.KeyArrowLeft, ebiten.KeyA, ebiten.KeyH})
	b.SetActionKeys("move-right", []ebiten.Key{ebiten.KeyArrowRight, ebiten.KeyD, ebiten.KeyL})
	b.SetActionKeys("move-up", []ebiten.Key{ebiten.KeyArrowUp, ebiten.KeyW, ebiten.KeyK})
	b.SetActionKeys("move-down", []ebiten.Key{ebiten.KeyArrowDown, ebiten.KeyS, ebiten.KeyJ})
	b.SetMultiAction("move-upleft", []Action{"move-left", "move-up"})
	b.SetMultiAction("move-upright", []Action{"move-right", "move-up"})
	b.SetMultiAction("move-downleft", []Action{"move-left", "move-down"})
	b.SetMultiAction("move-downright", []Action{"move-right", "move-down"})
}

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

func (b *Binds) SetActionKeys(a Action, k []ebiten.Key) {
	b.Actions[a] = k
	for _, key := range k {
		b.Keys[key] = a
	}
}

func (b *Binds) SetMultiAction(action Action, actions []Action) {
	b.multiActions = append(b.multiActions, MultiAction{Action: action, Actions: actions})
}

type MultiAction struct {
	Actions []Action
	Action  Action
}

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

type Action string

func (b *Binds) ActiveActions() []Action {
	return b.activeActions
}

func (b *Binds) IsActionActive(a Action) bool {
	for _, act := range b.activeActions {
		if act == a {
			return true
		}
	}
	return false
}

func (b *Binds) IsActionHeld(a Action) int {
	if count, ok := b.heldActions[a]; ok {
		return count
	}
	return -1
}
