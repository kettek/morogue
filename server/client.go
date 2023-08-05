package server

import (
	"time"

	"github.com/kettek/morogue/game"
	"github.com/kettek/morogue/net"
)

type clientState int

const (
	clientStateWaiting clientState = iota
	clientStateLoggedIn
	clientStateSelectedCharacter
	clientStateInWorld
)

// client contains a network connection, account, character, and more.
type client struct {
	account          Account
	character        string // Character the client is joining as.
	currentCharacter *game.Character
	currentLocation  *location
	state            clientState
	conn             *net.Connection
	msgChan          chan net.Message
	closedChan       chan error
	lastWorldsSent   time.Time
}
