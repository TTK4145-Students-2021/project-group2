package orders

import (
	"time"
)
func initOrderList() [orderListLength]HallOrder {
	OrderList := [orderListLength]HallOrder{}

	// initalizing HallUp Orders
	initialCosts := [numElevators]int{10000}

	for idx := 0; idx < orderListLength; idx ++{
		OrderList[idx] = HallOrder{
			HasOrder:   false,
			Floor:      orderListIdxToFloor(idx),
			Direction:  orderListIdxToBtnTyp(idx),
			VersionNum: 0,
			Costs:      initialCosts,
			TimeStamp:  time.Now(),
		}
	}
	return OrderList
}

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
	
	allElevatorSatuses[thisID].IsOnline = true
	allElevatorSatuses[thisID].IsAvailable = true

	chans.Get_init_position <- true
	initFloor := <-chans.New_floor
	initDoor := <-chans.Door_status

	allElevatorSatuses[thisID].Pos = initFloor
	allElevatorSatuses[thisID].DoorOpen = initDoor
}