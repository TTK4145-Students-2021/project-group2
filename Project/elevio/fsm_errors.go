package elevio

import (
	"errors"
)

// ErrEventRejected returns if the elevator cannot process the event in the state that it's in
var ErrEventRejected = errors.New("event rejected")

// ErrIllegalMoveUp returns if the elevator can't move up
var ErrIllegalMoveUp = errors.New("illegal move up, already at top")

// ErrIllegalMoveDown returns if the elevator can't move down
var ErrIllegalMoveDown = errors.New("illegal move down, already at bottom")

// ErrNoStateAction return if there are errors in the State: Action configurations
var ErrNoStateAction = errors.New("there are no actions in the next state")
