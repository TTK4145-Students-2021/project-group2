package orders

import (
	//"os"
	"time"
	"fmt"
	"../config"
	"../elevio"
)

// assumes buttonEvent as input
// bcast elevID, targetFloor, dir

type order struct {
	hasOrder bool
	floor int
	direction int	//up = 0, down = 1
	versionNum int
	costs [config.NumElevators] int
	timeStamp time.Time
}

type ElevatorStatus struct {
	ID  int
	pos int
	orderList [config.NumFloors*2-2] order  // 
	dir elevio.MotorDirection
	isOnline bool
	doorOpen bool
	cabOrder[config.NumFloors] bool
	isAvailable bool
}

func returnOrder() {
	// Dummy function to test elevatorController
}

func initAllElevatorStatuses() [config.NumElevators]ElevatorStatus {

	OrderList := [config.NumFloors*2-2]order{} 

	//files in the Orderlist with fake defualt order statuses
	for i := 0; i < (config.NumFloors-1); i++ {
		OrderList[i] = order{
			hasOrder : false,
			floor : i+1,
			direction : 0,	//up = 0, down = 1
			versionNum : 0,
			costs : [config.NumElevators]int{},
			timeStamp : time.Now(),
		}
	}
	for i := config.NumFloors-1 ; i < (config.NumFloors*2-2); i++ {
		OrderList[i] = order{
			hasOrder : false,
			floor : i+3-config.NumFloors,
			direction : 1,	//up = 0, down = 1
			versionNum : 0,
			costs : [config.NumElevators]int{},
			timeStamp : time.Now(),
		}
	}
	CabOrder := [config.NumFloors]bool{false}


	listOfElevators :=[config.NumElevators]ElevatorStatus{}

	for i := 1; i < (config.NumElevators+1); i++{
		Status := ElevatorStatus{
			ID : i,
			pos : 0,
			orderList :OrderList,  
			dir : elevio.MD_Stop,
			isOnline : false,        //NB isOline is deafult false
			doorOpen : false,
			cabOrder: CabOrder,
			isAvailable : false,        //NB isAvaliable is deafult false
		}
		listOfElevators[i-1] = Status
	}
	
	return listOfElevators
}



//updates your orderList

func updateOrderListButton(btnEvent elevio.ButtonEvent, status *ElevatorStatus) {	
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

//Updates the information list with the incoming update
func updateOtherElev(incomingStatus ElevatorStatus, list *[config.NumElevators]ElevatorStatus)  {
	list[incomingStatus.ID-1] = incomingStatus    //accountign for zero indexing
}
func updateElevatorStatusDoor(value bool,list *[config.NumElevators]ElevatorStatus){
	list[config.ID-1].doorOpen = value
}
func updateElevatorStatusFloor(pos int,list *[config.NumElevators]ElevatorStatus){
	list[config.ID-1].pos = pos
}

func costFunction(num int, list *[config.NumElevators]ElevatorStatus){
	curOrder := list[config.ID-1].orderList[num]
	for i := 0; i < config.NumElevators; i++{//i is already 0 indexed
		elevator := list[i]
		cost := 0
		if !elevator.isAvailable || !elevator.isOnline{
			cost += 1000
		}

		// for loop checking if the elevator has a cabOrder
		for j := 0; j < config.NumFloors; j++{
			if elevator.cabOrder[j] == true{
				cost += 1000
				break
			}
		}
		
		floorsAway:= curOrder.floor - list[i].pos
		if floorsAway < 0 {                         // takeing abs of floorsAway
			floorsAway = - floorsAway
		} 
		cost += floorsAway                                   
		timeWaited := time.Now().Sub(curOrder.timeStamp)
		//larger than 1 min
		if timeWaited > time.Duration(60000000000){
			cost += -1
		}
		//larger than 3 min
		if timeWaited > time.Duration(180000000000){
			cost += -1
		}
		//fmt.Println(cost)
		list[config.ID-1].orderList[num].costs[i] = cost
	}

}


func pickOrder(list *[config.NumElevators]ElevatorStatus) int {


	for i := 0; i < len(list[config.ID-1].orderList); i++{
		if list[config.ID-1].orderList[i].hasOrder {
			costFunction(i,list)
		}
	}
	pickedOrder := 0

	// Pick the order where you have the lowest cost
	for i := 0; i < len(list[config.ID-1].orderList); i++{
		if list[config.ID-1].orderList[i].hasOrder {

			isTaken := 0
			for j, num := range list[config.ID-1].orderList[i].costs{
				lowestCost := 500
				if num < lowestCost{
					isTaken = j + 1
				}
			if isTaken > 0{
				list[isTaken-1].isAvailable = false
				if isTaken-1 == config.ID-1{
				pickedOrder = list[config.ID-1].orderList[i].floor
				}
			}

			}		
		}
	}

	return pickedOrder    //return the floor that the elevator should go to
}



func updateOrderListOther(incomingStatus ElevatorStatus, list *[config.NumElevators]ElevatorStatus) {

	for i := 0; i < len(list[config.ID-1].orderList); i++{

		// the first if sentences checks for faiure is the version controll
		if incomingStatus.orderList[i].versionNum == list[config.ID-1].orderList[i].versionNum && incomingStatus.orderList[i].hasOrder != list[config.ID-1].orderList[i].hasOrder {
			fmt.Println("There is smoething worng with the version controll of the order system")
			list[config.ID-1].orderList[i].hasOrder = true
			list[config.ID-1].orderList[i].versionNum += 1
		}
		// this if statement checks if the order should be updated and updates it
		if incomingStatus.orderList[i].versionNum > list[config.ID-1].orderList[i].versionNum{
			list[config.ID-1].orderList[i].hasOrder = incomingStatus.orderList[i].hasOrder
			list[config.ID-1].orderList[i].versionNum = incomingStatus.orderList[i].versionNum
		}
	}
}
func updateOrderListCompleted(list *[config.NumElevators]ElevatorStatus){
	curFloor := list[config.ID-1].pos

	for i := 0; i < len(list[config.ID-1].orderList); i++{
		if list[config.ID-1].orderList[i].hasOrder && list[config.ID-1].orderList[i].floor == curFloor{
			list[config.ID-1].orderList[i].hasOrder = false
			break
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

func runOrders(button_press <- chan elevio.ButtonEvent,
	 received_elevator_update <- chan ElevatorStatus,
	 new_floor <- chan int, 
	 door_status <- chan bool,
	 send_status chan <- ElevatorStatus,
	 go_to_floor chan <- int){
	
	allElevators := initAllElevatorStatuses()
	allElevators[config.ID-1].isOnline = true
	allElevators[config.ID-1].isAvailable = true

	for {
		select{
		case buttonEvent := <- button_press:
			updateOrderListButton(buttonEvent, &allElevators[config.ID-1])
			go sendOutStatus(send_status, allElevators)
		
		case elevatorStatus := <- received_elevator_update: // new update
			// update own orderlist and otherElev with the incomming elevatorStatus
			updateOtherElev(elevatorStatus, &allElevators)
			updateOrderListOther(elevatorStatus, &allElevators)

		case floor := <- new_floor:
			updateElevatorStatusFloor(floor, &allElevators)
			go sendOutStatus(send_status, allElevators)
	
		case isOpen := <-door_status: //recives a bool value
			updateElevatorStatusDoor(isOpen, &allElevators)
			if(isOpen == true){
				updateOrderListCompleted(&allElevators)
			}
			go sendOutStatus(send_status, allElevators)

		}
		
		goTo := pickOrder(&allElevators)
		fmt.Println(goTo) 
		go sendingElevatorToFloor(go_to_floor, goTo)
		printElevatorStatus(allElevators[config.ID-1]) 
	}                                                                                                                                                                                                              
}





// Functions for printing out the Elevator staus and OrderList
func printOrderList(list [config.NumFloors*2-2]order){
	fmt.Println("This is the elevators orderlist")
	for i := 0; i < len(list); i++ {
		fmt.Println("Floor: ", list[i].floor, " |direction: ", list[i].direction, " |Order: ",list[i].hasOrder, " |Version", list[i].versionNum)
		fmt.Println("Costs", list[i].costs,"| time", list[i].timeStamp)
	}
}

func printElevatorStatus(status ElevatorStatus){
	fmt.Println("ID: ", status.ID, " |Currentfloor: ", status.pos, " |going: ", status.dir, " |cabOrder ",status.cabOrder)
	printOrderList(status.orderList)

}
