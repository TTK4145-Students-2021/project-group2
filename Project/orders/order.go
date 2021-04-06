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

func costFunction(curOrder order, list *[config.NumElevators]ElevatorStatus){
	for i := 0; i < config.NumElevators; i++{//i is already 0 indexed
		cost := 0
		if !list[i].isAvailable || !list[i].isOnline{
			cost += 1000
		}

		for j  // for loop checking if the elevator has a cabOrder
		
		
		floorsAway:= curOrder.floor - list[i].pos
		if floorsAway < 0 {                         // takeing abs of floorsAway
			floorsAway = - floorsAway
		} 
		cost += floorsAway                                   
		

	}

}

/*

func updateOrderListOther(incomingStatus ElevatorStatus, list *[config.NumElevators]ElevatorStatus) {
	This function is suppoesed to update the order list of allElvartors[config.ID] based on the order list from the incoming ElevatorStatus.
	Chech if their is a newer version number on all orders. Updated if there is. 
}


//set the costs for an order in the Orderlist


func pickOrder(list [config.NumElevators]ElevatorStatus) {
	update costs for all the orders
	pick the order where you have the lowest cost. 
}

/*
	assignedFloor = 0     //Defalut assigned floor
	// Itterate through up orders and calcuate all the costs

	for order := range allElevators[config.ID].orderList:    
		if (order.hasOrder):
			order.costs = costFunction(order)  //*****Should you take a copy here????
		else 
			var zeros[]              //TODO: Needs to be finished. Idea is just to inseart a slice of zeros. 
			order.costs = [config.NumElevators]
	
	//TODO:itterate through down orders and calcuate all the costs	
	

}


func runOrders("channel for reciving buttonEvemts from elevator", "channel for reciving changes in Door from elevator"
"channel for reciving changes in current floor from Elevator", "channel for sending "go to floor" (int) too Elevator", 
"channel for reciving incoming Elevatorstatuses","channel for sending out outgoing Elevatorstatuses" ) {

	
	allElevators := initAllElevatorStatuses()
	myElevatorStatus := initMyElevatorStatus()
	allElevators[config.ID-1] = myElevatorStatus

	for {
		select{
		case buttonEvent := <- HallOrder:
			updateOrderListButton(buttonEvent)
		
		case elevatorStatus := <- ReceivedElevatorUpdate: // new update
			// update own orderlist and otherElev with the incomming elevatorStatus
			updateOtherElev(elevatorStatus, &allElevators)
			updateOrderListOther()

		case floor := <- NewFloor:
			updateElevatorStatusFloor(floor)
	
		case isOpen := <-isDoorOpen: //recives a bool value
			updateElevatorStatusDoor(isOpen)
			if(isOpen == true){
				updateOrderListCompleted(myElevatorStatus.floor)
			}
		}  
	}                                                                                                                                                                                                              
	pickOrder()
}
*/




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
