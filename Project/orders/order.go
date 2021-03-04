package orders

import (
	"os"

	"./config"
	"./elevio"
)


func returnOrder() {
	// Dummy function to test elevatorController

}

// assumes buttonEvent as input
// bcast elevID, targetFloor, dir

type order struct {
	hasOrder   bool
	versionNum int
	assignedElevator int
}

type elevatorStatus struct {
	ID  int
	pos int
	orderList[config.NumFloors][2] order
	dir elevio.MotorDirection
	available bool
	doorOpen bool
}


otherElev[NumElevators - 1] elevatorStatus


func initElevatorStatus(startPos int) elevatorStatus{
	initOrder = order{false, 0}
	initOrderList = orderList[config.NumFloors][2]

	for i, _ := range initOrderList {
		for j, _ := range initOrderList[i] {
			initOrderList[i][j] = initOrder
		}
	}
	
	return elevatorStatus{
		ID: config.ID,
		pos: startPos,
		orderList: initOrderList
		dir: elevio.MD_Stop
		available: true
		doorOpen: false	
	}
}


//Returns cost value for every activate order
func costFunction(elevatorStatus elevatorStatus){
}

//updates your orderList
func updateOrderList(ButtonEvent){

}

func main() {
	for {
		select{
		case s := <- "order from yout own buttons"	
			updateOrderList(s) := <-"order from yout own buttons":
			
		case a := <-"updateds from other elevators":
			updateOrderList(s)
		
		case b := <-"reached a new floor":
			updateOrderList(s)
	
		case d := <-"open doors":
			if(the floor that you open the door at is in the orderList)
				updateOrderList(mission completed)
				updateOrderList(is UNAVAILABLE)	

		case e := <-"cab order completed or timed out"
			updateOrderList(is available)
		}
	}                                                                                                                                                                                                                

	cost := costFunction()
	if (you have the lowest score and is available){
		// Drive toward targetfloor
	}
}
