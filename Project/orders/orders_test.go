package orders

import (
	"fmt"
	"testing"
	"../elevio"
	"fmt"
	//"../config"
)

func TestOrders(t *testing.T) {

	fmt.Println("Kjørrr test:------------------")
	allElevators := initAllElevatorStatuses()
	btnEvent1 := elevio.ButtonEvent{
		Floor : 3,
		Button : 0, // 0 = up, 1 = down , 2 = hall
	}
	btnEvent2 := elevio.ButtonEvent{
		Floor : 1,
		Button : 2, // 0 = up, 1 = down , 2 = hall
	}
	btnEvent3 := elevio.ButtonEvent{
		Floor : 4,
		Button : 1, // 0 = up, 1 = down , 2 = hall
	}
	dummyStatus := allElevators[2]
	updateOrderListButton(btnEvent1, &dummyStatus)
	updateOrderListButton(btnEvent2, &dummyStatus)
	updateOrderListButton(btnEvent3, &dummyStatus)

	updateOtherElev(dummyStatus, &allElevators)
	printElevatorStatus(allElevators[2])

	// creating a btnEvent
}





