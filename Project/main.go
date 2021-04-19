// run main with as follows:
// go run main.go [SimPort] [elevatorID]

package main

import (
	"fmt"

	"./config"
	"./messages"
	"./network/bcast"
	"./orderDistribution"
	"./elevController"
)

func main() {

	config.InitConfig()
	fmt.Println(config.NumFloors, config.NumElevators, config.Port, config.ID)

	// Setup all the channels we need
	ElevatorUpdateRequest := make(chan bool)
	DoorOpen := make(chan bool, 1)
	DoorObstructed := make(chan bool,1)
	CurrentFloor := make(chan int, 1)                         
	ButtonAction := make(chan messages.ButtonEvent, 1) 
	ExecutionOrder := make(chan int, 1)                           
	ControllerReady := make(chan bool)                      
	ElevatorStatusTx := make(chan orderDistribution.ElevatorStatus)
	ElevatorStatusRx := make(chan orderDistribution.ElevatorStatus)
	LampUpdates := make(chan messages.LampUpdate, 1)

	// Bundle controller channels in a struct
	ctrChannels := elevController.ControllerChannels{
		DoorOpen             : DoorOpen,
		CurrentFloor         : CurrentFloor,
		RedirectButtonAction : ButtonAction,
		ReceiveOrder         : ExecutionOrder,
		ElevatorUpdateRequest: ElevatorUpdateRequest,
		ControllerReady      : ControllerReady,
		DoorObstructed       : DoorObstructed,
		ReceiveLampUpdate    : LampUpdates,
	}

	orderChannels := orderDistribution.OrderChannels{
		ButtonPress          : ButtonAction,
		ReceiveElevatorUpdate: ElevatorStatusRx,
		NewFloor             : CurrentFloor,
		DoorOpen             : DoorOpen,
		SendStatus		     : ElevatorStatusTx,
		GotoFloor            : ExecutionOrder,
		RequestElevatorStatus: ElevatorUpdateRequest,
		DoorObstructed       : DoorObstructed,
		UpdateLampMessage    : LampUpdates,
	}

	controller := elevController.NewController(ctrChannels)
	go controller.Run()
	<-ControllerReady

	go bcast.Transmitter(config.BcastPort, ElevatorStatusTx)
	go bcast.Receiver(config.BcastPort, ElevatorStatusRx)

	orderDistribution.RunOrders(orderChannels)
}
