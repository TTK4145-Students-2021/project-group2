package orders

import (
	"os"
	"time"

	"./config"
	"./elevio"
)

// assumes buttonEvent as input
// bcast elevID, targetFloor, dir

type order struct {
	hasOrder   bool
	versionNum int
	assignedElevator int
	costs[config.NumElevators] int
	timeStamp time
}

type ElevatorStatus struct {
	ID  int
	pos int
	orderList[config.NumFloors][2] order  // column 0: UP, column 1: DOWN
	dir elevio.MotorDirection
	isOnline bool
	doorOpen bool
	cabOrder[config.NumFloors] bool
}

func returnOrder() {
	// Dummy function to test elevatorController

}


func initElevatorStatus(startPos int) ElevatorStatus {
	initOrder = order{false, 0}
	initOrderList = orderList[config.NumFloors][2]

	for i, _ := range initOrderList {
		for j, _ := range initOrderList[i] {
			initOrderList[i][j] = initOrder
		}
	}
	
	return ElevatorStatus{
		ID: config.ID,
		pos: <-NewFloor,
		orderList: initOrderList,
		dir: elevio.MD_Stop,
		available: true,
		doorOpen: false,
	}
}



//updates your orderList
func updateOrderListButton(btnEvent ButtonEvent, status ElevatorStatus) {	
	if btnEvent.Button == BT_Cab {
		status.cabOrder[btnEvent.Floor] = true
	}
	else {
		newOrder = order{
			hasOrder: true,
			versionNum: status.orderList[btnEvent.Floor][btnEvent.Button] + 1,
			assignedElevator: nil,
			timeStamp: time.now(),
		}
		status.orderList[btnEvent.Floor][btnEvent.Button] = newOrder
	}
}

func updateOtherElev(status ElevatorStatus, otherElev *map[int]ElevatorStatus)  {
	otherElev[status.ID] = status
}

func updateOrderListOther(status ElevatorStatus) {

}

//Returns cost value for order
func costFunction(floor,direction){
	
	return "list of costs"
}

func pickOrder(elevatorStatus ElevatorStatus) {
	orderList[]
	floor = 0
	for order := range myElevatorStatus.orderList[0]:
		floor += 1
		if (order.hasOrder):
			costFunction(floor,"up")
		else break
	//run costfunction on all orders in orderlist
	//itterate through the list and check if you have the lowest score
	//if you have the lowest increase your score for every other
}


func runOrders(buttonPressChan chan) {

	myElevatorStatus := initElevatorStatus(floor)
	otherElevators = map[int]ElevatorStatus{}

	for {
		select{
		case buttonEvent := <- HallOrder:
			updateOrderListButton(buttonEvent):
		
		case elevatorStatus := <- ReceivedElevatorUpdate: // new update
			// update own orderlist and otherElev with the incomming elevatorStatus
			updateOtherElev(elevatorStatus, &otherElevators)
			updateOrderListOther()

		case floor := <- NewFloor:
			updateElevatorStatusFloor(floor)
	
		case isOpen := <-isDoorOpen: //recives a bool value
			updateElevatorStatusDoor(isOpen)
			if(isOpen == true)
				updateOrderListCompleted(myElevatorStatus.floor)
	}                                                                                                                                                                                                                

	pickOrder()
}




