package elevator

type MotorDirection int

const (
	MD_Up   MotorDirection = 1
	MD_Down                = -1
	MD_Stop                = 0
)

type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
)

type ElevatorState int

const (
	UNAVAILABLE ElevatorState = -1 //Is this really needed? We'll see
	IDLE                      = 0
	GOING_UP                  = 1
	GOING_DOWN                = 2
)

const NoFloor int = -1
