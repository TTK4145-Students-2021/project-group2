package orders

import (
	"fmt"
	"testing"
	"../elevio"
	"../config"
	"time"
)

func TestOrders(t *testing.T) {
	fmt.Println("Start")
	// creating 1 ElevatorStatus
	
	OrderList := [config.NumFloors*2-2]order{} 
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
	
	

	Status2 := ElevatorStatus{
		ID : 2,
		pos : 0,
		orderList :OrderList,  
		dir : elevio.MD_Stop,
		isOnline : true,        //NB isOline is deafult false
		doorOpen : false,
		cabOrder: CabOrder,
		isAvailable : true,        //NB isAvaliable is deafult false
	}
	Status3 := ElevatorStatus{
		ID : 3,
		pos : 0,
		orderList :OrderList,  
		dir : elevio.MD_Stop,
		isOnline : true,        //NB isOline is deafult false
		doorOpen : false,
		cabOrder: CabOrder,
		isAvailable : true,        //NB isAvaliable is deafult false
	}
	

	//for testing orders module
	// for tesing
	button_press := make(chan elevio.ButtonEvent)
	received_elevator_update := make(chan ElevatorStatus)
	new_floor := make(chan int)
	door_status := make(chan bool)
	send_status := make(chan ElevatorStatus)
	go_to_floor := make(chan int)

	go runOrders(button_press,
		received_elevator_update,
		new_floor,
		door_status,
		send_status,
		go_to_floor)

	
	btnEvent1 := elevio.ButtonEvent{
		Floor : 3,
		Button : 0, // 0 = up, 1 = down , 2 = hall
	}
	btnEvent2 := elevio.ButtonEvent{
		Floor : 1,
		Button : 0,
	 } // 0 = up, 1 = down , 2 = hall

	new_floor <- 2
	time.Sleep(5000000000)
	received_elevator_update <- Status2
	received_elevator_update <- Status3
	time.Sleep(5000000000)
	button_press <- btnEvent1
	time.Sleep(5000000000)
	new_floor <- 3
	time.Sleep(2000000000)
	door_status <- true
	time.Sleep(2000000000)
	button_press <- btnEvent2
	time.Sleep(2000000000)



	/*
	btnEvent2 := elevio.ButtonEvent{
		Floor : 1,
		Button : 2, // 0 = up, 1 = down , 2 = hall
	}
	btnEvent3 := elevio.ButtonEvent{
		Floor : 4,
		Button : 1, // 0 = up, 1 = down , 2 = hall
	}
	dummyStatus := allElevators[2]
	
	updateOrderListButton(btnEvent2, &dummyStatus)
	updateOrderListButton(btnEvent3, &dummyStatus)

	updateOtherElev(dummyStatus, &allElevators)
	*/
	// creating a btnEvent
}





