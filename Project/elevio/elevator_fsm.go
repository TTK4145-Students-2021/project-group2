package elevio

type ElevatorState int

const (
	AT_FLOOR   ElevatorState = 0
	GOING_UP                 = 1
	GOING_DOWN               = 2
)

func takeOrder(button ButtonType, floor int, value bool) {

	switch button {
	case button == BT_HallUp:
		//
	case button.BT_HallDown:
		//What to do?
	case button.BT_Cab:
		//What to do?
	}
}
