package elevio

import (
	"errors"
	"log"
	"sync"
	"time"

	"../config"
	"../messages"
)

// ElevatorStatus keeps a local track of the status of each elevator
type ElevatorStatus struct {
	ID           int
	Available    bool
	CurrentFloor int
	AtTop        bool
	AtBottom     bool
	mutex        sync.Mutex
}

// TODO: Have to implement a mapping of ID to ElevatorStatus' to track all elevators

// OrderChannels contain all the channels needed for the order module and are passed by main
type OrderChannels struct {
	elev_available <-chan bool
	current_floor  <-chan int
	button_event   <-chan messages.ButtonType_message
	send_order     chan<- int
	order_response <-chan error
}

// ExecOrder contains an order request from the order module to the elevator controller
type ExecOrder struct {
	targetFloor int
}

// Order contains an order from the button panel
type Order struct {
	ButtonType ButtonType
	Floor      int
}

// Orderlist contains and can process all orders regardless of type
type OrderList struct {
	hallOrders []Order
	cabOrders  []Order
	mutex      sync.Mutex
}

/*Note on OrderList: The reason for splitting into these two sublists is initially that
cab orders always should be performed before hall orders. Then we can put a priority on
the cabOrders list and only execute hallOrders if cabOrders == nil
*/

/*Note on Order Simulator: In this case I can just use a function, as the main purpose for this
is testing. However, you might find that there are some problems with trying to make
the order module using just a function! I think it's a better approach to create a
struct.

Why? Let me give an example: I now have only one elevator passed into the function.
What if I want to add a new elevator? I can't just redefine how many elevators that
should be passed to our module every time! It's better to store a map of all elevators
with their ID, and then create a method that can 'add' each elevator with their proper ID's.
This is just one example of a difficulty that may arise
*/

// OrderSimulator simulates simple order management for testing the controller
func OrderSimulator(chans *OrderChannels, status *ElevatorStatus) error {
	log.Println("Order simulator instantiated")

	orderList := newOrderList()

	/*This routine updates the elevators availability and floor status*/
	go elevStatusListener(chans, status)

	/*This routine should continuously process orders in the list and send
	them to the elevator controller. Should ideally not be a simple function, but
	a state machine!*/
	go manageOrders(orderList, chans, status)

	// TODO: Make all goroutines properly stoppable

	for {
		/*
			select {

			case a := <-chans.button_event:
				log.Println("Button event received")

				switch a.Button {
				case BT_HallDown:
					err := orderList.addHallOrder(a)
					if err != nil {
						fmt.Println("Hall order already exists")
					}
				case BT_HallUp:
					err := orderList.addHallOrder(a)
					if err != nil {
						fmt.Println("Hall order already exists")
					}
				case BT_Cab:
					err := orderList.addCabOrder(a)
					if err != nil {
						fmt.Println("Cab order already exists")
					}
				}
			}
		*/
	}
}

/* manageOrders should continuously check the status of the orderlist at a constant
rate and modify it accordingly if there are any priority changes that have to be made.
If there are any new orders, they will be added to this list and the processer finds it's
calculates it's place in the que.

It's final purpose is to send actual executional orders to the elevator!

This should also probably be it's own struct with methods. Just using a function is ok
i guess, but as complexity increases, so will the difficulty of implementing it drastically.
*/

func manageOrders(ol *OrderList, chans *OrderChannels, status *ElevatorStatus) {

	response := chans.order_response

	for status.Available {
		if len(ol.cabOrders) != 0 {

			order := ol.cabOrders[0]
			chans.send_order <- order.Floor
			res := <-response
			if res != nil {
				log.Println("[OrderManager] Execution request rejected")
			}
			ol.cabOrders = ol.cabOrders[1:] // Terrible solution :)
			// TODO: Create a queue data type with build in intuitive functionality

		} else {
			if len(ol.hallOrders) != 0 {

				order := ol.hallOrders[0]
				chans.send_order <- order.Floor
				res := <-response
				if res != nil {
					log.Println("[OrderManager] Execution request rejected")
				}
				ol.hallOrders = ol.hallOrders[1:] // Terrible solution here too :)
				// TODO: Create a queue data type with build in intuitive functionality
			}
		}
		// Just allow for some time to avoid errors
		time.Sleep(time.Millisecond * 20)
	}
	/*Tasks to perform:
	1. Send an order to the elevator when available.
	2. Calculate which task is next to send regularly.

	An immediate conclusion we can draw is that we don't want to try to send execution
	request all the time and spam the controller. We also don't want to sort tasks all the
	time if there are no reasons to do so. We should treat this function carefully, as
	a lot of 'traffic' will go through here (cost function etc.). Therefore, let's be careful
	about how many operations we perform simultaneously.*/

	return
}

func elevStatusListener(chans *OrderChannels, status *ElevatorStatus) {

	for {
		select {
		case a := <-chans.elev_available:
			log.Println("Elevator availability status changed")
			status.mutex.Lock()
			defer status.mutex.Unlock()
			status.Available = a

		case a := <-chans.current_floor:
			log.Printf("CurrentFloor %v", a)
			status.mutex.Lock()
			defer status.mutex.Unlock()
			status.CurrentFloor = a

			// Check if top or bottom
			switch a {
			case 0:
				status.AtBottom = true
				status.AtTop = false
			case config.NumFloors - 1:
				status.AtBottom = false
				status.AtTop = true
			default:
				status.AtBottom = false
				status.AtTop = false
			}
		}
	}
}

func newOrderList() *OrderList {
	return &OrderList{
		hallOrders: make([]Order, 0),
		cabOrders:  make([]Order, 0),
	}
}

// Add HallOrder to Order List
func (ol *OrderList) addHallOrder(btnEvent ButtonEvent) error {
	ol.mutex.Lock()
	defer ol.mutex.Unlock()
	log.Println("Request to add hall order received")

	// Check if similar order already there. If true, return order-exists error
	for _, order := range ol.hallOrders {
		if order.ButtonType == btnEvent.Button && order.Floor == btnEvent.Floor {
			return ErrOrderExists
		}
	}
	// Create and add new order
	newOrder := Order{
		ButtonType: btnEvent.Button,
		Floor:      btnEvent.Floor,
	}
	ol.hallOrders = append(ol.hallOrders, newOrder)
	return nil
}

// Add CabOrder to OrderList
func (ol *OrderList) addCabOrder(btnEvent ButtonEvent) error {
	ol.mutex.Lock()
	defer ol.mutex.Unlock()
	log.Println("Request to add cab order received")

	for _, order := range ol.cabOrders {
		if order.Floor == btnEvent.Floor {
			return ErrOrderExists
		}
	}
	// Create and add new order
	newOrder := Order{
		ButtonType: BT_Cab,
		Floor:      btnEvent.Floor,
	}
	ol.cabOrders = append(ol.cabOrders, newOrder)
	return nil
}

/***************ERROR MESSAGES*****************/

// ErrOrderExists returned if order already in list
var ErrOrderExists = errors.New("order already exists")
