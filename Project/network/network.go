package network

import (
	"fmt"
	"strconv"

	"./bcast"
	"./peers"
)

/*
Channels overview
peerUpdateCh := make(chan peers.PeerUpdate)				Receiving updates on the id's of the peers that are alive on the network
peerTxEnable := make(chan bool)							We can disable/enable the transmitter after it has been started. This could be used to signal that we are somehow "unavailable".
ElevatorStatusTx := make(chan orders.ElevatorStatus)	Transmit the ElevatorStatus on the network
ElevatorStatusRx := make(chan orders.ElevatorStatus)	Receive ElevatorStatus on the network
x - GetElevatorStatus := make(chan orders.ElevatorStatus)	Get ElevatorStatus from the order-module
x - SendEle
*/

type NetworkChans struct {
	PeerUpdate       chan peers.PeerUpdate
	ElevatorStatusTx chan Msg // Change to orders.ElevatorStatus
	ElevatorStatusRx chan Msg // Change to orders.ElevatorStatus
}

type OrderChans struct {
	ElevatorStatusTx chan Msg // Change to orders.ElevatorStatus
	ElevatorStatusRx chan Msg // Change to orders.ElevatorStatus
}

type Msg struct {
	ID          int
	pos         int
	isAvailable bool
	cabOrders   []bool
}

func RunNetwork(PeersPort int, BcastPort int, id int, networkCh NetworkChans) {
	id_str := strconv.Itoa(id)
	go peers.Transmitter(PeersPort, id_str)
	go peers.Receiver(PeersPort, networkCh.PeerUpdate)
	go bcast.Transmitter(BcastPort, networkCh.ElevatorStatusTx)
	go bcast.Receiver(BcastPort, networkCh.ElevatorStatusRx)

	for {
		select {
		case p := <-networkCh.PeerUpdate:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-networkCh.ElevatorStatusRx:
			fmt.Printf("Received: %#v\n", a)
		}
	}
}

func SendElevatorStatus() {

}

func RecvElevatorStatus() {

}
