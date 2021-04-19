package orderDistribution

/*The order module has a main function called runOrders() which handels alle the inputs and ouputs from the module*/

import (
	"fmt"
	"time"
	"../config"
	"../messages"
)

// type HallOrder = messages.HallOrder
// type ElevatorStatus = messages.ElevatorStatus
type ButtonEvent = messages.ButtonEvent
type ButtonType = messages.ButtonType
type LampUpdate = messages.LampUpdate
// type OrderList = messages.OrderList

const _orderListLength = config.NumFloors*2 - 2
const _numElevators = config.NumElevators
const _numFloors = config.NumFloors
const _thisID = config.ID

const _pollOfflineRate = config.CheckOfflineIntervall
const _lampUpdateRate = config.LampUpdateIntervall
const _offlineTimeout = config.OfflineTimeout
const _bcastRate = config.BcastIntervall
const _checkOfflineIntervall = config.CheckOfflineIntervall

const UNDEFINED = messages.UNDEFINED
const BT_HallUp = messages.BT_HallUp
const BT_HallDown = messages.BT_HallDown
const BT_Cab = messages.BT_Cab


type OrderChannels struct {
	ButtonPress            	<-chan ButtonEvent 
	ReceiveElevatorUpdate 	<-chan ElevatorStatus    
	NewFloor                <-chan int                          
	DoorOpen              	<-chan bool                        
	SendStatus              chan<- ElevatorStatus      
	GotoFloor              	chan<- int
	RequestElevatorStatus   chan<- bool
	DoorObstructed          <-chan bool
	UpdateLampMessage       chan<- LampUpdate
}

func RunOrders(channels OrderChannels) {
	fmt.Println("Order module initializing....")

	// Initialize list of all 
	allElevatorStatuses := &AllElevatorStatuses{}
	allElevatorStatuses.init()
	

	// Default assignedOrder and executeOrder
	assignedOrder := ButtonEvent{-1, messages.UNDEFINED}
	executeOrder := ButtonEvent{-1, messages.UNDEFINED}

	thisElevatorStatus := &allElevatorStatuses[_thisID]
	thisElevatorStatus.init(channels)
	thisOrderList := &thisElevatorStatus.OrderList
	
	go thisElevatorStatus.broadcast(channels.SendStatus)
	go thisElevatorStatus.sendOutLightsUpdate(channels.UpdateLampMessage)
	go allElevatorStatuses.checkIfElevatorOffline()

	for {
		select {
		case buttonEvent := <-channels.ButtonPress:
			fmt.Println("--------Button pressed------")
			thisElevatorStatus.updateOrders(buttonEvent)
	
			//go thisElevatorStatus.Broadcast(channels.Send_status)

		case incomingElevatorStatus := <-channels.ReceiveElevatorUpdate: // new update
			// update own orderlist and otherElev with the incomming elevatorStatus     
			
			if thisElevatorStatus.ID != _thisID {
				incomingOrderList := incomingElevatorStatus.OrderList
				thisOrderList.update(incomingOrderList)
				allElevatorStatuses.update(incomingElevatorStatus)
			}

		case floor := <-channels.NewFloor:
			thisElevatorStatus.updateFloor(floor)
	
			//go sendOutStatus(channels.Send_status, thisElevatorStatus])
			fmt.Println("-> Status sendt: New floor:", floor)
			printElevatorStatus(allElevatorStatuses[_thisID])

		case isOpen := <-channels.DoorOpen:
			thisElevatorStatus.updateDoor(isOpen)
			
			if isOpen == true {
				thisElevatorStatus.updateOrderCompleted(assignedOrder)
			
			}
			//go sendOutStatus(channels.Send_status, thisElevatorStatus)
			fmt.Println("-> Status sendt: Door: ", isOpen)
			printElevatorStatus(allElevatorStatuses[_thisID])


		case IsObstructed := <-channels.DoorObstructed:
			allElevatorStatuses[_thisID].IsObstructed = IsObstructed //TODO burde denne abstraheres ut i en egen funksjon
			//go sendOutStatus(channels.Send_status, thisElevatorStatus)

		case <-time.After(500 * time.Millisecond):
		}
		
		newAssignedOrder := assignOrder(allElevatorStatuses)
		
		if thisElevatorStatus.canAcceptOrder(newAssignedOrder){
		//if newAssignedOrder.Floor != -1 && allElevatorStatuses[_thisID].DoorOpen && !allElevatorStatuses[_thisID].IsObstructed && time.Since(orderTimeOut) > config.OrderTimeOut && allElevatorStatuses[_thisID].Pos != -1 {
			assignedOrder = newAssignedOrder

			if assignedOrder != executeOrder{
				go sendingElevatorToFloor(channels.GotoFloor, assignedOrder.Floor)
				//go thisElevatorStatus.sendOutStatus(channels.Send_status)
				// fmt.Println("-> Status sendt: Go to floor ", assignedOrder.Floor)
				go thisElevatorStatus.checkEngineFailure(assignedOrder.Floor, channels.SendStatus)
				executeOrder = assignedOrder
			}
			if assignedOrder.Floor == allElevatorStatuses[_thisID].Pos {
				thisElevatorStatus.updateOrderCompleted(assignedOrder)
			}
			//printElevatorStatus(allElevatorStatuses[_thisID])
		}
	}
}


// Trengs ikke?
func sendingElevatorToFloor(channel chan<- int, goToFloor int) {
	channel <- goToFloor
}



