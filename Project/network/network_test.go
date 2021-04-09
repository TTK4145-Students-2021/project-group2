package network

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"../config"
	"./peers"
)

const interval = 1000 * time.Millisecond

/*
type order struct {
	hasOrder   bool
	floor      int
	direction  int //up = 0, down = 1
	versionNum int
	costs      []int
	timeStamp  time.Time
}

type ElevatorStatus struct {
	ID          int
	pos         int
	orderList   [config.NumFloors*2 - 2]order //
	dir         elevio.MotorDirection
	isOnline    bool
	doorOpen    bool
	cabOrder    [config.NumFloors]bool
	isAvailable bool
}

func initElevatorStatus(startPos int) ElevatorStatus {
	var emptyIntList [config.NumFloors]int
	var emptyBoolList [config.NumFloors]bool

	initOrder := order{false, 0, emptyIntList, time.Now()}
	initOrderList := orderList[config.NumFloors][2]

	for i, _ := range initOrderList {
		for j, _ := range initOrderList[i] {
			initOrderList[i][j] = initOrder
		}
	}

	return ElevatorStatus{
		ID:          config.ID,
		pos:         startPos,
		orderList:   initOrderList,
		dir:         elevio.MD_Stop,
		isOnline:    true,
		doorOpen:    false,
		cabOrder:    emptyBoolList,
		isAvailable: true,
	}
}
*/

func RunOrders(OrdersCh OrderChans) {
	msg := Msg{
		ID:          1,
		pos:         1,
		isAvailable: true,
		cabOrders:   make([]bool, config.NumFloors),
	}
	for {
		select {
		case o := <-OrdersCh.ElevatorStatusRx:
			fmt.Printf("Elevator stats received: \n\t%+v", o)
		case <-time.After(interval):
			OrdersCh.ElevatorStatusTx <- msg
			msg.pos += 1
		}
	}
}

func TestNetwork(t *testing.T) {

	peerUpdateCh := make(chan peers.PeerUpdate)
	elevatorStatusTx := make(chan Msg) // Change to orders.ElevatorStatus
	elevatorStatusRx := make(chan Msg) // Change to orders.ElevatorStatus

	networkCh := NetworkChans{
		PeerUpdate:       peerUpdateCh,
		ElevatorStatusTx: elevatorStatusTx,
		ElevatorStatusRx: elevatorStatusRx,
	}

	orderCh := OrderChans{
		ElevatorStatusTx: elevatorStatusTx,
		ElevatorStatusRx: elevatorStatusRx,
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		RunNetwork(config.PeersPort, config.BcastPort, config.ID, networkCh)
		wg.Done()
	}()
	go func() {
		RunOrders(orderCh)
		wg.Done()
	}()
	wg.Wait()
}

/*
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16569, helloTx)
	go bcast.Receiver(16569, helloRx)

	// The example message. We just send one of these every second.
	go func() {
		helloMsg := HelloMsg{"Hello from " + id, 0}
		for {
			helloMsg.Iter++
			helloTx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-helloRx:
			fmt.Printf("Received: %#v\n", a)
		}
	}
}
*/
