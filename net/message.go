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

func (m ArchetypesMessage) Type() string {
	return "archetypes"
}

type CharactersMessage struct {
	Characters []game.Character `json:"c,omitempty"`
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
