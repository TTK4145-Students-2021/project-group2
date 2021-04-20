package orderDistribution

import (
	"fmt"
	"time"
	"../config"
	"../messages"
)


/*=============================================================================
 * @Description:
 * Contains the toplevel run function for the order module. It mediates the flow
 * of information between modules and connects it with the general decision-making
 * logic (cost funcition, order assignment etc..). 
 * 
/*=============================================================================
*/

func RunOrders(chans OrderChannels) {
	fmt.Println("Order module initializing....")

	allElevatorStatuses := initAllElevatorStatuses()
	initThisElevatorStatus(&allElevatorStatuses, chans)

	thisElevatorStatus := &allElevatorStatuses[_thisID]

	backup := NewBackup("CabOrderBackup.txt", config.ID)
	backup.RecoverCabOrders(&thisElevatorStatus.CabOrders)

	// Default assignedOrder and executeOrder
	assignedOrder := ButtonEvent{-1, messages.UNDEFINED}
	executeOrder := ButtonEvent{-1, messages.UNDEFINED}
	
	TimeOfLastCompletion := time.Now().Add(-time.Hour) 

	go broadcast(chans.SendStatus, thisElevatorStatus)
	go sendOutLightsUpdate(chans.UpdateLampMessage, thisElevatorStatus)
	go checkIfElevatorOffline(&allElevatorStatuses)

	for {
		select {
		case buttonEvent := <-chans.ButtonPress:
			updateOrders(buttonEvent, thisElevatorStatus)   
			go backup.SaveCabOrders(thisElevatorStatus.CabOrders)
			 
		case incomingElevatorStatus := <-chans.ReceivedElevatorUpdate:     
			if incomingElevatorStatus.ID != _thisID {
				updateElevatorStatuses(incomingElevatorStatus, &allElevatorStatuses)
				updateOrderList(incomingElevatorStatus, &allElevatorStatuses)
			}

		case floor := <-chans.NewFloor:
			updateElevatorStatusFloor(floor, &allElevatorStatuses)
			go backup.SaveCabOrders(thisElevatorStatus.CabOrders)

		case isOpen := <-chans.DoorOpen:
			updateElevatorStatusDoor(isOpen, &allElevatorStatuses)
			if isOpen == true {
				updateOrderListCompleted(&allElevatorStatuses, assignedOrder, &TimeOfLastCompletion)
			}
			go backup.SaveCabOrders(thisElevatorStatus.CabOrders)

		case IsObstructed := <-chans.DoorObstructed:
			allElevatorStatuses[_thisID].IsObstructed = IsObstructed 

		case <-time.After(_selectBreakOutRate):
		}

		newAssignedOrder := assignOrder(&allElevatorStatuses)
		if canAcceptOrder(newAssignedOrder, *thisElevatorStatus, TimeOfLastCompletion){
		
			assignedOrder = newAssignedOrder
			if assignedOrder != executeOrder{
				go sendingElevatorToFloor(chans.GotoFloor, assignedOrder.Floor)
				go checkIfEngineFails(assignedOrder.Floor, &allElevatorStatuses[_thisID], chans.SendStatus)
				executeOrder = assignedOrder
			}
			if assignedOrder.Floor == allElevatorStatuses[_thisID].Pos {
				updateOrderListCompleted(&allElevatorStatuses, assignedOrder, &TimeOfLastCompletion)
				go backup.SaveCabOrders(thisElevatorStatus.CabOrders)
			}

		}
	}
}

