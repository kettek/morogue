package game

// MoveDirection is a direction that an object can move.
type MoveDirection uint8

// Position returns the x and y direction of the MoveDirection.
func (d MoveDirection) Position() (int, int) {
	switch d {
	case LeftMoveDirection:
		return -1, 0
	case RightMoveDirection:
		return 1, 0
	case UpMoveDirection:
		return 0, -1
	case DownMoveDirection:
		return 0, 1
	case UpLeftMoveDirection:
		return -1, -1
	case UpRightMoveDirection:
		return 1, -1
	case DownLeftMoveDirection:
		return -1, 1
	case DownRightMoveDirection:
		return 1, 1
	}
	return 0, 0
}

// Our movement directions.
const (
	UpMoveDirection        MoveDirection = 8
	LeftMoveDirection      MoveDirection = 4
	RightMoveDirection     MoveDirection = 6
	DownMoveDirection      MoveDirection = 2
	UpLeftMoveDirection    MoveDirection = 7
	UpRightMoveDirection   MoveDirection = 9
	DownRightMoveDirection MoveDirection = 3
	DownLeftMoveDirection  MoveDirection = 1
	CenterMoveDirection    MoveDirection = 5
)
