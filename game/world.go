package game

import "github.com/kettek/morogue/id"

// WorldInfo is the information of a world and is used to send worlds to
// clients and for clients to use to join a world.
type WorldInfo struct {
	Name       string
	ID         id.UUID // random UUIDv4 to identify the world.
	Private    bool
	Players    int
	MaxPlayers int
}
