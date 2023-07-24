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
	case "ping":
		var m PingMessage
		json.Unmarshal(w.Data, &m)
		return &m
	}
	return nil
}

type Message interface {
	Type() string
}

type PingMessage struct {
}

func (m *PingMessage) Type() string {
	return "ping"
}
