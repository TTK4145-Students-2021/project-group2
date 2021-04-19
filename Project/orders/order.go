package orders

import (
	//"os"
	"errors"
	"fmt"
	"time"

	"../config"
	"../messages"
)

// assumes buttonEvent as input
// bcast elevID, targetFloor, dir
/*
type HallOrder struct {
	HasOrder   bool
	Floor      int
	Direction  int //up = 0, down = 1
	VersionNum int
	Costs      [config.NumElevators]int
	TimeStamp  time.Time
}

type ElevatorStatus struct {
	ID          int
	Pos         int
	OrderList   [config.NumFloors*2 - 2]HallOrder //
	Dir         int                               //up = 0, down = 1, stop = 2
	IsOnline    bool
	DoorOpen    bool
	CabOrders   [config.NumFloors]bool
	IsAvailable bool
}
*/
//type HallOrder = messages.HallOrder
//type ElevatorStatus = messages.ElevatorStatus

// type Direction int

// const (
// 	UP   Direction = 0
// 	DOWN Direction = 1
// )

type OrderChannels struct {
	Button_press             <-chan messages.ButtonEvent_message //Elevator communiaction
	Received_elevator_update <-chan messages.ElevatorStatus      //Network communication
	New_floor                <-chan int                          //Elevator communiaction
	Door_status              <-chan bool                         //Elevator communiaction
	Send_status              chan<- messages.ElevatorStatus      //Network communication
	Go_to_floor              chan<- int                          //Elevator communiaction
	AskElevatorForUpdate     chan<- bool
	DoorObstructed           <-chan bool
	UpdateLampMessage        chan<- messages.LampUpdate_message
}

type ButtonEvent_message = messages.ButtonEvent_message
type ButtonType_msg = messages.ButtonType_msg

func initOrderList() [config.NumFloors*2 - 2]messages.HallOrder {
	OrderList := [config.NumFloors*2 - 2]messages.HallOrder{}

	// initalizing HallUp
	initialCosts := [config.NumElevators]int{}
	for i := 0; i < config.NumElevators; i++ {
		initialCosts[i] = 10000
	}
	for i := 0; i < (config.NumFloors - 1); i++ {
		OrderList[i] = messages.HallOrder{
			HasOrder:   false,
			Floor:      i,
			Direction:  0, //up = 0, down = 1
			VersionNum: 0,
			Costs:      initialCosts,
			TimeStamp:  time.Now(),
		}
	}

	// initalizing HallDown
	for i := config.NumFloors - 1; i < (config.NumFloors*2 - 2); i++ {
		OrderList[i] = messages.HallOrder{
			HasOrder:   false,
			Floor:      i + 2 - config.NumFloors,
			Direction:  1, //up = 0, down = 1
			VersionNum: 0,
			Costs:      initialCosts,
			TimeStamp:  time.Now(),
		}
	}
	return OrderList
}

func initAllElevatorStatuses() [config.NumElevators]messages.ElevatorStatus {

	CabOrders := [config.NumFloors]bool{}
	OrderList := initOrderList()

	listOfElevators := [config.NumElevators]messages.ElevatorStatus{}

	for i := 0; i < (config.NumElevators); i++ {
		Status := messages.ElevatorStatus{
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
func updateOrderListButton(btnEvent messages.ButtonEvent_message, thisElevator *messages.ElevatorStatus) {

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
func updateOtherElev(incomingStatus messages.ElevatorStatus, list *[config.NumElevators]messages.ElevatorStatus) {
	list[incomingStatus.ID] = incomingStatus //accountign for zero indexing
}
func updateElevatorStatusDoor(value bool, list *[config.NumElevators]messages.ElevatorStatus) {
	list[config.ID].DoorOpen = value
}
func updateElevatorStatusFloor(pos int, list *[config.NumElevators]messages.ElevatorStatus) {
	list[config.ID].Pos = pos
}

func costFunction(num int, list *[config.NumElevators]messages.ElevatorStatus) {
	curOrder := list[config.ID].OrderList[num]
	for i := 0; i < config.NumElevators; i++ { //i is already 0 indexed

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
		for j := 0; j < config.NumFloors; j++ { //organize into function
			if elevator.CabOrders[j] == true {
				cost += 3000
				break
			}
		}
		//cost += waitingTimePenalty(curOrder)
		cost += distanceFromOrder(list[i].Pos, curOrder.Floor)

		list[config.ID].OrderList[num].Costs[i] = cost
	}
}

func waitingTimePenalty(curOrder messages.HallOrder) int {
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

func goToFloor(list *[config.NumElevators]messages.ElevatorStatus) ButtonEvent_message {
	// runs through the costFunction for all true orders
	for i := 0; i < len(list[config.ID].OrderList); i++ {
		if list[config.ID].OrderList[i].HasOrder { // TODO: rename list
			costFunction(i, list)
		} else {
			for j := 0; j < config.NumElevators; j++ {
				list[config.ID].OrderList[i].Costs[j] = 10000
			}
		}
	}
	pickedOrder := ButtonEvent_message{-1, -1} // buttonType -1 = UNDEFINED

	//	checking if the elevator has any cabOrders and if assigning that ass the goToFloor
	shortestCabDistance := config.NumFloors * 10 // random variable larger than config.numfloors
	for floor := 0; floor < config.NumFloors; floor++ {
		if list[config.ID].CabOrders[floor] == true {
			curCabDistance := distanceFromOrder(list[config.ID].Pos, floor)
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

	costMatrix := make([][config.NumElevators]int, config.NumFloors*2-2)
	for i := range costMatrix {
		costMatrix[i] = list[config.ID].OrderList[i].Costs
	}
	lowestCost := 500

	for i := 0; i < config.NumElevators; i++ {

		orderIdx, ID := minumumRowCol(costMatrix)
		cost := costMatrix[orderIdx][ID]
		if ID == config.ID && cost < lowestCost && list[config.ID].OrderList[orderIdx].HasOrder {
			pickedOrder.Floor = orderListIdxToFloor(orderIdx)
			pickedOrder.Button = ButtonType_msg(list[config.ID].OrderList[orderIdx].Direction)
			return pickedOrder
		} else {
			for i := range list[config.ID].OrderList {
				costMatrix[i][ID] = 10000
			}
			for i := 0; i < config.NumElevators; i++ {
				costMatrix[orderIdx][i] = 10000
			}
		}
	}
	return pickedOrder //return the floor that the elevator should go
}

func orderListIdxToFloor(idx int) int {
	if idx < config.NumFloors-1 {
		return idx
	} else {
		return idx - config.NumFloors + 2
	}
}

func floorToOrderListIdx(floor int, dir messages.ButtonType_msg) int {
	if dir == messages.BT_HallUp {
		// fmt.Println("DOWN")
		return floor
	} else {
		// fmt.Println("UP:", floor + config.NumFloors - 2)
		return floor + config.NumFloors - 2
	}
}

func minumumRowCol(list [][config.NumElevators]int) (int, int) {
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

func updateOrderListOther(incomingStatus messages.ElevatorStatus, list *[config.NumElevators]messages.ElevatorStatus) {

	for i := 0; i < len(list[config.ID].OrderList); i++ {

		// the first if sentences checks for faiure is the version controll
		if incomingStatus.OrderList[i].VersionNum == list[config.ID].OrderList[i].VersionNum && incomingStatus.OrderList[i].HasOrder != list[config.ID].OrderList[i].HasOrder {
			fmt.Println("There is smoething worng with the version controll of the order system")
			list[config.ID].OrderList[i].HasOrder = true
			list[config.ID].OrderList[i].VersionNum += 1
		}
		// this if statement checks if the order should be updated and updates it
		if incomingStatus.OrderList[i].VersionNum > list[config.ID].OrderList[i].VersionNum {
			list[config.ID].OrderList[i].HasOrder = incomingStatus.OrderList[i].HasOrder
			list[config.ID].OrderList[i].VersionNum = incomingStatus.OrderList[i].VersionNum
		}
	}
}

func updateOrderListCompleted(list *[config.NumElevators]messages.ElevatorStatus, assignedOrder ButtonEvent_message, orderTimeOut *time.Time) {
	curFloor := list[config.ID].Pos

	// checks if the elevator has a cab call for the current floor
	if list[config.ID].CabOrders[curFloor] == true && list[config.ID].DoorOpen == true {
		list[config.ID].CabOrders[curFloor] = false
		*orderTimeOut = time.Now()
		return
	}

	orderIdx := floorToOrderListIdx(assignedOrder.Floor, assignedOrder.Button)

	if list[config.ID].OrderList[orderIdx].HasOrder && list[config.ID].OrderList[orderIdx].Floor == curFloor {
		list[config.ID].OrderList[orderIdx].HasOrder = false
		list[config.ID].OrderList[orderIdx].VersionNum += 1
		*orderTimeOut = time.Now()
		fmt.Println("Hall order compleete:", orderListIdxToFloor(orderIdx))
	}
}

func sendOutLightsUpdate(sendLampUpdate chan<- messages.LampUpdate_message, status *messages.ElevatorStatus) {
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

func checkIfElevatorOffline(allElevators *[config.NumElevators]messages.ElevatorStatus) {
	for elevID := range allElevators {
		if time.Since(allElevators[elevID].Timestamp) > config.OfflineTimeout && elevID != config.ID {
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

func sendOutStatus(channel chan<- messages.ElevatorStatus, status messages.ElevatorStatus) {
	status.Timestamp = time.Now()
	channel <- status
}
func sendingElevatorToFloor(channel chan<- int, goToFloor int) {
	channel <- goToFloor
}

func contSend(channel chan<- messages.ElevatorStatus, list *[config.NumElevators]messages.ElevatorStatus) {
	for {
		status := *list
		go sendOutStatus(channel, status[config.ID])
		time.Sleep(config.BcastIntervall)
	}
}
func checkIfEngineFails(assignedFloor int, status *messages.ElevatorStatus, send_status chan<- messages.ElevatorStatus) {
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
	allElevators[config.ID].IsOnline = true
	allElevators[config.ID].IsAvailable = true

	backup := NewBackup("CabOrderBackup.txt", config.ID)
	backup.RecoverCabOrders(&allElevators[config.ID].CabOrders)

	chans.AskElevatorForUpdate <- true
	initFloor := <-chans.New_floor
	initDoor := <-chans.Door_status

	allElevators[config.ID].Pos = initFloor
	allElevators[config.ID].DoorOpen = initDoor
	assignedOrder := ButtonEvent_message{-1, messages.UNDEFINED}

	orderTimeOut := time.Now().Add(-time.Hour)

	go contSend(chans.Send_status, &allElevators)
	go sendOutLightsUpdate(chans.UpdateLampMessage, &allElevators[config.ID])

	for {
		select {
		case buttonEvent := <-chans.Button_press:
			fmt.Println("--------Button pressed------")
			updateOrderListButton(buttonEvent, &allElevators[config.ID])
			// if buttonEvent.Floor == allElevators[config.ID].Pos && time.Since(orderTimeOut) > config.OrderTimeOut{
			// 	updateOrderListCompleted(&allElevators, buttonEvent, &orderTimeOut) //This causes hall orders to get lost when it is at floor and hall order up and down is pressed
			// }
			fmt.Println("-> Status sendt: Buttonpress,", buttonEvent.Button)
			go backup.SaveCabOrders(allElevators[config.ID].CabOrders)
			go sendOutStatus(chans.Send_status, allElevators[config.ID])

		case elevatorStatus := <-chans.Received_elevator_update: // new update
			// update own orderlist and otherElev with the incomming elevatorStatus         COMMENTED OUT BEACUSE OF TESTING WITHOUT NETWORK MODULE

			if elevatorStatus.ID != config.ID {
				updateOtherElev(elevatorStatus, &allElevators)
				updateOrderListOther(elevatorStatus, &allElevators)
				// fmt.Println("<- Recived status")
			}

		case floor := <-chans.New_floor:
			updateElevatorStatusFloor(floor, &allElevators)
			updateOrderListCompleted(&allElevators, assignedOrder, &orderTimeOut)
			go backup.SaveCabOrders(allElevators[config.ID].CabOrders)
			go sendOutStatus(chans.Send_status, allElevators[config.ID])
			fmt.Println("-> Status sendt: New floor:", floor)

		case isOpen := <-chans.Door_status:
			updateElevatorStatusDoor(isOpen, &allElevators)
			if isOpen == true {
				// time.Sleep(time.Second * 3)
				updateOrderListCompleted(&allElevators, assignedOrder, &orderTimeOut)
				go backup.SaveCabOrders(allElevators[config.ID].CabOrders)

			}
			go sendOutStatus(chans.Send_status, allElevators[config.ID])
			fmt.Println("-> Status sendt: Door: ", isOpen)

		case IsObstructed := <-chans.DoorObstructed:
			allElevators[config.ID].IsObstructed = IsObstructed //TODO burde denne abstraheres ut i en egen funksjon
			go sendOutStatus(chans.Send_status, allElevators[config.ID])

		case <-time.After(500 * time.Millisecond):
		}

		// allElevators[config.ID].IsAvailable = true
		newAssignedOrder := goToFloor(&allElevators)
		// printElevatorStatus(allElevators[config.ID])
		//if newAssignedOrder.Floor == -1 {
		//	assignedOrder = newAssignedOrder
		//}
		// it isn't able to complete a order if it recives a matching order ass the previous while its TimedOut
		// if newAssignedOrder.Floor != -1 && newAssignedOrder != assignedOrder && allElevators[config.ID].DoorOpen && !allElevators[config.ID].IsObstructed && time.Since(orderTimeOut) > config.OrderTimeOut {
		if newAssignedOrder.Floor != -1 && allElevators[config.ID].DoorOpen && !allElevators[config.ID].IsObstructed && time.Since(orderTimeOut) > config.OrderTimeOut {
			assignedOrder = newAssignedOrder
			//fmt.Println("Sending out order:", assignedFloor)

			go sendingElevatorToFloor(chans.Go_to_floor, assignedOrder.Floor)

			if assignedOrder.Floor == allElevators[config.ID].Pos {
				updateOrderListCompleted(&allElevators, assignedOrder, &orderTimeOut)
			}
			//allElevators[config.ID].IsAvailable = false
			//printElevatorStatus(allElevators[config.ID])
			go sendOutStatus(chans.Send_status, allElevators[config.ID])
			fmt.Println("-> Status sendt: Go to floor ", assignedOrder.Floor)
			go checkIfEngineFails(assignedOrder.Floor, &allElevators[config.ID], chans.Send_status)
			printElevatorStatus(allElevators[config.ID])

		}

		checkIfElevatorOffline(&allElevators)
	}
}

// Functions for printing out the Elevator staus and OrderList
func printOrderList(list [config.NumFloors*2 - 2]messages.HallOrder) {
	for i := 0; i < len(list); i++ {
		fmt.Println("Floor: ", list[i].Floor, " |direction: ", list[i].Direction, " |Order: ", list[i].HasOrder, " |Version", list[i].VersionNum, "Costs", list[i].Costs)
	}
}

func printElevatorStatus(status messages.ElevatorStatus) {
	fmt.Println("ID:", status.ID, " |Currentfloor:", status.Pos, "| Door open:", status.DoorOpen, "|avaliable:", status.IsAvailable, "|isOnline:", status.IsOnline, "|isObstructed:", status.IsObstructed, "\nCabOrders:", status.CabOrders)

	printOrderList(status.OrderList)
	fmt.Println("-----------------------------------------------")
}
