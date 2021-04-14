package network

import (
	"fmt"
	"strconv"
	// "time"
	"./bcast"
	"./peers"
	"../config"
)

/*
Channels overview
peerUpdateCh := make(chan peers.PeerUpdate)				Receiving updates on the id's of the peers that are alive on the network
ElevatorStatusTx := make(chan orders.ElevatorStatus)	Transmit the ElevatorStatus on the network
ElevatorStatusRx := make(chan orders.ElevatorStatus)	Receive ElevatorStatus on the network
x - peerTxEnable := make(chan bool)						We can disable/enable the transmitter after it has been started. This could be used to signal that we are somehow "unavailable".
*/

type ElevatorStatus struct {
	Message string
	Iter    int
}

type NetworkChans struct {
	PeerUpdate       chan peers.PeerUpdate
	ElevatorStatusTx chan ElevatorStatus // Change to orders.ElevatorStatus
	ElevatorStatusRx chan ElevatorStatus // Change to orders.ElevatorStatus
}


func RunNetwork(networkCh NetworkChans) {

	go peers.Transmitter(15647, strconv.Itoa(config.ID))
	go peers.Receiver(15647, networkCh.PeerUpdate)
	go bcast.Transmitter(16569, networkCh.ElevatorStatusTx)
	fmt.Println("Network started")
	bcast.Receiver(16569, networkCh.ElevatorStatusRx)

}
