package network

import (
	"fmt"
	"testing"
	"time"

	"../config"
	"../messages"

	// "./peers"
	"./bcast"
)

type OrderChans struct {
	// ElevatorStatusTx chan messages.ElevatorStatus // Change to orders.ElevatorStatus
	// ElevatorStatusRx chan messages.ElevatorStatus // Change to orders.ElevatorStatus

	ElevatorStatusTx chan messages.Outer // Change to orders.ElevatorStatus
	ElevatorStatusRx chan messages.Outer // Change to orders.ElevatorStatus
}

func RunOrders(OrdersCh OrderChans) {

	orderList := [config.NumFloors*2 - 2]messages.HallOrder{}
	initOrderList(orderList)

	elevStatus := messages.Outer{
		ID:  0,
		Msg: orderList,
		Pos: 1,
		Dir: 2,
		//OrderList:   orderList,
		IsOnline:    false, //NB isOline is deafult false
		DoorOpen:    false,
		CabOrders:   [config.NumFloors]bool{false},
		IsAvailable: false, //NB isAvaliable is deafult false
	}
	//[config.NumFloors*2 - 2]messages.HallOrder{hall},

	/*
		elevStatus := messages.ElevatorStatus{
			ID:          0,
			Pos:         1,
			Dir:         2,
			OrderList:   orders,
			IsOnline:    false, //NB isOline is deafult false
			DoorOpen:    false,
			CabOrders:   [config.NumFloors]bool{false},
			IsAvailable: false, //NB isAvaliable is deafult false
		}
	*/

	for {
		select {
		case o := <-OrdersCh.ElevatorStatusRx:
			fmt.Printf("Elevator stats received: \n\t%+v\n", o)
		case <-time.After(config.BcastIntervall):
			OrdersCh.ElevatorStatusTx <- elevStatus
			elevStatus.ID++
		}
	}
}

func initOrderList(OrderList [config.NumFloors*2 - 2]messages.HallOrder) {

	// OrderList := [config.NumFloors*2 - 2]messages.HallOrder{}

	// initalizing HallUp

	for i := 0; i < (config.NumFloors - 1); i++ {
		OrderList[i] = messages.HallOrder{
			HasOrder:   false,
			Floor:      i,
			Direction:  0, //up = 0, down = 1
			VersionNum: 0,
			Costs:      [config.NumElevators]int{},
			TimeStamp:  time.Now(),
		}
	}

	// initalizing HallDown
	for i := config.NumFloors - 1; i < (config.NumFloors*2 - 2); i++ {
		OrderList[i] = messages.HallOrder{
			HasOrder:   false,
			Floor:      i + 2 - config.NumFloors,
			Direction:  1, //up = 0, down = 1
			VersionNum: 0,
			Costs:      [config.NumElevators]int{},
			TimeStamp:  time.Now(),
		}
	}
}

func TestNetwork(t *testing.T) {
	fmt.Println("[+] Started testing")

	// peerUpdateCh := make(chan peers.PeerUpdate)
	// elevatorStatusTx := make(chan messages.ElevatorStatus) // Change to orders.ElevatorStatus
	// elevatorStatusRx := make(chan messages.ElevatorStatus) // Change to orders.ElevatorStatus

	elevatorStatusTx := make(chan messages.Outer) // Change to orders.ElevatorStatus
	elevatorStatusRx := make(chan messages.Outer) // Change to orders.ElevatorStatus

	orderCh := OrderChans{
		ElevatorStatusTx: elevatorStatusTx,
		ElevatorStatusRx: elevatorStatusRx,
	}

	// networkCh := NetworkChans{
	// 	PeerUpdate:       peerUpdateCh,
	// 	ElevatorStatusTx: elevatorStatusTx,
	// 	ElevatorStatusRx: elevatorStatusRx,
	// }

	go RunOrders(orderCh)
	// RunNetwork(networkCh)
	go bcast.Transmitter(16569, elevatorStatusTx)
	bcast.Receiver(16569, elevatorStatusRx)

}
