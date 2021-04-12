// run main with as follows:
// go run main.go [port] [elevatorID]

package main

import (
	"fmt"

	"./config"
	"./elevator"
	"./messages"
	"./orders"
	"./network/bcast"
)

func main() {

	fmt.Println(config.NumFloors, config.NumElevators, config.Port, config.ID)

	// Setup all the channels we need
	getElevatorUpdate := make(chan bool)
	elevAvailable := make(chan bool)                        // Door_open*
	currentFloor := make(chan int)                          //new  floor
	buttonAction := make(chan messages.ButtonEvent_message) //pressed button
	order := make(chan int)                                 // goTo floor
	orderResponse := make(chan error)                       //TODO: Order response can be a const int instead
	controllerReady := make(chan bool)                      // Check when controller is done
	elevatorStatusTx := make(chan messages.ElevatorStatus) // Change to orders.ElevatorStatus
	elevatorStatusRx := make(chan messages.ElevatorStatus) // Change to orders.ElevatorStatus
	

	// Bundle controller channels in a struct
	ctrChans := elevator.ControllerChannels{
		Elev_available:         elevAvailable, // TODO: Rename
		Current_floor:          currentFloor,
		Redirect_button_action: buttonAction,
		Receive_order:          order,
		Respond_order:          orderResponse,
		ElevatorUpdateRequest:  getElevatorUpdate,
		ControllerReady:        controllerReady,
	}

	controller := elevator.NewController(ctrChans)
	go controller.Run()
	<-controllerReady // Check when controller is ready

	//define channels

	fmt.Println("Should start running orders")
	//elevator.RunElevator("localhost:"+config.Port, config.NumFloors)
	go orders.RunOrders(buttonAction,
		elevatorStatusRx, //<- chan ElevatorStatus,    //Network communication
		currentFloor,
		elevAvailable,
		elevatorStatusTx, //Network communication
		order,
		getElevatorUpdate)

	go bcast.Transmitter(config.BcastPort, elevatorStatusTx)
	bcast.Receiver(config.BcastPort, elevatorStatusRx)
	

	//newOrder := make(chan elevator.ButtonEvent)
	//assignedOrder := make(chan elevator.ButtonEvent)

	// TestElevatorMachine() ??

}
