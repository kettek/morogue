package net

import (
	"encoding/json"

	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/id"
)

type Wrapper struct {
	Type string          `json:"t"`
	Data json.RawMessage `json:"d"`
}

func (w *Wrapper) Message() Message {
	switch w.Type {
	case (PingMessage{}).Type():
		var m PingMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (LoginMessage{}).Type():
		var m LoginMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (RegisterMessage{}).Type():
		var m RegisterMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (LogoutMessage{}).Type():
		var m LogoutMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (ArchetypesMessage{}).Type():
		var m ArchetypesMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (CharactersMessage{}).Type():
		var m CharactersMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (CreateCharacterMessage{}).Type():
		var m CreateCharacterMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (DeleteCharacterMessage{}).Type():
		var m DeleteCharacterMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (JoinCharacterMessage{}).Type():
		var m JoinCharacterMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (UnjoinCharacterMessage{}).Type():
		var m UnjoinCharacterMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (WorldsMessage{}).Type():
		var m WorldsMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (CreateWorldMessage{}).Type():
		var m CreateWorldMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (JoinWorldMessage{}).Type():
		var m JoinWorldMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (LocationMessage{}).Type():
		var m LocationMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (TileMessage{}).Type():
		var m TileMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (DesireMessage{}).Type():
		var m DesireMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (OwnerMessage{}).Type():
		var m OwnerMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (InventoryMessage{}).Type():
		var m InventoryMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (SkillsMessage{}).Type():
		var m SkillsMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (EventMessage{}).Type():
		var m EventMessage
		json.Unmarshal(w.Data, &m)
		return m
	case (EventsMessage{}).Type():
		var m EventsMessage
		json.Unmarshal(w.Data, &m)
		return m

	}
	return nil
}

type Message interface {
	Type() string
}

type PingMessage struct {
}

func (m PingMessage) Type() string {
	return "ping"
}

type LoginMessage struct {
	User       string `json:"u,omitempty"`
	Password   string `json:"p,omitempty"`
	Result     string `json:"r,omitempty"`
	ResultCode int    `json:"c,omitempty"`
}

func (m LoginMessage) Type() string {
	return "login"
}

type RegisterMessage struct {
	User       string `json:"u,omitempty"`
	Password   string `json:"p,omitempty"`
	Result     string `json:"r,omitempty"`
	ResultCode int    `json:"c,omitempty"`
}

func (m RegisterMessage) Type() string {
	return "register"
}

type LogoutMessage struct {
	Result     string `json:"r,omitempty"`
	ResultCode int    `json:"c,omitempty"`
}

func (m LogoutMessage) Type() string {
	return "logout"
}

type ArchetypesMessage struct {
	Archetypes []game.Archetype `json:"a,omitempty"`
}

// archetypeWrapper is used to wrap a game.Archetype interface for safe traversal.
type archetypeWrapper struct {
	Type string          `json:"t"`
	Data json.RawMessage `json:"d"`
}

// MarshalJSON marshals a game.Archetype interface into a wrapper JSON object.
func (m *ArchetypesMessage) UnmarshalJSON(data []byte) error {
	var archetypes []archetypeWrapper
	if err := json.Unmarshal(data, &archetypes); err != nil {
		return err
	}

	for _, a := range archetypes {
		switch a.Type {
		case (game.CharacterArchetype{}).Type():
			var archetype game.CharacterArchetype
			json.Unmarshal(a.Data, &archetype)
			m.Archetypes = append(m.Archetypes, archetype)
		case (game.ItemArchetype{}).Type():
			var archetype game.ItemArchetype
			json.Unmarshal(a.Data, &archetype)
			m.Archetypes = append(m.Archetypes, archetype)
		case (game.WeaponArchetype{}).Type():
			var archetype game.WeaponArchetype
			json.Unmarshal(a.Data, &archetype)
			m.Archetypes = append(m.Archetypes, archetype)
		case (game.ArmorArchetype{}).Type():
			var archetype game.ArmorArchetype
			json.Unmarshal(a.Data, &archetype)
			m.Archetypes = append(m.Archetypes, archetype)
		}
	}

	return nil
}

// UnmarshalJSON unmarshals a wrapper JSON object into a game.Archetype interface.
func (m ArchetypesMessage) MarshalJSON() ([]byte, error) {
	var archetypes []archetypeWrapper
	for _, a := range m.Archetypes {
		aw, err := json.Marshal(a)
		if err != nil {
			panic(err)
		}
		archetypes = append(archetypes, archetypeWrapper{
			Type: a.Type(),
			Data: aw,
		})
	}
	return json.Marshal(archetypes)
}

func (m ArchetypesMessage) Type() string {
	return "archetypes"
}

type CharactersMessage struct {
	Characters []*game.Character `json:"c,omitempty"`
}

func (m CharactersMessage) Type() string {
	return "characters"
}

type CreateCharacterMessage struct {
	Result     string  `json:"r,omitempty"`
	ResultCode int     `json:"c,omitempty"`
	Name       string  `json:"n,omitempty"`
	Archetype  id.UUID `json:"a,omitempty"`
}

func (m CreateCharacterMessage) Type() string {
	return "create-character"
}

type DeleteCharacterMessage struct {
	Result     string `json:"r,omitempty"`
	ResultCode int    `json:"c,omitempty"`
	Name       string `json:"n,omitempty"`
}

func (m DeleteCharacterMessage) Type() string {
	return "delete-character"
}

type JoinCharacterMessage struct {
	Result     string `json:"r,omitempty"`
	ResultCode int    `json:"c,omitempty"`
	Name       string `json:"n,omitempty"`
}

func (m JoinCharacterMessage) Type() string {
	return "join-character"
}

type UnjoinCharacterMessage struct {
	Result     string `json:"r,omitempty"`
	ResultCode int    `json:"c,omitempty"`
}

func (m UnjoinCharacterMessage) Type() string {
	return "unjoin-character"
}

type WorldsMessage struct {
	Result     string           `json:"r,omitempty"`
	ResultCode int              `json:"c,omitempty"`
	Worlds     []game.WorldInfo `json:"w,omitempty"`
}

func (m WorldsMessage) Type() string {
	return "worlds"
}

type CreateWorldMessage struct {
	Result     string `json:"r,omitempty"`
	ResultCode int    `json:"c,omitempty"`
	Name       string `json:"n,omitempty"`
	Password   string `json:"p,omitempty"`
}

func (m CreateWorldMessage) Type() string {
	return "create-world"
}

type JoinWorldMessage struct {
	Result     string  `json:"r,omitempty"`
	ResultCode int     `json:"c,omitempty"`
	World      id.UUID `json:"w,omitempty"`
	Password   string  `json:"p,omitempty"`
}

func (m JoinWorldMessage) Type() string {
	return "join-world"
}

type LocationMessage struct {
	Result     string       `json:"r,omitempty"`
	ResultCode int          `json:"c,omitempty"`
	ID         id.UUID      `json:"id,omitempty"`
	Objects    game.Objects `json:"o,omitempty"`
	Cells      game.Cells   `json:"g,omitempty"`
}

func (m LocationMessage) Type() string {
	return "location"
}

type TileMessage struct {
	Result     string             `json:"r,omitempty"`
	ResultCode int                `json:"c,omitempty"`
	ID         id.UUID            `json:"id,omitempty"`
	Tile       game.TileArchetype `sjon:"t,omitempty"`
}

func (m TileMessage) Type() string {
	return "tile"
}

type OwnerMessage struct {
	WID       id.WID        `json:"wid,omitempty"`
	Inventory []game.Object `json:"i,omitempty"`
	Skills    game.Skills   `json:"s,omitempty"`
}

func (m OwnerMessage) Type() string {
	return "owner"
}

type InventoryMessage struct {
	WID       id.WID        `json:"wid,omitempty"`
	Inventory []game.Object `json:"i,omitempty"`
}

func (m InventoryMessage) Type() string {
	return "inventory"
}

type SkillsMessage struct {
	WID    id.WID      `json:"wid,omitempty"`
	Skills game.Skills `json:"s,omitempty"`
}

func (m SkillsMessage) Type() string {
	return "skills"
}

type DesireMessage struct {
	Result     string             `json:"r,omitempty"`
	ResultCode int                `json:"c,omitempty"`
	WID        id.WID             `json:"wid,omitempty"`
	Desire     game.DesireWrapper `json:"d,omitempty"`
}

func (m DesireMessage) Type() string {
	return "desire"
}

type EventMessage struct {
	Event game.EventWrapper `json:"e,omitempty"`
}

func (m EventMessage) Type() string {
	return "event"
}

type EventsMessage struct {
	Events []game.EventWrapper `json:"e,omitempty"`
}

func (m EventsMessage) Type() string {
	return "events"
}
