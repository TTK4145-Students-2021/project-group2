package orders

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
type ButtonEvent_message = messages.ButtonEvent_message
type ButtonType_msg = messages.ButtonType_msg
type LampUpdate_message = messages.LampUpdate_message
const orderListLength = numFloors*2 - 2
const numElevators = config.NumElevators
const numFloors = config.NumFloors
const curID = config.ID

const UNDEFINED = messages.UNDEFINED
const BT_HallUp = messages.BT_HallUp
const BT_HallDown = messages.BT_HallDown
const BT_Cab = messages.BT_Cab


type OrderChannels struct {
	Button_press             <-chan ButtonEvent_message 
	Received_elevator_update <-chan ElevatorStatus    
	New_floor                <-chan int                          
	Door_status              <-chan bool                        
	Send_status              chan<- ElevatorStatus      
	Go_to_floor              chan<- int                         
	Get_init_position 	     chan<- bool
	DoorObstructed           <-chan bool
	UpdateLampMessage        chan<- LampUpdate_message
}

// updates OrderList and CabOrders based on button presses
func updateOrderListButton(btnEvent messages.ButtonEvent_message, thisElevator *ElevatorStatus) {

	if btnEvent.Button == messages.BT_Cab {
		thisElevator.CabOrders[btnEvent.Floor] = true

	} else { //TODO remove placeInList or look at OrderList's structure
		idx := floorToOrderListIdx(btnEvent.Floor, btnEvent.Button)
		if thisElevator.OrderList[idx].HasOrder == false {
			thisElevator.OrderList[idx].HasOrder = true
			thisElevator.OrderList[idx].VersionNum += 1
		}
	}
}

//Updates the information list with the incoming update
func updateOtherElev(incomingStatus ElevatorStatus, list *[numElevators]ElevatorStatus) {
	list[incomingStatus.ID] = incomingStatus //accountign for zero indexing
}
func updateElevatorStatusDoor(value bool, list *[numElevators]ElevatorStatus) {
	list[curID].DoorOpen = value
}
func updateElevatorStatusFloor(pos int, list *[numElevators]ElevatorStatus) {
	list[curID].Pos = pos
}

func updateOrderListOther(incomingStatus ElevatorStatus, list *[numElevators]ElevatorStatus) {

	for i := 0; i < len(list[curID].OrderList); i++ {

		// the first if sentences checks for faiure is the version controll
		if incomingStatus.OrderList[i].VersionNum == list[curID].OrderList[i].VersionNum && incomingStatus.OrderList[i].HasOrder != list[curID].OrderList[i].HasOrder {
			fmt.Println("There is smoething worng with the version controll of the order system")
			list[curID].OrderList[i].HasOrder = true
			list[curID].OrderList[i].VersionNum += 1
		}
		// this if statement checks if the order should be updated and updates it
		if incomingStatus.OrderList[i].VersionNum > list[curID].OrderList[i].VersionNum {
			list[curID].OrderList[i].HasOrder = incomingStatus.OrderList[i].HasOrder
			list[curID].OrderList[i].VersionNum = incomingStatus.OrderList[i].VersionNum
		}
	}
}

func updateOrderListCompleted(list *[numElevators]ElevatorStatus, assignedOrder ButtonEvent_message, orderTimeOut *time.Time) {
	curFloor := list[curID].Pos

	// checks if the elevator has a cab call for the current floor
	if list[curID].CabOrders[curFloor] == true && list[curID].DoorOpen == true {
		list[curID].CabOrders[curFloor] = false
		*orderTimeOut = time.Now()
		return
	}

	orderIdx := floorToOrderListIdx(assignedOrder.Floor, assignedOrder.Button)

	if list[curID].OrderList[orderIdx].HasOrder && list[curID].OrderList[orderIdx].Floor == curFloor {
		list[curID].OrderList[orderIdx].HasOrder = false
		list[curID].OrderList[orderIdx].VersionNum += 1
		*orderTimeOut = time.Now()
		fmt.Println("Hall order compleete:", orderListIdxToFloor(orderIdx))
	}
}

func sendOutLightsUpdate(sendLampUpdate chan<- messages.LampUpdate_message, status *ElevatorStatus) {
	for {
		for _, order := range status.OrderList {
			light := messages.LampUpdate_message{
				Floor:  order.Floor,
				Button: messages.ButtonType_msg(order.Direction), // = 0 up og 1 down
				Turn:   order.HasOrder,
			}
			sendLampUpdate <- light
			//send light at channel
		}
		for i, order := range status.CabOrders {

			light := messages.LampUpdate_message{
				Floor:  i,
				Button: messages.BT_Cab, // = 0 up og 1 down
				Turn:   order,
			}
			sendLampUpdate <- light //send light at channel
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func checkIfElevatorOffline(allElevatorSatuses *[numElevators]ElevatorStatus) {
	for{
		for elevID := range allElevatorSatuses {
			if time.Since(allElevatorSatuses[elevID].Timestamp) > config.OfflineTimeout && elevID != curID {
				if allElevatorSatuses[elevID].IsOnline == true {
					fmt.Printf("[-] Elevator offline ID: %02d\n", elevID)
				}
				allElevatorSatuses[elevID].IsOnline = false
			} else {
				if allElevatorSatuses[elevID].IsOnline == false {
					fmt.Printf("[+] Elevator online ID: %02d\n", elevID)
				}
				allElevatorSatuses[elevID].IsOnline = true
			}
		}
	}
}

func sendOutStatus(channel chan<- ElevatorStatus, status ElevatorStatus) {
	status.Timestamp = time.Now()
	channel <- status
}
func sendingElevatorToFloor(channel chan<- int, goToFloor int) {
	channel <- goToFloor
}

func contSend(channel chan<- ElevatorStatus, status *ElevatorStatus) {
	for {
		printElevatorStatus(*status)
		go sendOutStatus(channel, *status)
		time.Sleep(config.BcastIntervall)
	}
}
func checkIfEngineFails(assignedFloor int, status *ElevatorStatus, send_status chan<- ElevatorStatus) {
	faultCounter := 0
	curPos := -1
	lastPos := -1
	for assignedFloor != -1 {
		curPos = status.Pos
		fmt.Println(curPos)
		if curPos == assignedFloor {
			return
		}
		if curPos != lastPos {
			lastPos = curPos
		} else {
			faultCounter += 1
			if faultCounter > 10 {
				status.IsAvailable = false
				printElevatorStatus(*status)
				go sendOutStatus(send_status, *status)
				err := errors.New("engine failure")
				panic(err)
				return
			}
		}
		time.Sleep(time.Second * 1)
	}
}

func RunOrders(chans OrderChannels) {
	fmt.Println("Order module initializing....")

	allElevatorSatuses := initAllElevatorStatuses()
	initThisElevatorStatus(&allElevatorSatuses, chans)

	// Default assignedOrder and executeOrder
	assignedOrder := ButtonEvent_message{-1, messages.UNDEFINED}
	executeOrder := ButtonEvent_message{-1, messages.UNDEFINED}
	// 
	orderTimeOut := time.Now().Add(-time.Hour)

	thisElevatorStatus := &allElevatorSatuses[curID]
	//thread continuinling sending out this elevatorStatus to other elevators
	go contSend(chans.Send_status, thisElevatorStatus)
	//thread continuinling sending out this LampMessages to elevator module
	go sendOutLightsUpdate(chans.UpdateLampMessage, thisElevatorStatus)
	// Thread to check if elevators have gone offline or come back online
	go checkIfElevatorOffline(&allElevatorSatuses)

	for {
		select {
		case buttonEvent := <-chans.Button_press:
			fmt.Println("--------Button pressed------")
			updateOrderListButton(buttonEvent, &allElevatorSatuses[curID])
			
			fmt.Println("-> Status sendt: Buttonpress,", buttonEvent.Button)

			go sendOutStatus(chans.Send_status, allElevatorSatuses[curID])

		case elevatorStatus := <-chans.Received_elevator_update: // new update
			// update own orderlist and otherElev with the incomming elevatorStatus     

			if elevatorStatus.ID != curID {
				updateOtherElev(elevatorStatus, &allElevatorSatuses)
				updateOrderListOther(elevatorStatus, &allElevatorSatuses)
				fmt.Println("<- Recived status")
			}

		case floor := <-chans.New_floor:
			updateElevatorStatusFloor(floor, &allElevatorSatuses)
			updateOrderListCompleted(&allElevatorSatuses, assignedOrder, &orderTimeOut)
			go sendOutStatus(chans.Send_status, allElevatorSatuses[curID])
			fmt.Println("-> Status sendt: New floor:", floor)

		case isOpen := <-chans.Door_status:
			updateElevatorStatusDoor(isOpen, &allElevatorSatuses)
			if isOpen == true {
				updateOrderListCompleted(&allElevatorSatuses, assignedOrder, &orderTimeOut)
			}
			go sendOutStatus(chans.Send_status, allElevatorSatuses[curID])
			fmt.Println("-> Status sendt: Door: ", isOpen)

		case IsObstructed := <-chans.DoorObstructed:
			allElevatorSatuses[curID].IsObstructed = IsObstructed //TODO burde denne abstraheres ut i en egen funksjon
			go sendOutStatus(chans.Send_status, allElevatorSatuses[curID])

		case <-time.After(500 * time.Millisecond):
		}
		newAssignedOrder := assignOrder(&allElevatorSatuses)
	
		if newAssignedOrder.Floor != -1 && allElevatorSatuses[curID].DoorOpen && !allElevatorSatuses[curID].IsObstructed && time.Since(orderTimeOut) > config.OrderTimeOut {
			assignedOrder = newAssignedOrder

			if assignedOrder != executeOrder{
				go sendingElevatorToFloor(chans.Go_to_floor, assignedOrder.Floor)
				go sendOutStatus(chans.Send_status, allElevatorSatuses[curID])
				fmt.Println("-> Status sendt: Go to floor ", assignedOrder.Floor)
				go checkIfEngineFails(assignedOrder.Floor, &allElevatorSatuses[curID], chans.Send_status)
				executeOrder = assignedOrder
			}
			if assignedOrder.Floor == allElevatorSatuses[curID].Pos {
				updateOrderListCompleted(&allElevatorSatuses, assignedOrder, &orderTimeOut)
			}
			//printElevatorStatus(allElevatorSatuses[curID])
		}
	}
}

