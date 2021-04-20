package elevController

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

const _doorPollRate  = 20 * time.Millisecond
const _floorPollRate = 20 * time.Millisecond
const _obstrPollRate = 20 * time.Millisecond

type ControllerChannels struct {
	DoorOpen              chan<- bool
	CurrentFloor          chan<- int
	RedirectButtonAction  chan<- messages.ButtonEvent
	ReceiveOrder          <-chan int
	RespondOrder          chan<- error
	ElevatorUpdateRequest <-chan bool
	ControllerReady       chan<- bool
	DoorObstructed        chan<- bool
	ReceiveLampUpdate     <-chan messages.LampUpdate
}

type Controller struct {
	Channels ControllerChannels
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

	drvButtons := make(chan ButtonEvent)
	drvFloors := make(chan int)
	drvObstr := make(chan bool)

	elevChans := ElevChannels{
		drv_floors: drvFloors,
		drv_obstr:  drvObstr,
	}
	err := elev.Initialize(elevChans)
	if err != nil {
		return err
	}

	SetFloorIndicator(elev.CurrentFloor)

	// Start elevator driver routines
	go PollButtons(drvButtons)
	go PollFloorSensor(drvFloors)

	// Start routines for communication with external modules
	go ctr.pollDoorOpen(elev)
	go ctr.pollElevFloor(elev)
	go ctr.pollDoorObstructed(elev)

	ctr.Channels.ControllerReady <- true

	fmt.Println("Controller waiting for action...")
	for {
		select {
		case a := <-drvButtons:

			// Redirect message to order module
			message := messages.ButtonEvent{
				Floor:  a.Floor,
				Button: messages.ButtonType(a.Button),
			}
			channels.RedirectButtonAction <- message

		case a := <-channels.ReceiveOrder:
			go elev.GotoFloor(a)

		case a := <-channels.ReceiveLampUpdate:
			go handleLampUpdate(a)

		case <-channels.ElevatorUpdateRequest:
			channels.CurrentFloor <- elev.CurrentFloor
			channels.DoorOpen <- elev.Available
			channels.DoorObstructed <- GetObstruction()
		}
	}
}

func handleLampUpdate(message messages.LampUpdate) {
	SetButtonLamp(ButtonType(message.Button), message.Floor, message.Turn)
}

// SendElevatorStatus is used to send a full status report to channel recipients if requested
func (ctr *Controller) SendElevatorStatus() {
	ctr.Elev.mutex.Lock()
	defer ctr.Elev.mutex.Unlock()
	ctr.Channels.DoorOpen <- ctr.Elev.Available
	ctr.Channels.CurrentFloor <- ctr.Elev.CurrentFloor
}

// Updates listener if door opens or closes
func (ctr *Controller) pollDoorOpen(elev *ElevatorMachine) {
	prev := elev.Available
	ctr.Channels.DoorOpen <- prev
	for {
		time.Sleep(_doorPollRate)
		v := ctr.Elev.Available

		if v != prev {
			ctr.Channels.DoorOpen <- v
		}
		prev = v
	}
}

// Updates listener when elevator enters new floor
func (ctr *Controller) pollElevFloor(elev *ElevatorMachine) {
	prev := elev.CurrentFloor
	ctr.Channels.CurrentFloor <- prev
	for {
		time.Sleep(_floorPollRate)
		v := GetFloor()
		if v != prev && v != -1 {
			SetFloorIndicator(v)
			ctr.Channels.CurrentFloor <- v
			prev = v
		}
	}
}

// Updates listener if elevator door is obstructed
func (ctr *Controller) pollDoorObstructed(elev *ElevatorMachine) {
	prev := GetObstruction()
	ctr.Channels.DoorObstructed <- prev
	for {
		time.Sleep(_obstrPollRate)
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
