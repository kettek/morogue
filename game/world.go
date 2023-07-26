package game

import "github.com/kettek/morogue/id"

type World struct {
	Name       string
	ID         id.UUID // random UUIDv4 to identify the world.
	Private    bool
	Players    int
	MaxPlayers int
}
