package game

import (
	"github.com/kettek/morogue/game"
)

// Mover converts bind actions into move directions.
type Mover struct {
	targetX, targetY int
	currX, currY     int
}

// Init does nothing at the moment.
func (mover *Mover) Init() {
}

// Update returns the current move direction from the current actions.
func (mover *Mover) Update(binds Binds) game.MoveDirection {
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

	return d
}
