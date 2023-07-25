package main

import "github.com/kettek/morogue/net"

type clientState int

const (
	clientStateWaiting clientState = iota
	clientStateLoggedIn
	clientStateSelectedCharacter
	clientStateInWorld
)

type client struct {
	account    Account
	character  string // Character the client is joining as.
	state      clientState
	conn       *net.Connection
	msgChan    chan net.Message
	closedChan chan error
}
