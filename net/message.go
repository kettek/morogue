package net

import (
	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/id"
	"github.com/vmihailenco/msgpack/v5"
)

type Wrapper struct {
	Type string             `msgpack:"t"`
	Data msgpack.RawMessage `msgpack:"d"`
}

func (w *Wrapper) Message() Message {
	switch w.Type {
	case (PingMessage{}).Type():
		var m PingMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (LoginMessage{}).Type():
		var m LoginMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (RegisterMessage{}).Type():
		var m RegisterMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (LogoutMessage{}).Type():
		var m LogoutMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (ArchetypesMessage{}).Type():
		var m ArchetypesMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (CharactersMessage{}).Type():
		var m CharactersMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (CreateCharacterMessage{}).Type():
		var m CreateCharacterMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (DeleteCharacterMessage{}).Type():
		var m DeleteCharacterMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (JoinCharacterMessage{}).Type():
		var m JoinCharacterMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (UnjoinCharacterMessage{}).Type():
		var m UnjoinCharacterMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (WorldsMessage{}).Type():
		var m WorldsMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (CreateWorldMessage{}).Type():
		var m CreateWorldMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (JoinWorldMessage{}).Type():
		var m JoinWorldMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (LocationMessage{}).Type():
		var m LocationMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (TileMessage{}).Type():
		var m TileMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (DesireMessage{}).Type():
		var m DesireMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (OwnerMessage{}).Type():
		var m OwnerMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (InventoryMessage{}).Type():
		var m InventoryMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (SkillsMessage{}).Type():
		var m SkillsMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (EventMessage{}).Type():
		var m EventMessage
		msgpack.Unmarshal(w.Data, &m)
		return m
	case (EventsMessage{}).Type():
		var m EventsMessage
		msgpack.Unmarshal(w.Data, &m)
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
	User       string `msgpack:"u,omitempty"`
	Password   string `msgpack:"p,omitempty"`
	Result     string `msgpack:"r,omitempty"`
	ResultCode int    `msgpack:"c,omitempty"`
}

func (m LoginMessage) Type() string {
	return "login"
}

type RegisterMessage struct {
	User       string `msgpack:"u,omitempty"`
	Password   string `msgpack:"p,omitempty"`
	Result     string `msgpack:"r,omitempty"`
	ResultCode int    `msgpack:"c,omitempty"`
}

func (m RegisterMessage) Type() string {
	return "register"
}

type LogoutMessage struct {
	Result     string `msgpack:"r,omitempty"`
	ResultCode int    `msgpack:"c,omitempty"`
}

func (m LogoutMessage) Type() string {
	return "logout"
}

type ArchetypesMessage struct {
	Archetypes []game.Archetype `msgpack:"a,omitempty"`
}

// archetypeWrapper is used to wrap a game.Archetype interface for safe traversal.
type archetypeWrapper struct {
	Type string             `msgpack:"t"`
	Data msgpack.RawMessage `msgpack:"d"`
}

// UnmarshalMsgpack unmarshals a wrapper Msgpack object into a game.Archetype interface.
func (m *ArchetypesMessage) UnmarshalMsgpack(data []byte) error {
	var archetypes []archetypeWrapper
	if err := msgpack.Unmarshal(data, &archetypes); err != nil {
		return err
	}

	for _, a := range archetypes {
		switch a.Type {
		case (game.CharacterArchetype{}).Type():
			var archetype game.CharacterArchetype
			if err := msgpack.Unmarshal(a.Data, &archetype); err != nil {
				panic(err)
			}
			m.Archetypes = append(m.Archetypes, archetype)
		case (game.ItemArchetype{}).Type():
			var archetype game.ItemArchetype
			if err := msgpack.Unmarshal(a.Data, &archetype); err != nil {
				panic(err)
			}
			m.Archetypes = append(m.Archetypes, archetype)
		case (game.WeaponArchetype{}).Type():
			var archetype game.WeaponArchetype
			if err := msgpack.Unmarshal(a.Data, &archetype); err != nil {
				panic(err)
			}
			m.Archetypes = append(m.Archetypes, archetype)
		case (game.ArmorArchetype{}).Type():
			var archetype game.ArmorArchetype
			if err := msgpack.Unmarshal(a.Data, &archetype); err != nil {
				panic(err)
			}
			m.Archetypes = append(m.Archetypes, archetype)
		}
	}

	return nil
}

// MarshalMsgpack marshals a game.Archetype interface into a wrapper Msgpack object.
func (m ArchetypesMessage) MarshalMsgpack() ([]byte, error) {
	var archetypes []archetypeWrapper
	for _, a := range m.Archetypes {
		aw, err := msgpack.Marshal(a)
		if err != nil {
			panic(err)
		}
		archetypes = append(archetypes, archetypeWrapper{
			Type: a.Type(),
			Data: aw,
		})
	}

	return msgpack.Marshal(archetypes)
}

func (m ArchetypesMessage) Type() string {
	return "archetypes"
}

type CharactersMessage struct {
	Characters []*game.Character `msgpack:"c,omitempty"`
}

func (m CharactersMessage) Type() string {
	return "characters"
}

type CreateCharacterMessage struct {
	Result     string  `msgpack:"r,omitempty"`
	ResultCode int     `msgpack:"c,omitempty"`
	Name       string  `msgpack:"n,omitempty"`
	Archetype  id.UUID `msgpack:"a,omitempty"`
}

func (m CreateCharacterMessage) Type() string {
	return "create-character"
}

type DeleteCharacterMessage struct {
	Result     string `msgpack:"r,omitempty"`
	ResultCode int    `msgpack:"c,omitempty"`
	Name       string `msgpack:"n,omitempty"`
}

func (m DeleteCharacterMessage) Type() string {
	return "delete-character"
}

type JoinCharacterMessage struct {
	Result     string `msgpack:"r,omitempty"`
	ResultCode int    `msgpack:"c,omitempty"`
	Name       string `msgpack:"n,omitempty"`
}

func (m JoinCharacterMessage) Type() string {
	return "join-character"
}

type UnjoinCharacterMessage struct {
	Result     string `msgpack:"r,omitempty"`
	ResultCode int    `msgpack:"c,omitempty"`
}

func (m UnjoinCharacterMessage) Type() string {
	return "unjoin-character"
}

type WorldsMessage struct {
	Result     string           `msgpack:"r,omitempty"`
	ResultCode int              `msgpack:"c,omitempty"`
	Worlds     []game.WorldInfo `msgpack:"w,omitempty"`
}

func (m WorldsMessage) Type() string {
	return "worlds"
}

type CreateWorldMessage struct {
	Result     string `msgpack:"r,omitempty"`
	ResultCode int    `msgpack:"c,omitempty"`
	Name       string `msgpack:"n,omitempty"`
	Password   string `msgpack:"p,omitempty"`
}

func (m CreateWorldMessage) Type() string {
	return "create-world"
}

type JoinWorldMessage struct {
	Result     string  `msgpack:"r,omitempty"`
	ResultCode int     `msgpack:"c,omitempty"`
	World      id.UUID `msgpack:"w,omitempty"`
	Password   string  `msgpack:"p,omitempty"`
}

func (m JoinWorldMessage) Type() string {
	return "join-world"
}

type LocationMessage struct {
	Result     string       `msgpack:"r,omitempty"`
	ResultCode int          `msgpack:"c,omitempty"`
	ID         id.UUID      `msgpack:"id,omitempty"`
	Objects    game.Objects `msgpack:"o,omitempty"`
	Cells      game.Cells   `msgpack:"g,omitempty"`
}

func (m LocationMessage) Type() string {
	return "location"
}

type TileMessage struct {
	Result     string             `msgpack:"r,omitempty"`
	ResultCode int                `msgpack:"c,omitempty"`
	ID         id.UUID            `msgpack:"id,omitempty"`
	Tile       game.TileArchetype `sjon:"t,omitempty"`
}

func (m TileMessage) Type() string {
	return "tile"
}

type OwnerMessage struct {
	WID       id.WID       `msgpack:"wid,omitempty"`
	Inventory game.Objects `msgpack:"i,omitempty"`
	Skills    game.Skills  `msgpack:"s,omitempty"`
}

func (m OwnerMessage) Type() string {
	return "owner"
}

type InventoryMessage struct {
	WID       id.WID        `msgpack:"wid,omitempty"`
	Inventory []game.Object `msgpack:"i,omitempty"`
}

func (m InventoryMessage) Type() string {
	return "inventory"
}

type SkillsMessage struct {
	WID    id.WID      `msgpack:"wid,omitempty"`
	Skills game.Skills `msgpack:"s,omitempty"`
}

func (m SkillsMessage) Type() string {
	return "skills"
}

type DesireMessage struct {
	Result     string             `msgpack:"r,omitempty"`
	ResultCode int                `msgpack:"c,omitempty"`
	WID        id.WID             `msgpack:"wid,omitempty"`
	Desire     game.DesireWrapper `msgpack:"d,omitempty"`
}

func (m DesireMessage) Type() string {
	return "desire"
}

type EventMessage struct {
	Event game.EventWrapper `msgpack:"e,omitempty"`
}

func (m EventMessage) Type() string {
	return "event"
}

type EventsMessage struct {
	Events []game.EventWrapper `msgpack:"e,omitempty"`
}

func (m EventsMessage) Type() string {
	return "events"
}
