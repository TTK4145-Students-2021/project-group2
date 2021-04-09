// run main with as follows:
// go run main.go [port] [elevatorID]

package main

import (
	"fmt"

	"./config"
	"./elevio"
	"./messages"
	"./orders"
)

func main() {
	//var HallOrder = make(chan elevio.ButtonEvent)
	//var ReceivedElevatorUpdate = make(chan orders.ElevatorStatus)
	//var NewFloor = make(chan int)
	//var IsDoorOpen = make(chan bool)

	fmt.Println(config.NumFloors, config.NumElevators, config.Port, config.ID)

	// Setup all the channels we need
	elev_available := make(chan bool)                        // Door_open*
	current_floor := make(chan int)                          //new  floor
	button_action := make(chan messages.ButtonEvent_message) //pressed button
	order := make(chan int)                                  // goTo floor
	order_response := make(chan error)                       //TODO: Order response can be a const int instead

	// Bundle controller channels in a struct
	ctrChans := elevio.ControllerChannels{
		Elev_available:         elev_available, // TODO: Rename
		Current_floor:          current_floor,
		Redirect_button_action: button_action,
		Receive_order:          order,
		Respond_order:          order_response,
	}

	//config.InitConfig()

	controller := elevio.NewController(ctrChans)
	go controller.Run()

	//elev := controller.GetElevator()

	//define channels

	fmt.Println("Should start running orders")
	//elevio.RunElevator("localhost:"+config.Port, config.NumFloors)
	go orders.RunOrders(button_action,
		//received_elevator_update <- chan ElevatorStatus,    //Network communication
		current_floor,
		elev_available,
		// send_status chan <- ElevatorStatus, //Network communication
		order)
	
	for{	
	}

	//newOrder := make(chan elevio.ButtonEvent)
	//assignedOrder := make(chan elevio.ButtonEvent)

	// TestElevatorMachine() ??

}
