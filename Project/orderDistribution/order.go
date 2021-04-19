package orderDistribution

/*The order module has a main function called runOrders() which handels alle the inputs and ouputs from the module*/

import (
	//"os"
	"errors"
	"fmt"
	"time"
	"../config"
	"../messages"
)

type HallOrder = messages.HallOrder
type ElevatorStatus = messages.ElevatorStatus
type ButtonEvent = messages.ButtonEvent
type ButtonType = messages.ButtonType
type LampUpdate = messages.LampUpdate

const orderListLength = numFloors*2 - 2
const numElevators = config.NumElevators
const numFloors = config.NumFloors
const thisID = config.ID
const checkOfflineIntervall = config.CheckOfflineIntervall
const lampUpdateIntervall = config.LampUpdateIntervall

const UNDEFINED = messages.UNDEFINED
const BT_HallUp = messages.BT_HallUp
const BT_HallDown = messages.BT_HallDown
const BT_Cab = messages.BT_Cab


type OrderChannels struct {
	Button_press             <-chan ButtonEvent 
	Received_elevator_update <-chan ElevatorStatus    
	New_floor                <-chan int                          
	Door_status              <-chan bool                        
	Send_status              chan<- ElevatorStatus      
	Go_to_floor              chan<- int                         
	Get_init_position 	     chan<- bool
	DoorObstructed           <-chan bool
	UpdateLampMessage        chan<- LampUpdate
}

func RunOrders(chans OrderChannels) {
	fmt.Println("Order module initializing....")

	// Initialize list of all 
	allElevatorStatuses := &AllElevatorStatuses{}
	allElevatorStatuses.initAllStatuses()
	

	// Default assignedOrder and executeOrder
	assignedOrder := ButtonEvent{-1, messages.UNDEFINED}
	executeOrder := ButtonEvent{-1, messages.UNDEFINED}
	// a variable used to keep the door open when reaching a floor
	orderTimeOut := time.Now().Add(-time.Hour)

	thisElevatorStatus := &allElevatorStatuses[thisID]
	thisElevatorStatus.initStatus(chans)
	thisOrderList := &thisElevatorStatus.OrderList
	
	// TODO: Gi intuitive og beskrivende navn så vi slipper å kommentere disse!
	go thisElevatorStatus.broadcast(chans.Send_status)
	go thisElevatorStatus.sendOutLightsUpdate(chans.UpdateLampMessage)
	go checkIfElevatorsIsOffline(&allElevatorStatuses)

	for {
		select {
		case buttonEvent := <-chans.Button_press:
			fmt.Println("--------Button pressed------")
			thisElevatorStatus.updateOrders(buttonEvent)
	
			//go thisElevatorStatus.Broadcast(chans.Send_status)

		case incomingElevatorStatus := <-chans.Received_elevator_update: // new update
			// update own orderlist and otherElev with the incomming elevatorStatus     

			if elevatorStatus.ID != thisID {
				incomingOrderList := incomingElevatorStatus.OrderList
				OrderList.update(incomingOrderList)
				allElevatorStatuses.update(incomingElevatorStatus)
			}

		case floor := <-chans.New_floor:
			thisElevatorStatus.updateFloor(floor)
	
			//go sendOutStatus(chans.Send_status, thisElevatorStatus])
			fmt.Println("-> Status sendt: New floor:", floor)
			printElevatorStatus(allElevatorStatuses[thisID])

		case isOpen := <-chans.Door_status:
			thisElevatorStatuses.updateDoor(isOpen)
			
			if isOpen == true {
				allElevatorStatuses.updateOrderListCompleted(assignedOrder, &orderTimeout)
			
			}
			//go sendOutStatus(chans.Send_status, thisElevatorStatus)
			fmt.Println("-> Status sendt: Door: ", isOpen)
			printElevatorStatus(allElevatorStatuses[thisID])


		case IsObstructed := <-chans.DoorObstructed:
			allElevatorStatuses[thisID].IsObstructed = IsObstructed //TODO burde denne abstraheres ut i en egen funksjon
			//go sendOutStatus(chans.Send_status, thisElevatorStatus)

		case <-time.After(500 * time.Millisecond):
		}
		
		newAssignedOrder := assignOrder(allElevatorStatuses)
		
		thisEle
		if newAssignedOrder.Floor != -1 && allElevatorStatuses[thisID].DoorOpen && !allElevatorStatuses[thisID].IsObstructed && time.Since(orderTimeOut) > config.OrderTimeOut && allElevatorStatuses[thisID].Pos != -1 {
			assignedOrder = newAssignedOrder

			if assignedOrder != executeOrder{
				go sendingElevatorToFloor(chans.Go_to_floor, assignedOrder.Floor)
				go thisElevatorStatus.sendOutStatus(chans.Send_status)
				// fmt.Println("-> Status sendt: Go to floor ", assignedOrder.Floor)
				go thisElevatorStatus.checkIfEngineFails(assignedOrder.Floor, chans.Send_status)
				executeOrder = assignedOrder
			}
			if assignedOrder.Floor == allElevatorStatuses[thisID].Pos {
				allElevatorStatuses.updateOrderListCompleted(assignedOrder, &orderTimeOut)
			}
			//printElevatorStatus(allElevatorStatuses[thisID])
		}
	}
}


// Trengs ikke?
func sendingElevatorToFloor(channel chan<- int, goToFloor int) {
	channel <- goToFloor
}



