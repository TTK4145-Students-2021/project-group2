package orderDistribution

import (
	"time"
)
/*
*=============================================================================
 * @Description:
 * Contains all functions used for initallizing the list of information stored
 * in the runOrder() function
 *
 * 
/*=============================================================================
*/

func initOrderList() [_orderListLength]HallOrder {
	OrderList := [_orderListLength]HallOrder{}

	initialCosts := [_numElevators]int{10000}

	for idx := 0; idx < _orderListLength; idx ++{
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

func initAllElevatorStatuses() [_numElevators]ElevatorStatus {

	CabOrders := [_numFloors]bool{false}
	OrderList := initOrderList()

	listOfElevators := [_numElevators]ElevatorStatus{}

	for i := 0; i < (_numElevators); i++ {
		Status := ElevatorStatus{
			ID:           i,
			Pos:          1,
			OrderList:    OrderList,
			IsOnline:     false, 
			DoorOpen:     false,
			CabOrders:    CabOrders,
			IsAvailable:  false, 
			IsObstructed: false,
		}
		listOfElevators[i] = Status
	}
	return listOfElevators
}


func initThisElevatorStatus(allElevatorSatuses *[_numElevators]ElevatorStatus, chans OrderChannels){
	 
	allElevatorSatuses[_thisID].IsOnline = true
	allElevatorSatuses[_thisID].IsAvailable = true

	chans.GetInitPosition <- true
	initFloor := <-chans.NewFloor
	initDoor := <-chans.DoorOpen

	allElevatorSatuses[_thisID].Pos = initFloor
	allElevatorSatuses[_thisID].DoorOpen = initDoor
}