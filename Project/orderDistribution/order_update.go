package orderDistribution

import(
	"../messages"
	"time"
	"fmt"
)

/*
*=============================================================================
 * @Description:
 * Contains functions used for updating the information stored and used in runOrders()
 * 
 *
 * 
/*=============================================================================
*/


func updateOrders(btnEvent messages.ButtonEvent, thisElevator *ElevatorStatus) {

	if btnEvent.Button == messages.BT_Cab {
		thisElevator.CabOrders[btnEvent.Floor] = true

	} else {
		idx := floorToOrderListIdx(btnEvent.Floor, btnEvent.Button)
		if thisElevator.OrderList[idx].HasOrder == false {
			thisElevator.OrderList[idx].HasOrder = true
			thisElevator.OrderList[idx].VersionNum += 1
		}
	}
}

//Updates the information list with the incoming update
func updateElevatorStatuses(incomingStatus ElevatorStatus, list *[_numElevators]ElevatorStatus) {
	list[incomingStatus.ID] = incomingStatus
}
func updateElevatorStatusDoor(value bool, list *[_numElevators]ElevatorStatus) {
	list[_thisID].DoorOpen = value
}
func updateElevatorStatusFloor(pos int, list *[_numElevators]ElevatorStatus) {
	list[_thisID].Pos = pos
}

func updateOrderList(incomingStatus ElevatorStatus, allElevatorSatuses *[_numElevators]ElevatorStatus) {
	for i := 0; i < len(allElevatorSatuses[_thisID].OrderList); i++ {

		// the first if sentences checks for faiure is the version controll
		if incomingStatus.OrderList[i].VersionNum == allElevatorSatuses[_thisID].OrderList[i].VersionNum && incomingStatus.OrderList[i].HasOrder != allElevatorSatuses[_thisID].OrderList[i].HasOrder {
			fmt.Println("There is smoething worng with the version controll of the order system")
			allElevatorSatuses[_thisID].OrderList[i].HasOrder = true
			allElevatorSatuses[_thisID].OrderList[i].VersionNum += 1
		}
		// this if statement checks if the order should be updated and updates it
		if incomingStatus.OrderList[i].VersionNum > allElevatorSatuses[_thisID].OrderList[i].VersionNum {
			allElevatorSatuses[_thisID].OrderList[i].HasOrder = incomingStatus.OrderList[i].HasOrder
			allElevatorSatuses[_thisID].OrderList[i].VersionNum = incomingStatus.OrderList[i].VersionNum
		}
	}
}

func updateOrderListCompleted(list *[_numElevators]ElevatorStatus, assignedOrder ButtonEvent, TimeOfLastCompletion *time.Time) {
	curFloor := list[_thisID].Pos

	// checks if the elevator has a cab call for the current floor
	if list[_thisID].CabOrders[curFloor] == true && list[_thisID].DoorOpen == true && assignedOrder.Floor==curFloor {
		list[_thisID].CabOrders[curFloor] = false
		*TimeOfLastCompletion = time.Now()
		return
	}

	orderIdx := floorToOrderListIdx(assignedOrder.Floor, assignedOrder.Button)
	if assignedOrder.Button == -1{
		return
	}

	if list[_thisID].OrderList[orderIdx].HasOrder && list[_thisID].OrderList[orderIdx].Floor == curFloor {
		list[_thisID].OrderList[orderIdx].HasOrder = false
		list[_thisID].OrderList[orderIdx].VersionNum += 1
		*TimeOfLastCompletion = time.Now()
	}
}