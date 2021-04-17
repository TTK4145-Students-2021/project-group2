package orders

import (
	"fmt"
	"testing"
	"time"

	"../config"
	"../messages"
)

func TestOrders(t *testing.T) {
	fmt.Println("Start")
	// creating 1 ElevatorStatus

	OrderList := [config.NumFloors*2 - 2]HallOrder{}
	for i := 0; i < (config.NumFloors - 1); i++ {
		OrderList[i] = HallOrder{
			HasOrder:   false,
			Floor:      i + 1,
			Direction:  0, //up = 0, down = 1
			VersionNum: 0,
			Costs:      [config.NumElevators]int{},
			TimeStamp:  time.Now(),
		}
	}
	for i := config.NumFloors - 1; i < (config.NumFloors*2 - 2); i++ {
		OrderList[i] = HallOrder{
			HasOrder:   false,
			Floor:      i + 3 - config.NumFloors,
			Direction:  1, //up = 0, down = 1
			VersionNum: 0,
			Costs:      [config.NumElevators]int{},
			TimeStamp:  time.Now(),
		}
	}
	CabOrders := [config.NumFloors]bool{false}

	Status2 := ElevatorStatus{
		ID:          1,
		Pos:         0,
		OrderList:   OrderList,
		Dir:         2,
		IsOnline:    true, //NB isOline is deafult false
		DoorOpen:    false,
		CabOrders:   CabOrders,
		IsAvailable: true, //NB isAvaliable is deafult false
	}
	Status3 := ElevatorStatus{
		ID:          2,
		Pos:         3,
		OrderList:   OrderList,
		Dir:         2,
		IsOnline:    true, //NB isOline is deafult false
		DoorOpen:    false,
		CabOrders:   CabOrders,
		IsAvailable: true, //NB isAvaliable is deafult false
	}

	//for testing orders module
	button_press := make(chan messages.ButtonEvent_message)
	received_elevator_update := make(chan ElevatorStatus)
	new_floor := make(chan int)
	door_status := make(chan bool)
	send_status := make(chan ElevatorStatus)
	go_to_floor := make(chan int)

	go RunOrders(button_press,
		received_elevator_update,
		new_floor,
		door_status,
		send_status,
		go_to_floor)

	btnEvent1 := messages.ButtonEvent_message{
		Floor:  1,
		Button: 0, // 0 = up, 1 = down , 2 = hall
	}
	btnEvent2 := messages.ButtonEvent_message{
		Floor:  2,
		Button: 0,
	} // 0 = up, 1 = down , 2 = hall
	new_floor <- 0
	time.Sleep(1000000000)
	button_press <- btnEvent2
	time.Sleep(1000000000)
	return
	button_press <- btnEvent1
	new_floor <- 2
	door_status <- true
	time.Sleep(10000000000)
	door_status <- false
	time.Sleep(5000000000)
	received_elevator_update <- Status2
	received_elevator_update <- Status3
	button_press <- btnEvent2
	time.Sleep(10000000000)

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
