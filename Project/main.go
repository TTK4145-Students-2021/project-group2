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
	getElevatorUpdate := make(chan bool)
	doorOpen := make(chan bool)
	currentFloor := make(chan int)                          //new  floor
	buttonAction := make(chan messages.ButtonEvent_message) //pressed button
	order := make(chan int)                                 // goTo floor
	orderResponse := make(chan error)                       //TODO: Order response can be a const int instead
	controllerReady := make(chan bool)                      // Check when controller is done

	// Bundle controller channels in a struct
	ctrChans := elevator.ControllerChannels{
		DoorOpen:              doorOpen,
		CurrentFloor:          currentFloor,
		RedirectButtonAction:  buttonAction,
		ReceiveOrder:          order,
		RespondOrder:          orderResponse,
		ElevatorUpdateRequest: getElevatorUpdate,
		ControllerReady:       controllerReady,
	}

	controller := elevator.NewController(ctrChans)
	go controller.Run()
	<-controllerReady // Check when controller is ready

	fmt.Println("Should start running orders")
	go orders.RunOrders(buttonAction,
		//received_elevator_update //<- chan ElevatorStatus,    //Network communication
		currentFloor,
		doorOpen,
		//send_status chan <- ElevatorStatus, //Network communication
		order,
		getElevatorUpdate)

	for {
	}

}
