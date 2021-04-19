package orders

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
	Button_press             <-chan ButtonEvent_message //Elevator communiaction
	Received_elevator_update <-chan ElevatorStatus      //Network communication
	New_floor                <-chan int                          //Elevator communiaction
	Door_status              <-chan bool                         //Elevator communiaction
	Send_status              chan<- ElevatorStatus      //Network communication
	Go_to_floor              chan<- int                          //Elevator communiaction
	AskElevatorForUpdate     chan<- bool
	DoorObstructed           <-chan bool
	UpdateLampMessage        chan<- LampUpdate_message
}

func initOrderList() [orderListLength]HallOrder {
	OrderList := [orderListLength]HallOrder{}

	// initalizing HallUp
	initialCosts := [numElevators]int{}
	for i := 0; i < numElevators; i++ {
		initialCosts[i] = 10000
	}
	for i := 0; i < (numFloors - 1); i++ {
		OrderList[i] = HallOrder{
			HasOrder:   false,
			Floor:      i,
			// Direction:  0, //up = 0, down = 1
			// Direction:  messages.BT_HallUp, //up = 0, down = 1
			Direction:  BT_HallUp, //up = 0, down = 1
			VersionNum: 0,
			Costs:      initialCosts,
			TimeStamp:  time.Now(),
		}
	}

	// initalizing HallDown
	for i := numFloors - 1; i < orderListLength; i++ {
		OrderList[i] = HallOrder{
			HasOrder:   false,
			Floor:      i + 2 - numFloors,
			// Direction:  1, //up = 0, down = 1
			// Direction:  messages.BT_HallDown, //up = 0, down = 1
			Direction:  BT_HallDown, //up = 0, down = 1
			VersionNum: 0,
			Costs:      initialCosts,
			TimeStamp:  time.Now(),
		}
	}
	return OrderList
}

func initAllElevatorStatuses() [numElevators]ElevatorStatus {

	CabOrders := [numFloors]bool{}
	OrderList := initOrderList()

	listOfElevators := [numElevators]ElevatorStatus{}

	for i := 0; i < (numElevators); i++ {
		Status := ElevatorStatus{
			ID:           i,
			Pos:          1,
			OrderList:    OrderList,
			Dir:          2,
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

// updates OrderList and CabOrders based on button presses
func updateOrderListButton(btnEvent ButtonEvent_message, thisElevator *ElevatorStatus) {

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

func costFunction(num int, list *[numElevators]ElevatorStatus) {
	curOrder := list[curID].OrderList[num]
	for i := 0; i < numElevators; i++ { //i is already 0 indexed

		elevator := list[i]
		cost := 0
		cost += i
		if !elevator.IsAvailable || !elevator.IsOnline {
			cost += 4000
		}
		if elevator.IsObstructed {
			cost += 100
		}

		// for loop checking if the elevator has a CabOrders
		for j := 0; j < numFloors; j++ { //organize into function
			if elevator.CabOrders[j] == true {
				cost += 3000
				break
			}
		}
		//cost += waitingTimePenalty(curOrder)
		cost += distanceFromOrder(list[i].Pos, curOrder.Floor)

		list[curID].OrderList[num].Costs[i] = cost
	}
}

func waitingTimePenalty(curOrder HallOrder) int {
	extraCost := 0
	timeWaited := time.Now().Sub(curOrder.TimeStamp) //organize into function
	//larger than 1 min
	if timeWaited > time.Duration(time.Minute*1) {
		extraCost += -15
	}
	//larger than 3 min
	if timeWaited > time.Duration(time.Minute*3) {
		extraCost += -45
	}
	return extraCost
}

func distanceFromOrder(elevatorFloor int, goingToFloor int) int {
	floorsAway := elevatorFloor - goingToFloor
	if floorsAway < 0 { // takeing abs of floorsAway
		floorsAway = -floorsAway
	}
	return floorsAway * 10
}

func goToFloor(list *[numElevators]ElevatorStatus) ButtonEvent_message {
	// runs through the costFunction for all true orders
	for i := 0; i < len(list[curID].OrderList); i++ {
		if list[curID].OrderList[i].HasOrder { // TODO: rename list
			costFunction(i, list)
		} else {
			for j := 0; j < numElevators; j++ {
				list[curID].OrderList[i].Costs[j] = 10000
			}
		}
	}
	pickedOrder := ButtonEvent_message{-1, -1} // buttonType -1 = UNDEFINED

	//	checking if the elevator has any cabOrders and if assigning that ass the goToFloor
	shortestCabDistance := numFloors * 10 // random variable larger than numFloors
	for floor := 0; floor < numFloors; floor++ {
		if list[curID].CabOrders[floor] == true {
			curCabDistance := distanceFromOrder(list[curID].Pos, floor)
			if curCabDistance < shortestCabDistance {
				pickedOrder.Floor = floor
				pickedOrder.Button = messages.BT_Cab
			}
		}
	}
	// if there is an cab order return that as going to floor
	if pickedOrder.Button != messages.UNDEFINED {
		return pickedOrder
	}
	// Pick the order where you have the lowest cost

	costMatrix := make([][numElevators]int, numFloors*2-2)
	for i := range costMatrix {
		costMatrix[i] = list[curID].OrderList[i].Costs
	}
	lowestCost := 500

	for i := 0; i < numElevators; i++ {

		orderIdx, ID := minumumRowCol(costMatrix)
		cost := costMatrix[orderIdx][ID]
		if ID == curID && cost < lowestCost && list[curID].OrderList[orderIdx].HasOrder {
			pickedOrder.Floor = orderListIdxToFloor(orderIdx)
			pickedOrder.Button = ButtonType_msg(list[curID].OrderList[orderIdx].Direction)
			return pickedOrder
		} else {
			for i := range list[curID].OrderList {
				costMatrix[i][ID] = 10000
			}
			for i := 0; i < numElevators; i++ {
				costMatrix[orderIdx][i] = 10000
			}
		}
	}
	return pickedOrder //return the floor that the elevator should go
}

func orderListIdxToFloor(idx int) int {
	if idx < numFloors-1 {
		return idx
	} else {
		return idx - numFloors + 2
	}
}

func floorToOrderListIdx(floor int, dir ButtonType_msg) int {
	if dir == messages.BT_HallUp {
		// fmt.Println("DOWN")
		return floor
	} else {
		// fmt.Println("UP:", floor + numFloors - 2)
		return floor + numFloors - 2
	}
}

func minumumRowCol(list [][numElevators]int) (int, int) {
	min := config.MaxInt
	minRowIndex := -1
	minColIndex := -1
	for rowIndex, row := range list {
		for colIndex, number := range row {
			if number < min {
				min = number
				minRowIndex = rowIndex
				minColIndex = colIndex
			}
		}
	}
	return minRowIndex, minColIndex
}

func ifXinSliceInt(x int, slice []int) bool {
	for _, j := range slice {
		if x == j {
			return true
		}
	}
	return false
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

func sendOutLightsUpdate(sendLampUpdate chan<- LampUpdate_message, status *ElevatorStatus) {
	for {
		for _, order := range status.OrderList {
			light := LampUpdate_message{
				Floor:  order.Floor,
				Button: ButtonType_msg(order.Direction), // = 0 up og 1 down
				Turn:   order.HasOrder,
			}
			sendLampUpdate <- light
			//send light at channel
		}
		for i, order := range status.CabOrders {

			light := LampUpdate_message{
				Floor:  i,
				Button: messages.BT_Cab, // = 0 up og 1 down
				Turn:   order,
			}
			sendLampUpdate <- light //send light at channel
		}
		time.Sleep(time.Millisecond * 100)
	}
}
// -------------- Checked so far
func checkIfElevatorOffline(allElevators *[numElevators]ElevatorStatus) {
	for elevID := range allElevators {
		if time.Since(allElevators[elevID].Timestamp) > config.OfflineTimeout && elevID != curID {
			if allElevators[elevID].IsOnline == true {
				fmt.Printf("[-] Elevator offline ID: %02d\n", elevID)
			}
			allElevators[elevID].IsOnline = false
		} else {
			if allElevators[elevID].IsOnline == false {
				fmt.Printf("[+] Elevator online ID: %02d\n", elevID)
			}
			allElevators[elevID].IsOnline = true
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

func contSend(channel chan<- ElevatorStatus, list *[numElevators]ElevatorStatus) {
	for {
		status := *list
		go sendOutStatus(channel, status[curID])
		time.Sleep(config.BcastIntervall)
	}
}
func checkIfEngineFails(assignedFloor int, status *ElevatorStatus, send_status chan<- ElevatorStatus) {
	faultCounter := 0
	curPos := -1
	lastPos := -1
	for assignedFloor != -1 {
		curPos = status.Pos
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
	allElevators := initAllElevatorStatuses()
	allElevators[curID].IsOnline = true
	allElevators[curID].IsAvailable = true

	chans.AskElevatorForUpdate <- true
	initFloor := <-chans.New_floor
	initDoor := <-chans.Door_status

	allElevators[curID].Pos = initFloor
	allElevators[curID].DoorOpen = initDoor
	assignedOrder := ButtonEvent_message{-1, messages.UNDEFINED}

	orderTimeOut := time.Now().Add(-time.Hour)

	go contSend(chans.Send_status, &allElevators)
	go sendOutLightsUpdate(chans.UpdateLampMessage, &allElevators[curID])

	for {
		select {
		case buttonEvent := <-chans.Button_press:
			fmt.Println("--------Button pressed------")
			updateOrderListButton(buttonEvent, &allElevators[curID])
			// if buttonEvent.Floor == allElevators[curID].Pos && time.Since(orderTimeOut) > config.OrderTimeOut{
			// 	updateOrderListCompleted(&allElevators, buttonEvent, &orderTimeOut) //This causes hall orders to get lost when it is at floor and hall order up and down is pressed
			// }
			fmt.Println("-> Status sendt: Buttonpress,", buttonEvent.Button)

			go sendOutStatus(chans.Send_status, allElevators[curID])

		case elevatorStatus := <-chans.Received_elevator_update: // new update
			// update own orderlist and otherElev with the incomming elevatorStatus         COMMENTED OUT BEACUSE OF TESTING WITHOUT NETWORK MODULE

			if elevatorStatus.ID != curID {
				updateOtherElev(elevatorStatus, &allElevators)
				updateOrderListOther(elevatorStatus, &allElevators)
				// fmt.Println("<- Recived status")
			}

		case floor := <-chans.New_floor:
			updateElevatorStatusFloor(floor, &allElevators)
			updateOrderListCompleted(&allElevators, assignedOrder, &orderTimeOut)
			go sendOutStatus(chans.Send_status, allElevators[curID])
			fmt.Println("-> Status sendt: New floor:", floor)

		case isOpen := <-chans.Door_status:
			updateElevatorStatusDoor(isOpen, &allElevators)
			if isOpen == true {
				// time.Sleep(time.Second * 3)
				updateOrderListCompleted(&allElevators, assignedOrder, &orderTimeOut)

			}
			go sendOutStatus(chans.Send_status, allElevators[curID])
			fmt.Println("-> Status sendt: Door: ", isOpen)

		case IsObstructed := <-chans.DoorObstructed:
			allElevators[curID].IsObstructed = IsObstructed //TODO burde denne abstraheres ut i en egen funksjon
			go sendOutStatus(chans.Send_status, allElevators[curID])

		case <-time.After(500 * time.Millisecond):
		}

		// allElevators[curID].IsAvailable = true
		newAssignedOrder := goToFloor(&allElevators)
		// printElevatorStatus(allElevators[curID])
		//if newAssignedOrder.Floor == -1 {
		//	assignedOrder = newAssignedOrder
		//}
		// it isn't able to complete a order if it recives a matching order ass the previous while its TimedOut
		// if newAssignedOrder.Floor != -1 && newAssignedOrder != assignedOrder && allElevators[curID].DoorOpen && !allElevators[curID].IsObstructed && time.Since(orderTimeOut) > config.OrderTimeOut {
		if newAssignedOrder.Floor != -1 && allElevators[curID].DoorOpen && !allElevators[curID].IsObstructed && time.Since(orderTimeOut) > config.OrderTimeOut {
			assignedOrder = newAssignedOrder
			//fmt.Println("Sending out order:", assignedFloor)

			go sendingElevatorToFloor(chans.Go_to_floor, assignedOrder.Floor)

			if assignedOrder.Floor == allElevators[curID].Pos {
				updateOrderListCompleted(&allElevators, assignedOrder, &orderTimeOut)
			}
			//allElevators[curID].IsAvailable = false
			//printElevatorStatus(allElevators[curID])
			go sendOutStatus(chans.Send_status, allElevators[curID])
			fmt.Println("-> Status sendt: Go to floor ", assignedOrder.Floor)
			go checkIfEngineFails(assignedOrder.Floor, &allElevators[curID], chans.Send_status)
			printElevatorStatus(allElevators[curID])

		}

		checkIfElevatorOffline(&allElevators)
	}
}

// Functions for printing out the Elevator staus and OrderList
func printOrderList(list [orderListLength]HallOrder) {
	for i := 0; i < len(list); i++ {
		fmt.Println("Floor: ", list[i].Floor, " |direction: ", list[i].Direction, " |Order: ", list[i].HasOrder, " |Version", list[i].VersionNum, "Costs", list[i].Costs)
	}
}

func printElevatorStatus(status ElevatorStatus) {
	fmt.Println("ID:", status.ID, " |Currentfloor:", status.Pos, "| Door open:", status.DoorOpen, "|avaliable:", status.IsAvailable, "|isOnline:", status.IsOnline, "|isObstructed:", status.IsObstructed, "\nCabOrders:", status.CabOrders)

	printOrderList(status.OrderList)
	fmt.Println("-----------------------------------------------")
}
