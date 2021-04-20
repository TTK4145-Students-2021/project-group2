package main

import (
	"fmt"

	"./config"
	"./elevController"
	"./messages"
	"./network/bcast"
	"./orderDistribution"
)

func main() {

	config.InitConfig()
	fmt.Println(config.NumFloors, config.NumElevators, config.Port, config.ID)

	// Setup all the channels we need
	getElevatorUpdate := make(chan bool)
	doorOpen := make(chan bool, 1)
	doorObstructed := make(chan bool, 1)
	currentFloor := make(chan int, 1)
	buttonAction := make(chan messages.ButtonEvent, 1)
	order := make(chan int, 1)
	controllerReady := make(chan bool)
	elevatorStatusTx := make(chan orderDistribution.ElevatorStatus)
	elevatorStatusRx := make(chan orderDistribution.ElevatorStatus)
	lampNotifications := make(chan messages.LampUpdate, 1)

	// Bundle controller channels in a struct
	ctrChans := elevController.ControllerChannels{
		DoorOpen:              doorOpen,
		CurrentFloor:          currentFloor,
		RedirectButtonAction:  buttonAction,
		ReceiveOrder:          order,
		ElevatorUpdateRequest: getElevatorUpdate,
		ControllerReady:       controllerReady,
		DoorObstructed:        doorObstructed,
		ReceiveLampUpdate:     lampNotifications,
	}

	orderChans := orderDistribution.OrderChannels{
		ButtonPress:            buttonAction,
		ReceivedElevatorUpdate: elevatorStatusRx,
		NewFloor:               currentFloor,
		DoorOpen:               doorOpen,
		SendStatus:             elevatorStatusTx,
		GotoFloor:              order,
		GetInitPosition:        getElevatorUpdate,
		DoorObstructed:         doorObstructed,
		UpdateLampMessage:      lampNotifications,
	}

	controller := elevController.NewController(ctrChans)
	go controller.Run()
	<-controllerReady

	go bcast.Transmitter(config.BcastPort, elevatorStatusTx)
	go bcast.Receiver(config.BcastPort, elevatorStatusRx)

	orderDistribution.RunOrders(orderChans)
}
