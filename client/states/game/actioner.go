package game

import (
	"github.com/kettek/morogue/game"
)

// Actioner converts bind actions into move directions.
type Actioner struct {
	targetX, targetY int
	currX, currY     int
}

// Init does nothing at the moment.
func (actioner *Actioner) Init() {
}

// Update returns the current move direction from the current actions.
func (actioner *Actioner) Update(binds Binds) game.Desire {
	var d game.MoveDirection
	if binds.IsActionHeld("move-upleft") == 0 {
		d = game.UpLeftMoveDirection
	} else if binds.IsActionHeld("move-upright") == 0 {
		d = game.UpRightMoveDirection
	} else if binds.IsActionHeld("move-downleft") == 0 {
		d = game.DownLeftMoveDirection
	} else if binds.IsActionHeld("move-downright") == 0 {
		d = game.DownRightMoveDirection
	} else if binds.IsActionHeld("move-left") == 0 {
		d = game.LeftMoveDirection
	} else if binds.IsActionHeld("move-right") == 0 {
		d = game.RightMoveDirection
	} else if binds.IsActionHeld("move-down") == 0 {
		d = game.DownMoveDirection
	} else if binds.IsActionHeld("move-up") == 0 {
		d = game.UpMoveDirection
	}
	if d != 0 {
		return game.DesireMove{
			Direction: d,
		}
	}
	if binds.IsActionHeld("bash") == 0 {
		return game.DesireBash{}
	}

	return nil
}
