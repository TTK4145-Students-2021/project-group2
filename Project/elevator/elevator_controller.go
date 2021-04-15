package elevator

import (
	"fmt"
	"time"

	"../config"
	"../messages"
)

/*=============================================================================
 * @Description:
 * The Elevator Controller acts as an extension to an ElevatorMachine, working as
 * a broker between the input buttons, the order module and the physical machine.
 * It sets up and connects all the necessary communication channels in the way that
 * is needed for the modules to work together, redirecting inputs and outputs to their
 * intended destinations.
 *
 * @Author: Eldar Sandanger
/*=============================================================================
*/

const _ctrPollRate = 20 * time.Millisecond

type ControllerChannels struct {
	DoorOpen              chan<- bool
	CurrentFloor          chan<- int
	RedirectButtonAction  chan<- messages.ButtonEvent_message
	ReceiveOrder          <-chan int
	RespondOrder          chan<- error
	ElevatorUpdateRequest <-chan bool
	ControllerReady       chan<- bool
	DoorObstructed        chan<- bool
}

type Controller struct {
	// All the channels the controller needs for communication
	Channels ControllerChannels

	// Has to contain an elevator machine to control
	Elev *ElevatorMachine
}

// NewController returns a controller containing an initialized ElevatorMachine
func NewController(ctrChans ControllerChannels) *Controller {
	InitElevatorDriver("localhost:"+config.SimPort, config.NumFloors)
	elev := NewElevatorMachine()
	return &Controller{
		Channels: ctrChans,
		Elev:     elev,
	}
}

// Run is used to 'turn on' and run a controller
func (ctr *Controller) Run() error {
	fmt.Println("Controller running")
	channels := ctr.Channels
	elev := ctr.Elev

	drv_buttons := make(chan ButtonEvent)
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

	SetFloorIndicator(elev.CurrentFloor)

	// Start elevator driver routines
	go PollButtons(drv_buttons)
	go PollFloorSensor(drv_floors)
	go PollObstructionSwitch(drv_obstr)

	// Start polling routines for updating external listener
	go ctr.pollDoorOpen(elev)
	go ctr.pollElevFloor(elev)
	go ctr.pollDoorObstructed(elev)

	ctr.Channels.ControllerReady <- true

	fmt.Println("Controller waiting for action...")
	for {
		select {
		case a := <-drv_buttons:
			fmt.Println("Button press registered")
			fmt.Printf("%+v\n", a)

			// Convert to type ButtonEvent_message before sending
			btm := messages.UNDEFINED
			switch a.Button {
			case BT_HallUp:
				btm = messages.BT_HallUp
			case BT_HallDown:
				btm = messages.BT_HallDown
			case BT_Cab:
				btm = messages.BT_Cab
			}
			message := messages.ButtonEvent_message{
				Floor:  a.Floor,
				Button: btm,
			}

			channels.RedirectButtonAction <- message

		case a := <-channels.ReceiveOrder:
			fmt.Printf("New order received\n")
			fmt.Printf("Target floor: %v\n", a)

			go elev.GotoFloor(a)

		case <-channels.ElevatorUpdateRequest:
			channels.CurrentFloor <- elev.CurrentFloor
			channels.DoorOpen <- elev.Available
			channels.DoorObstructed <- GetObstruction()
		}
	}
}

// SendElevatorStatus is used to send a full status report to channel recipients if requested
func (ctr *Controller) SendElevatorStatus() {
	ctr.Elev.mutex.Lock()
	defer ctr.Elev.mutex.Unlock()
	ctr.Channels.DoorOpen <- ctr.Elev.Available
	ctr.Channels.CurrentFloor <- ctr.Elev.CurrentFloor
}

// pollDoorOpen sends update to listener if door opens or closes
func (ctr *Controller) pollDoorOpen(elev *ElevatorMachine) {
	prev := elev.Available
	ctr.Channels.DoorOpen <- prev
	for {
		time.Sleep(_ctrPollRate)
		v := ctr.Elev.Available

		if v != prev {
			ctr.Channels.DoorOpen <- v
		}
		prev = v
	}
}

// pollElevFloor sends update to listener if elevator enters new floor
func (ctr *Controller) pollElevFloor(elev *ElevatorMachine) {
	prev := elev.CurrentFloor
	ctr.Channels.CurrentFloor <- prev
	for {
		time.Sleep(_ctrPollRate)
		elev.mutex.Lock()
		v := ctr.Elev.CurrentFloor
		elev.mutex.Unlock()

		if v != prev && v != -1 {
			SetFloorIndicator(v)
			ctr.Channels.CurrentFloor <- v
			prev = v
		}
	}
}

// pollDoorObstructed sends update to listener if elevator door is obstructed
func (ctr *Controller) pollDoorObstructed(elev *ElevatorMachine) {
	prev := GetObstruction()
	ctr.Channels.DoorObstructed <- prev
	for {
		time.Sleep(_ctrPollRate)
		v := GetObstruction()
		if v != prev {
			ctr.Channels.DoorObstructed <- v
			prev = v
		}
	}
}

/*________________________Wrapper functions___________________________*
*
* These functions are used to simplify the elevator interface from the controllers perspective.
* To perform an operation on the elevator, you have to send an Event that takes inn the Event name
* and an Event context. The context has to be created with the appropriate content first, and these
* functions perform this operation.
 */

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
