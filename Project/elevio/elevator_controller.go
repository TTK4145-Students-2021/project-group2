package elevio

import (
	"fmt"
	"time"

	"../config"
)

/*
The Elevator Controller works like an extension to an elevator_fsm. It contains channels
needed for communicating between the elevator and the order-module, and methods for sending
operations to the elevator_fsm. This way, we conceptually abstract the physical state machine of
the elvator (the fsm) and the unit trying to perform operations on it.
*/

const _ctrPollRate = 20 * time.Millisecond

type ControllerChannels struct {
	elev_available         chan<- bool
	current_floor          chan<- int
	redirect_button_action chan<- ButtonEvent
	receive_order          <-chan int //Some defined ExecOrder message type to be implemented!!
	respond_order          chan<- error
	// TODO: How do we tell the controller to turn off lights and such??
}

type Controller struct {
	// All the channels the controller needs for communication
	Channels ControllerChannels

	// Has to contain an elevator machine to control
	Elev *ElevatorMachine
}

func NewController(ctrChans ControllerChannels) *Controller {
	Init("localhost:"+config.SimPort, config.NumFloors)
	// TODO: Change name of 'Init'. Bad name
	elev := NewElevatorMachine()
	return &Controller{
		Channels: ctrChans,
		Elev:     elev,
	}
}

func (ctr *Controller) Run() error {
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

	// Send information to recipients wether elevator is available
	go ctr.PollElevAvailability()

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			channels.redirect_button_action <- a // Send to entity that handles this

		case a := <-channels.receive_order:
			fmt.Printf("New order received\n")

			/*NOTE: When we want to tell our order module that the elevator is unavailable,
			there are different possible approaches. We can
			1. Continuously update the order module on a channel
			2. Let the order module try to send and order and then give it as a response
			3. Let the order module check it directly itself (GetAvailable or something like that)
			The implementation below is of case 2, but can be changed!
			*/

			if !elev.Available {
				channels.respond_order <- ErrElevUnavailable // Tell the order system error is unavailable
				break
			}

			//err := elev.GotoFloor(targetFloor)
			err := elev.GotoFloor(a)
			if err != nil {
				// TODO: Tell the order module that the order is rejected
				channels.respond_order <- ErrExecOrderFailed
			} else {
				channels.respond_order <- nil
			}

			// TODO: Do we need an option to shut down this loop? Probably
		}
	}
	return nil
}

// SendElevatorStatus is used to send a full status report to channel recipients if requested
func (ctr *Controller) SendElevatorStatus() {
	ctr.Elev.mutex.Lock()
	defer ctr.Elev.mutex.Unlock()
	ctr.Channels.elev_available <- ctr.Elev.Available
	ctr.Channels.current_floor <- ctr.Elev.CurrentFloor
}

func (ctr *Controller) GetElevator() *ElevatorMachine {
	return ctr.Elev
}

func (ctr *Controller) PollElevAvailability() {
	prev := false
	ctr.Channels.elev_available <- prev
	for {
		time.Sleep(_ctrPollRate)
		ctr.Elev.mutex.Lock()
		v := ctr.Elev.Available
		ctr.Elev.mutex.Unlock()

		if v != prev {
			ctr.Channels.elev_available <- v
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

	if elev.Available == false {
		return ErrElevUnavailable
	}

	moveCtx := &MoveContext{
		TargetFloor: targetFloor,
	}

	err := elev.SendEvent(Move, moveCtx)
	if err != nil {
		// Return any other errors
		return err
	}
	return nil
}
