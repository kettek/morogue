package game

import (
	"github.com/kettek/morogue/net"
)

type Mover struct {
	targetX, targetY int
	currX, currY     int
}

func (mover *Mover) Init() {
}

func (mover *Mover) Update(binds Binds) net.MoveDirection {
	var d net.MoveDirection
	if binds.IsActionHeld("move-upleft") == 0 {
		d = net.UpLeftMoveDirection
	} else if binds.IsActionHeld("move-upright") == 0 {
		d = net.UpRightMoveDirection
	} else if binds.IsActionHeld("move-downleft") == 0 {
		d = net.DownLeftMoveDirection
	} else if binds.IsActionHeld("move-downright") == 0 {
		d = net.DownRightMoveDirection
	} else if binds.IsActionHeld("move-left") == 0 {
		d = net.LeftMoveDirection
	} else if binds.IsActionHeld("move-right") == 0 {
		d = net.RightMoveDirection
	} else if binds.IsActionHeld("move-down") == 0 {
		d = net.DownMoveDirection
	} else if binds.IsActionHeld("move-up") == 0 {
		d = net.UpMoveDirection
	}

	return d
}
