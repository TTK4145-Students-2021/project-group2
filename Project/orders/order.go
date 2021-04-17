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
var thisID = config.ID

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
	list[thisID].DoorOpen = value
}
func updateElevatorStatusFloor(pos int, list *[numElevators]ElevatorStatus) {
	list[thisID].Pos = pos
}

func costFunction(num int, list *[numElevators]ElevatorStatus) {
	curOrder := list[thisID].OrderList[num]
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

		list[thisID].OrderList[num].Costs[i] = cost
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
	for i := 0; i < len(list[thisID].OrderList); i++ {
		if list[thisID].OrderList[i].HasOrder { // TODO: rename list
			costFunction(i, list)
		} else {
			for j := 0; j < numElevators; j++ {
				list[thisID].OrderList[i].Costs[j] = 10000
			}
		}
	}
	pickedOrder := ButtonEvent_message{-1, -1} // buttonType -1 = UNDEFINED

	//	checking if the elevator has any cabOrders and if assigning that ass the goToFloor
	shortestCabDistance := numFloors * 10 // random variable larger than config.numfloors
	for floor := 0; floor < numFloors; floor++ {
		if list[thisID].CabOrders[floor] == true {
			curCabDistance := distanceFromOrder(list[thisID].Pos, floor)
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
		costMatrix[i] = list[thisID].OrderList[i].Costs
	}
	lowestCost := 500

	for i := 0; i < numElevators; i++ {

		orderIdx, ID := minumumRowCol(costMatrix)
		cost := costMatrix[orderIdx][ID]
		if ID == thisID && cost < lowestCost && list[thisID].OrderList[orderIdx].HasOrder {
			pickedOrder.Floor = orderListIdxToFloor(orderIdx)
			pickedOrder.Button = ButtonType_msg(list[thisID].OrderList[orderIdx].Direction)
			return pickedOrder
		} else {
			for i := range list[thisID].OrderList {
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

func orderListIdxToBtnTyp(idx int) ButtonType_msg {
	if idx < numFloors-1 {
		return BT_HallUp
	} else {
		return BT_HallDown
	}
}






func floorToOrderListIdx(floor int, dir messages.ButtonType_msg) int {
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

	for i := 0; i < len(list[thisID].OrderList); i++ {

		// the first if sentences checks for faiure is the version controll
		if incomingStatus.OrderList[i].VersionNum == list[thisID].OrderList[i].VersionNum && incomingStatus.OrderList[i].HasOrder != list[thisID].OrderList[i].HasOrder {
			fmt.Println("There is smoething worng with the version controll of the order system")
			list[thisID].OrderList[i].HasOrder = true
			list[thisID].OrderList[i].VersionNum += 1
		}
		// this if statement checks if the order should be updated and updates it
		if incomingStatus.OrderList[i].VersionNum > list[thisID].OrderList[i].VersionNum {
			list[thisID].OrderList[i].HasOrder = incomingStatus.OrderList[i].HasOrder
			list[thisID].OrderList[i].VersionNum = incomingStatus.OrderList[i].VersionNum
		}
	}
}

func updateOrderListCompleted(list *[numElevators]ElevatorStatus, assignedOrder ButtonEvent_message, orderTimeOut *time.Time) {
	curFloor := list[thisID].Pos

	// checks if the elevator has a cab call for the current floor
	if list[thisID].CabOrders[curFloor] == true && list[thisID].DoorOpen == true {
		list[thisID].CabOrders[curFloor] = false
		*orderTimeOut = time.Now()
		return
	}

	orderIdx := floorToOrderListIdx(assignedOrder.Floor, assignedOrder.Button)

	if list[thisID].OrderList[orderIdx].HasOrder && list[thisID].OrderList[orderIdx].Floor == curFloor {
		list[thisID].OrderList[orderIdx].HasOrder = false
		list[thisID].OrderList[orderIdx].VersionNum += 1
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
	for elevID := range allElevatorSatuses {
		if time.Since(allElevatorSatuses[elevID].Timestamp) > config.OfflineTimeout && elevID != thisID {
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
		go sendOutStatus(channel, status[thisID])
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


	assignedOrder := ButtonEvent_message{-1, messages.UNDEFINED}
	executeOrder := ButtonEvent_message{-1, messages.UNDEFINED}

	orderTimeOut := time.Now().Add(-time.Hour)


	go contSend(chans.Send_status, &allElevatorSatuses)
	go sendOutLightsUpdate(chans.UpdateLampMessage, &allElevatorSatuses[thisID])

	for {
		select {
		case buttonEvent := <-chans.Button_press:
			fmt.Println("--------Button pressed------")
			updateOrderListButton(buttonEvent, &allElevatorSatuses[thisID])
			// if buttonEvent.Floor == allElevatorSatuses[thisID].Pos && time.Since(orderTimeOut) > config.OrderTimeOut{
			// 	updateOrderListCompleted(&allElevatorSatuses, buttonEvent, &orderTimeOut) //This causes hall orders to get lost when it is at floor and hall order up and down is pressed
			// }
			fmt.Println("-> Status sendt: Buttonpress,", buttonEvent.Button)

			go sendOutStatus(chans.Send_status, allElevatorSatuses[thisID])

		case elevatorStatus := <-chans.Received_elevator_update: // new update
			// update own orderlist and otherElev with the incomming elevatorStatus         COMMENTED OUT BEACUSE OF TESTING WITHOUT NETWORK MODULE

			if elevatorStatus.ID != thisID {
				updateOtherElev(elevatorStatus, &allElevatorSatuses)
				updateOrderListOther(elevatorStatus, &allElevatorSatuses)
				// fmt.Println("<- Recived status")
			}

		case floor := <-chans.New_floor:
			updateElevatorStatusFloor(floor, &allElevatorSatuses)
			updateOrderListCompleted(&allElevatorSatuses, assignedOrder, &orderTimeOut)
			go sendOutStatus(chans.Send_status, allElevatorSatuses[thisID])
			fmt.Println("-> Status sendt: New floor:", floor)

		case isOpen := <-chans.Door_status:
			updateElevatorStatusDoor(isOpen, &allElevatorSatuses)
			if isOpen == true {
				// time.Sleep(time.Second * 3)
				updateOrderListCompleted(&allElevatorSatuses, assignedOrder, &orderTimeOut)

			}
			go sendOutStatus(chans.Send_status, allElevatorSatuses[thisID])
			fmt.Println("-> Status sendt: Door: ", isOpen)

		case IsObstructed := <-chans.DoorObstructed:
			allElevatorSatuses[thisID].IsObstructed = IsObstructed //TODO burde denne abstraheres ut i en egen funksjon
			go sendOutStatus(chans.Send_status, allElevatorSatuses[thisID])

		case <-time.After(500 * time.Millisecond):
		}

		newAssignedOrder := goToFloor(&allElevatorSatuses)
	
		if newAssignedOrder.Floor != -1 && allElevatorSatuses[thisID].DoorOpen && !allElevatorSatuses[thisID].IsObstructed && time.Since(orderTimeOut) > config.OrderTimeOut {
			assignedOrder = newAssignedOrder

			if assignedOrder != executeOrder{
				go sendingElevatorToFloor(chans.Go_to_floor, assignedOrder.Floor)
				go sendOutStatus(chans.Send_status, allElevatorSatuses[thisID])
				fmt.Println("-> Status sendt: Go to floor ", assignedOrder.Floor)
				go checkIfEngineFails(assignedOrder.Floor, &allElevatorSatuses[thisID], chans.Send_status)
				executeOrder = assignedOrder
			}
			if assignedOrder.Floor == allElevatorSatuses[thisID].Pos {
				updateOrderListCompleted(&allElevatorSatuses, assignedOrder, &orderTimeOut)
			}
			printElevatorStatus(allElevatorSatuses[thisID])
		}
		checkIfElevatorOffline(&allElevatorSatuses)
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
