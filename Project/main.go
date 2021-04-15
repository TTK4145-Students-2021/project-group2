// run main with as follows:
// go run main.go [SimPort] [elevatorID]

package main

import (
	"fmt"

	"./config"
	"./elevator"
	"./messages"
	"./network/bcast"
	"./orders"
)

func main() {

	config.InitConfig()
	fmt.Println(config.NumFloors, config.NumElevators, config.Port, config.ID)

	// Setup all the channels we need
	getElevatorUpdate := make(chan bool)
	doorOpen := make(chan bool)
	doorObstructed := make(chan bool)
	currentFloor := make(chan int)                          //new  floor
	buttonAction := make(chan messages.ButtonEvent_message) //pressed button
	order := make(chan int)                                 // goTo floor
	orderResponse := make(chan error)                       //TODO: Order response can be a const int instead
	controllerReady := make(chan bool)                      // Check when controller is done
	elevatorStatusTx := make(chan messages.ElevatorStatus)  // Change to orders.ElevatorStatus
	elevatorStatusRx := make(chan messages.ElevatorStatus)  // Change to orders.ElevatorStatus
	//
	// Bundle controller channels in a struct
	ctrChans := elevator.ControllerChannels{
		DoorOpen:              doorOpen,
		CurrentFloor:          currentFloor,
		RedirectButtonAction:  buttonAction,
		ReceiveOrder:          order,
		RespondOrder:          orderResponse,
		ElevatorUpdateRequest: getElevatorUpdate,
		ControllerReady:       controllerReady,
		DoorObstructed:        doorObstructed,
	}

	controller := elevator.NewController(ctrChans)
	go controller.Run()
	<-controllerReady // Check when controller is ready

	fmt.Println("Should start running orders")

	go bcast.Transmitter(config.BcastPort, elevatorStatusTx)
	go bcast.Receiver(config.BcastPort, elevatorStatusRx)

	orders.RunOrders(buttonAction,
		elevatorStatusRx, //<- chan ElevatorStatus,    //Network communication
		currentFloor,
		doorOpen,
		elevatorStatusTx, //Network communication
		order,
		getElevatorUpdate,
		doorObstructed)

}
