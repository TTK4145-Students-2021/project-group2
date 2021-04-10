package orders

import (
	//"os"
	"time"
	"fmt"
	"../config"
	"../messages"
)

// assumes buttonEvent as input
// bcast elevID, targetFloor, dir

type order struct {
	HasOrder bool
	Floor int
	Direction int	//up = 0, down = 1
	VersionNum int
	Costs [config.NumElevators] int
	TimeStamp time.Time
}

type ElevatorStatus struct {
	ID  int
	Pos int
	OrderList [config.NumFloors*2-2] order  // 
	Dir int   //up = 0, down = 1, stop = 2
	IsOnline bool
	DoorOpen bool
	CabOrder[config.NumFloors] bool
	IsAvailable bool
}

func returnOrder() {
	// Dummy function to test elevatorController
}

func initAllElevatorStatuses() [config.NumElevators]ElevatorStatus {

	OrderList := [config.NumFloors*2-2]order{} 

	//files in the Orderlist with fake defualt order statuses
	for i := 0; i < (config.NumFloors-1); i++ {
		OrderList[i] = order{
			HasOrder : false,
			Floor : i+1,
			Direction : 0,	//up = 0, down = 1
			VersionNum : 0,
			Costs : [config.NumElevators]int{},
			TimeStamp : time.Now(),
		}
	}
	for i := config.NumFloors-1 ; i < (config.NumFloors*2-2); i++ {
		OrderList[i] = order{
			HasOrder : false,
			Floor : i+3-config.NumFloors,
			Direction : 1,	//up = 0, down = 1
			VersionNum : 0,
			Costs : [config.NumElevators]int{},
			TimeStamp : time.Now(),
		}
	}
	CabOrder := [config.NumFloors]bool{false}


	listOfElevators :=[config.NumElevators]ElevatorStatus{}

	for i := 1; i < (config.NumElevators+1); i++{
		Status := ElevatorStatus{
			ID : i,
			Pos : 1,
			OrderList :OrderList,  
			Dir : 2,
			IsOnline : false,        //NB isOline is deafult false
			DoorOpen : false,
			CabOrder: CabOrder,
			IsAvailable : false,        //NB isAvaliable is deafult false
		}
		listOfElevators[i-1] = Status
	}
	
	return listOfElevators
}



//updates your orderList
/*
func updateOrderListButton(btnEvent elevio.ButtonEvent, status *ElevatorStatus) {
	//if 	
	if btnEvent.Button == elevio.BT_Cab {
		status.cabOrder[btnEvent.Floor-1] = true
	} else {
		placeInList := 0
		if btnEvent.Button == elevio.BT_HallUp{
			placeInList = btnEvent.Floor - 1
		} else {
			placeInList = config.NumElevators + btnEvent.Floor - 2
		}	
		if status.orderList[placeInList].hasOrder == false {
			fmt.Println("im here")
			status.orderList[placeInList].hasOrder = true
		}
	}
}
*/
func updateOrderListButton(btnEvent messages.ButtonEvent_message, status *ElevatorStatus) {
	//if 	
	if btnEvent.Button == messages.BT_Cab {
		status.CabOrder[btnEvent.Floor] = true
	} else {
		placeInList := 0
		if btnEvent.Button == messages.BT_HallUp{
			placeInList = btnEvent.Floor
		} else {
			placeInList = config.NumElevators + btnEvent.Floor - 1
		}	
		if status.OrderList[placeInList].HasOrder == false {
			status.OrderList[placeInList].HasOrder = true
			status.OrderList[placeInList].VersionNum += 1
		}
	}
}

//Updates the information list with the incoming update
func updateOtherElev(incomingStatus ElevatorStatus, list *[config.NumElevators]ElevatorStatus)  {
	list[incomingStatus.ID-1] = incomingStatus    //accountign for zero indexing
}
func updateElevatorStatusDoor(value bool,list *[config.NumElevators]ElevatorStatus){
	list[config.ID-1].DoorOpen = value
}
func updateElevatorStatusFloor(pos int,list *[config.NumElevators]ElevatorStatus){
	list[config.ID-1].Pos = pos
}

func costFunction(num int, list *[config.NumElevators]ElevatorStatus){
	curOrder := list[config.ID-1].OrderList[num]
	for i := 0; i < config.NumElevators; i++{//i is already 0 indexed
		elevator := list[i]
		cost := 0
		if !elevator.IsAvailable || !elevator.IsOnline{
			cost += 1000
		}

		// for loop checking if the elevator has a cabOrder
		for j := 0; j < config.NumFloors; j++{
			if elevator.CabOrder[j] == true{
				cost += 1000
				break
			}
		}
		
		floorsAway:= curOrder.Floor - list[i].Pos
		if floorsAway < 0 {                         // takeing abs of floorsAway
			floorsAway = - floorsAway
		} 
		cost += floorsAway                                   
		timeWaited := time.Now().Sub(curOrder.TimeStamp)
		//larger than 1 min
		if timeWaited > time.Duration(60000000000){
			cost += -1
		}
		//larger than 3 min
		if timeWaited > time.Duration(180000000000){
			cost += -1
		}
		//fmt.Println(cost)
		list[config.ID-1].OrderList[num].Costs[i] = cost
	}

}


func pickOrder(list *[config.NumElevators]ElevatorStatus) int {


	for i := 0; i < len(list[config.ID-1].OrderList); i++{
		if list[config.ID-1].OrderList[i].HasOrder {
			costFunction(i,list)
		}
	}
	pickedOrder := 0

	// Pick the order where you have the lowest cost
	for i := 0; i < len(list[config.ID-1].OrderList); i++{
		if list[config.ID-1].OrderList[i].HasOrder {

			isTaken := 0
			for j, num := range list[config.ID-1].OrderList[i].Costs{
				lowestCost := 500
				if num < lowestCost{
					isTaken = j + 1
				}
			if isTaken > 0{
				list[isTaken-1].IsAvailable = false
				if isTaken-1 == config.ID-1{
				pickedOrder = list[config.ID-1].OrderList[i].Floor
				}
			}

			}		
		}
	}

	return pickedOrder    //return the floor that the elevator should go to
}



func updateOrderListOther(incomingStatus ElevatorStatus, list *[config.NumElevators]ElevatorStatus) {

	for i := 0; i < len(list[config.ID-1].OrderList); i++{

		// the first if sentences checks for faiure is the version controll
		if incomingStatus.OrderList[i].VersionNum == list[config.ID-1].OrderList[i].VersionNum && incomingStatus.OrderList[i].HasOrder != list[config.ID-1].OrderList[i].HasOrder {
			fmt.Println("There is smoething worng with the version controll of the order system")
			list[config.ID-1].OrderList[i].HasOrder = true
			list[config.ID-1].OrderList[i].VersionNum += 1
		}
		// this if statement checks if the order should be updated and updates it
		if incomingStatus.OrderList[i].VersionNum > list[config.ID-1].OrderList[i].VersionNum{
			list[config.ID-1].OrderList[i].HasOrder = incomingStatus.OrderList[i].HasOrder
			list[config.ID-1].OrderList[i].VersionNum = incomingStatus.OrderList[i].VersionNum
		}
	}
}
func updateOrderListCompleted(list *[config.NumElevators]ElevatorStatus){
	curFloor := list[config.ID-1].Pos
	fmt.Println(curFloor)
	// checkis if the elevator has a cab call for the current floor
	if list[config.ID-1].CabOrder[curFloor-1] == true{
		list[config.ID-1].CabOrder[curFloor-1] = false
	}

	//Note: this foorloop is here beacuse you have to check both up and down for the same floor, which are stored in diffrent paces in the list
	// TODO: remove this foor loop 
	for i := 0; i < len(list[config.ID-1].OrderList); i++{

		//checks if ther is a hall call at the current floor
		if list[config.ID-1].OrderList[i].HasOrder && list[config.ID-1].OrderList[i].Floor == curFloor{
			list[config.ID-1].OrderList[i].HasOrder = false
		}
	}
}

func sendOutStatus(channel chan<- ElevatorStatus ,list [config.NumElevators]ElevatorStatus){
	channel <- list[config.ID-1]
}
func sendingElevatorToFloor(channel chan<- int, goToFloor int){
	channel <- goToFloor
}


/*
func runOrders("channel for reciving buttonEvemts from elevator", "channel for reciving changes in Door from elevator"
"channel for reciving changes in current floor from Elevator", "channel for sending "go to floor" (int) too Elevator", 
"channel for reciving incoming Elevatorstatuses","channel for sending out outgoing Elevatorstatuses" ) {
*/

func RunOrders(button_press <- chan messages.ButtonEvent_message,  //Elevator communiaction
	 //received_elevator_update <- chan ElevatorStatus,    //Network communication
	 new_floor <- chan int,   //Elevator communiaction
	 door_status <- chan bool,   //Elevator communiaction
	 //send_status chan <- ElevatorStatus, //Network communication
	 go_to_floor chan <- int){    //Elevator communiaction
	
	fmt.Println("Order module initializing....")
	allElevators := initAllElevatorStatuses()
	allElevators[config.ID-1].IsOnline = true
	allElevators[config.ID-1].IsAvailable = true

	for {
		select{
		case buttonEvent := <- button_press:
			updateOrderListButton(buttonEvent, &allElevators[config.ID-1])
			// go sendOutStatus(send_status, allElevators)
		
		//case elevatorStatus := <- received_elevator_update: // new update
			// update own orderlist and otherElev with the incomming elevatorStatus
		//	updateOtherElev(elevatorStatus, &allElevators)
		//	updateOrderListOther(elevatorStatus, &allElevators)

		case floor := <- new_floor:
			updateElevatorStatusFloor(floor, &allElevators)
		//	go sendOutStatus(send_status, allElevators)
	
		case isOpen := <-door_status: //recives a bool value
			updateElevatorStatusDoor(isOpen, &allElevators)
			if(isOpen == true){
				updateOrderListCompleted(&allElevators)
			}
		//	go sendOutStatus(send_status, allElevators)

		}
		
		goTo := pickOrder(&allElevators)
		if goTo != 0 {
			fmt.Println("Sending out order:", goTo) 
			go sendingElevatorToFloor(go_to_floor, goTo)
		}
		printElevatorStatus(allElevators[config.ID-1]) 
	}                                                                                                                                                                                                              
}





// Functions for printing out the Elevator staus and OrderList
func printOrderList(list [config.NumFloors*2-2]order){
	fmt.Println("Orderlist")
	for i := 0; i < len(list); i++ {
		fmt.Println("Floor: ", list[i].Floor, " |direction: ", list[i].Direction, " |Order: ",list[i].HasOrder, " |Version", list[i].VersionNum, "Costs", list[i].Costs,"| time", list[i].TimeStamp)
	}
}

func printElevatorStatus(status ElevatorStatus){
	fmt.Println("ID: ", status.ID, " |Currentfloor: ", status.Pos, "|cabOrder ",status.CabOrder)
	printOrderList(status.OrderList)
}

// |Direction: ", status.Dir,
