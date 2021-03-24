package orders

import (
	"os"
	"time"

	"../config"
	"../elevio"
)

// assumes buttonEvent as input
// bcast elevID, targetFloor, dir

type order struct {
	hasOrder bool
	floor int
	direction int	//up = 0, down = 1
	versionNum int
	costs map[int]int
	timeStamp time
}

type ElevatorStatus struct {
	ID  int
	pos int
	orderList[config.NumFloors*2-2] order  // 
	dir elevio.MotorDirection
	isOnline bool
	doorOpen bool
	cabOrder[config.NumFloors] bool
	isAvailable bool
}

func returnOrder() {
	// Dummy function to test elevatorController
}


func initElevatorStatus(startPos int) ElevatorStatus {
	var emptyIntList[config.NumFloors] int
	var emptyBoolList[config.NumFloors] bool
	
	initOrder := order{false, 0, emptyList, time.now()}
	initOrderList := orderList[config.NumFloors][2]

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
		isOnline: true,
		doorOpen: false,
		cabOrder: emptyBoolList
		isAvailable: true
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





	assignedFloor = 0     //Defalut assigned floor
	// Itterate through up orders and calcuate all the costs

	for order := range allElevators[config.ID].orderList:    
		if (order.hasOrder):
			order.costs = costFunction(order)  //*****Should you take a copy here????
		else 
			var zeros[]              //TODO: Needs to be finished. Idea is just to inseart a slice of zeros. 
			order.costs = [config.NumElevators]
	
	//TODO:itterate through down orders and calcuate all the costs		
}


func runOrders(buttonPressChan chan) {

	myElevatorStatus := initElevatorStatus(floor)
	allElevators = map[int]ElevatorStatus{}

	for {
		select{
		case buttonEvent := <- HallOrder:
			updateOrderListButton(buttonEvent):
		
		case elevatorStatus := <- ReceivedElevatorUpdate: // new update
			// update own orderlist and otherElev with the incomming elevatorStatus
			updateOtherElev(elevatorStatus, &allElevators)
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




