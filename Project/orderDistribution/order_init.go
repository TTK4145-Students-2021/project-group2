package orderDistribution

import (
	"time"
)


func initAllElevatorStatuses() [numElevators]ElevatorStatus {

	CabOrders := [numFloors]bool{false}
	OrderList := initOrderList()

	listOfElevators := [numElevators]ElevatorStatus{}

	for i := 0; i < (numElevators); i++ {
		Status := ElevatorStatus{
			ID:           i,
			Pos:          1,
			OrderList:    OrderList,
			IsOnline:     false, //NB isOline is deafult false
			DoorOpen:     false,
			CabOrders:    CabOrders,
			IsAvailable:  false, //NB isAvaliable is deafult false
			IsObstructed: false,
		}
		listOfElevators[i] = Status
	}
	return listOfElevators
}

func initThisElevatorStatus(allElevatorSatuses *[numElevators]ElevatorStatus, chans OrderChannels){
	
	allElevatorSatuses[curID].IsOnline = true
	allElevatorSatuses[curID].IsAvailable = true

	chans.Get_init_position <- true
	initFloor := <-chans.New_floor
	initDoor := <-chans.Door_status

	allElevatorSatuses[curID].Pos = initFloor
	allElevatorSatuses[curID].DoorOpen = initDoor
}