package main

import "github.com/kettek/morogue/net"

type clientState int

const (
	clientStateWaiting clientState = iota
	clientStateLoggedIn
	clientStateInWorld
)

type client struct {
	account    Account
	state      clientState
	conn       *net.Connection
	msgChan    chan net.Message
	closedChan chan error
}
