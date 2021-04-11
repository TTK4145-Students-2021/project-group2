package network

import (
	"fmt"
	"testing"
	"time"
	"strconv"

	"../config"
	"./peers"
)

type OrderChans struct {
	ElevatorStatusTx chan ElevatorStatus // Change to orders.ElevatorStatus
	ElevatorStatusRx chan ElevatorStatus // Change to orders.ElevatorStatus
}

func RunOrders(OrdersCh OrderChans) {
	msg := ElevatorStatus{"Hello from " + strconv.Itoa(config.ID), 0}

	for {
		select {
		case o := <-OrdersCh.ElevatorStatusRx:
			fmt.Printf("Elevator stats received: \n\t%+v\n", o)
		case <-time.After(config.BcastIntervall):
			OrdersCh.ElevatorStatusTx <- msg
			msg.Iter++
		}
	}
}

func TestNetwork(t *testing.T) {
	fmt.Println("[+] Started testing")

	peerUpdateCh := make(chan peers.PeerUpdate)
	elevatorStatusTx := make(chan ElevatorStatus) // Change to orders.ElevatorStatus
	elevatorStatusRx := make(chan ElevatorStatus) // Change to orders.ElevatorStatus

	orderCh := OrderChans{
		ElevatorStatusTx: elevatorStatusTx,
		ElevatorStatusRx: elevatorStatusRx,
	}

	networkCh := NetworkChans{
		PeerUpdate:       peerUpdateCh,
		ElevatorStatusTx: elevatorStatusTx,
		ElevatorStatusRx: elevatorStatusRx,
	}

	go RunOrders(orderCh)
	RunNetwork(networkCh)
}
