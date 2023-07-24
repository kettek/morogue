package net

import (
	"encoding/json"
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
