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
	doorOpen := make(chan bool, 5)
	doorObstructed := make(chan bool,5)
	currentFloor := make(chan int, 10)                          //new  floor
	buttonAction := make(chan messages.ButtonEvent_message, 10) //pressed button
	order := make(chan int, 2)                                 // goTo floor
	controllerReady := make(chan bool)                      // Check when controller is done
	elevatorStatusTx := make(chan messages.ElevatorStatus)
	elevatorStatusRx := make(chan messages.ElevatorStatus)
	lampNotifications := make(chan messages.LampUpdate_message, 5)

	// Bundle controller channels in a struct
	ctrChans := elevator.ControllerChannels{
		DoorOpen             : doorOpen,
		CurrentFloor         : currentFloor,
		RedirectButtonAction : buttonAction,
		ReceiveOrder         : order,
		ElevatorUpdateRequest: getElevatorUpdate,
		ControllerReady      : controllerReady,
		DoorObstructed       : doorObstructed,
		ReceiveLampUpdate    : lampNotifications,
	}

	orderChans := orders.OrderChannels{
		Button_press            : buttonAction,
		Received_elevator_update: elevatorStatusRx,
		New_floor               : currentFloor,
		Door_status             : doorOpen,
		Send_status             : elevatorStatusTx,
		Go_to_floor             : order,
		Get_init_position  	 	: getElevatorUpdate,
		DoorObstructed          : doorObstructed,
		UpdateLampMessage       : lampNotifications,
	}

	controller := elevator.NewController(ctrChans)
	go controller.Run()
	<-controllerReady

	go bcast.Transmitter(config.BcastPort, elevatorStatusTx)
	go bcast.Receiver(config.BcastPort, elevatorStatusRx)

	orders.RunOrders(orderChans)
}
