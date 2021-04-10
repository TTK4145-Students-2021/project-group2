package elevator

import (
	"fmt"
	"time"

	"../config"
	"../messages"
)

/*
The Elevator Controller works like an extension to an elevator_fsm. It contains channels
needed for communicating between the elevator and the order-module, and methods for sending
operations to the elevator_fsm. This way, we conceptually abstract the physical state machine of
the elvator (the fsm) and the unit trying to perform operations on it.
*/

const _ctrPollRate = 20 * time.Millisecond

type ControllerChannels struct {
	Elev_available         chan<- bool
	Current_floor          chan<- int
	Redirect_button_action chan<- messages.ButtonEvent_message
	Receive_order          <-chan int //Some defined ExecOrder message type to be implemented!!
	Respond_order          chan<- error
	ControllerReady        chan<- bool
	// TODO: How do we tell the controller to turn off lights and such??
}

type Controller struct {
	// All the channels the controller needs for communication
	Channels ControllerChannels

	// Has to contain an elevator machine to control
	Elev *ElevatorMachine
}

func NewController(ctrChans ControllerChannels) *Controller {
	InitElevatorDriver("localhost:"+config.SimPort, config.NumFloors)
	elev := NewElevatorMachine()
	return &Controller{
		Channels: ctrChans,
		Elev:     elev,
	}
}

func (ctr *Controller) Run() error {
	fmt.Println("Controller running")
	channels := ctr.Channels
	elev := ctr.Elev

	// Setup "Button panel" interface
	drv_buttons := make(chan ButtonEvent)

	// Setup and run elevator machine instance
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	elevChans := ElevChannels{
		drv_floors: drv_floors,
		drv_obstr:  drv_obstr,
	}
	err := elev.Initialize(elevChans)
	if err != nil {
		return err
	}
	// Start elevio routines
	go PollButtons(drv_buttons)
	go PollFloorSensor(drv_floors)
	go PollObstructionSwitch(drv_obstr)

	// Send information to recipients whether elevator is available
	go ctr.PollElevAvailability(elev)

	ctr.Channels.ControllerReady <- true

	fmt.Println("Controller waiting for action...")
	for {
		select {
		case a := <-drv_buttons:
			fmt.Println("Button press registered")
			fmt.Printf("%+v\n", a)

			switch a.Button {
			case BT_HallUp:
				bt := messages.BT_HallUp
				message := messages.ButtonEvent_message{
					Floor:  a.Floor,
					Button: bt,
				}
				channels.Redirect_button_action <- message
			case BT_HallDown:
				bt := messages.BT_HallDown
				message := messages.ButtonEvent_message{
					Floor:  a.Floor,
					Button: bt,
				}
				channels.Redirect_button_action <- message
			case BT_Cab:
				bt := messages.BT_Cab
				message := messages.ButtonEvent_message{
					Floor:  a.Floor,
					Button: bt,
				}
				channels.Redirect_button_action <- message
			}

		case a := <-channels.Receive_order:
			fmt.Printf("New order received\n")
			fmt.Printf("Value %v\n", a)

			/*NOTE: When we want to tell our order module that the elevator is unavailable,
			there are different possible approaches. We can
			1. Continuously update the order module on a channel
			2. Let the order module try to send and order and then give it as a response
			3. Let the order module check it directly itself (GetAvailable or something like that)
			The implementation below is of case 2, but can be changed!
			*/

			// MIGHT NOT BE NECESSARY IF THE ORDER MODULE IS CONTINUOUSLY UPDATED
		//	if !elev.Available {
		//		channels.Respond_order <- ErrElevUnavailable // Tell the order system error is unavailable
		//		break
		//	}

			go elev.|Floor(a)
			/*
			//err := elev.GotoFloor(targetFloor)
			err := elev.GotoFloor(a)
			if err != nil {
				// TODO: Tell the order module that the order is rejected
				channels.Respond_order <- ErrExecOrderFailed
			} else {
				channels.Respond_order <- nil
			}
			*/

			// TODO: Do we need an option to shut down this loop? Probably
		}
	}
}

// SendElevatorStatus is used to send a full status report to channel recipients if requested
func (ctr *Controller) SendElevatorStatus() {
	ctr.Elev.mutex.Lock()
	defer ctr.Elev.mutex.Unlock()
	ctr.Channels.Elev_available <- ctr.Elev.Available
	ctr.Channels.Current_floor <- ctr.Elev.CurrentFloor
}

func (ctr *Controller) GetElevator() *ElevatorMachine {
	return ctr.Elev
}

func (ctr *Controller) PollElevAvailability(elev *ElevatorMachine) {
	prev := elev.Available
	ctr.Channels.Elev_available <- prev
	for {
		time.Sleep(_ctrPollRate)
		ctr.Elev.mutex.Lock()
		v := ctr.Elev.Available
		ctr.Elev.mutex.Unlock()

		if v != prev {
			ctr.Channels.Elev_available <- v
		}
		prev = v
	}
}



/*________________________Wrapper functions___________________________*/
/*These are used to simplify the elevator interface from the controllers perspective*/

// Initialize an elevator
func (elev *ElevatorMachine) Initialize(elevChans ElevChannels) error {
	initCtx := &InitContext{
		Channels: elevChans,
	}
	err := elev.SendEvent(Initialize, initCtx)
	if err != nil {
		return err
	}
	return nil
}

// GotoFloor tells an elevator to mote to targetFloor
func (elev *ElevatorMachine) GotoFloor(targetFloor int) error {

	fmt.Println("GotoFloor called")
	if elev.Available == false {
		fmt.Println("GotoFloor failed")
		return ErrElevUnavailable
	}

	moveCtx := &MoveContext{
		TargetFloor: targetFloor,
	}

	fmt.Println("Sending order event to elevator\n")
	err := elev.SendEvent(Move, moveCtx)
	if err != nil {
		fmt.Println("Failed to send order to elevator")
		// Return any other errors
		return err
	}
	return nil
}
