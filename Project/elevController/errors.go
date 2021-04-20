package elevController

import (
	"errors"
)

// ErrEventRejected returns if the elevator cannot process the event in the state that it's in
var ErrEventRejected = errors.New("event rejected")

// ErrElevUnavailable returns if the elevator can't accept an order in it's current state
var ErrElevUnavailable = errors.New("elevator unavailable")

// ErrIllegalMoveUp returns if the elevator can't move up
var ErrIllegalMoveUp = errors.New("illegal move up, already at top")

// ErrIllegalMoveDown returns if the elevator can't move down
var ErrIllegalMoveDown = errors.New("illegal move down, already at bottom")

// ErrNoStateAction return if there are errors in the State: Action configurations
var ErrNoStateAction = errors.New("there are no actions in the next state")

// ErrExecOrderFailed returns if an order could not be executed
var ErrExecOrderFailed = errors.New("order execution failed")
