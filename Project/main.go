// run main with as follows:
// go run main.go [port] [elevatorID]

package main

import (
	"fmt"

	"./config"
	"./elevator"
	"./messages"
	"./orders"
)

func main() {

	fmt.Println(config.NumFloors, config.NumElevators, config.Port, config.ID)

	// Setup all the channels we need
	elevAvailable := make(chan bool)                        // Door_open*
	currentFloor := make(chan int)                          //new  floor
	buttonAction := make(chan messages.ButtonEvent_message) //pressed button
	order := make(chan int)                                  // goTo floor
	orderResponse := make(chan error)                       //TODO: Order response can be a const int instead
	controllerReady := make(chan bool)    					 // Check when controller is done

	// Bundle controller channels in a struct
	ctrChans := elevator.ControllerChannels{
		Elev_available        : elevAvailable,    // TODO: Rename
		Current_floor         : currentFloor,
		Redirect_button_action: buttonAction,
		Receive_order         : order,
		Respond_order         : orderResponse,
		ControllerReady       : controllerReady,
	}

	controller := elevator.NewController(ctrChans)
	go controller.Run()
	<-controllerReady // Check when controller is ready

	//elev := controller.GetElevator()

	//define channels

	fmt.Println("Should start running orders")
	//elevator.RunElevator("localhost:"+config.Port, config.NumFloors)
	go orders.RunOrders(buttonAction,
		//received_elevator_update <- chan ElevatorStatus,    //Network communication
		currentFloor,
		elevAvailable,
		// send_status chan <- ElevatorStatus, //Network communication
		order)
	
	for{	
	}

	//newOrder := make(chan elevator.ButtonEvent)
	//assignedOrder := make(chan elevator.ButtonEvent)

	// TestElevatorMachine() ??

}
