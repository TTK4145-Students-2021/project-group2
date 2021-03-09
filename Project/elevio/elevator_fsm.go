package elevio

import (
	"fmt"
	"sync"
	"time"
)

// language spec: https://golang.org/ref/spec#Function_literals
// Framework: https://venilnoronha.io/a-simple-state-machine-framework-in-go

const (
	// States of the elevator
	Uninitialized      StateType = "Uninitialized"
	Initializing       StateType = "Initializing"
	MovingUp           StateType = "MovingUp"
	MovingDown         StateType = "MovingDown"
	Obstructed         StateType = "Obstructed"
	AtFloorDoorsClosed StateType = "AtFloorDoorsClosed"
	AtFloorDoorsOpen   StateType = "AtFloorDoorsOpen"
	IdleBetweenFloors  StateType = "IdleBetweenFloors"
	// CheckingPosition   StateType = "CheckingPosition" (?) In an attempt to act in idle between floors
	//Stopped StateType = "Stopped" (?)
	//BottomFloor StateType = "BottomFloor" (?)
	//TopFloor StateType = "TopFloor" (?)

	// Events that can be processed
	// NoOp represents a no-op event
	NoOp           EventType = ""
	Initialize     EventType = "Initialize"
	SetInitialized EventType = "SetInitialized"
	OpenDoors      EventType = "OpenDoors"
	CloseDoors     EventType = "CloseDoors"
	MoveUp         EventType = "MoveUp"
	MoveDown       EventType = "MoveDown"
	Stop           EventType = "Stop"
	Obstruct       EventType = "Obstruct"
	ArriveAtFloor  EventType = "ArriveAtFloor"
	// CheckPos   EventType = "CheckPos" (?) In an attemt to check if inbetween floors
)

// StateType is an extensible state type in the elevator
type StateType string

// EventType is an extensible event type in the elevator
type EventType string

// EventContext represents the context to be passed to the action implementation
type EventContext interface{}

// Action represents the action to be executed in a given state
type Action interface {
	// Pass an eventContext and receive an event type
	Execute(elev *ElevatorMachine, eventCtx EventContext) EventType
}

// Events represents the mapping of events that can be performed in given states
type Events map[EventType]StateType

// State binds a state with an action and a set of events it can handle in that given state
type State struct {
	Action Action
	Events Events
}

// States represents a mapping of states and their implementations
type States map[StateType]State

// ElevatorMachine represents the elevator itself
type ElevatorMachine struct {
	// Previous state
	Previous StateType

	// Current state
	Current StateType

	// Holds the configuration of states and events handled by the state machine
	States States

	// Contains all channels the elevator needs from elevio
	Channels ElevChannels

	// Number of floors
	TotalFloors int

	// Top floor number ?
	TopFloor int

	// Bottom floor number ?
	BottomFloor int

	// Track if at top floor ?
	AtTop bool

	// Track if at bottom floor ?
	AtBottom bool

	// DO WE NEED ONE FOR CURRENT FLOOR? Maybe

	// Ensures only one event is processed by the machine at a time
	mutex sync.Mutex
}

/*The first function argument states that we want the function to act on an elevator object
The second is our method, and the third is our return types*/

// Return next state for given event, or an error if the event can't be handled in this state
func (elev *ElevatorMachine) getNextState(event EventType) (StateType, error) {

	// Check if the event is possible in the current state
	if state, ok := elev.States[elev.Current]; ok {
		if state.Events != nil {
			if next, ok := state.Events[event]; ok {
				return next, nil
			}
		}
	}
	return elev.Current, ErrEventRejected //This is most likely wrong!!
}

// SendEvent sends an event to the state machine
func (elev *ElevatorMachine) SendEvent(event EventType, eventCtx EventContext) error {
	elev.mutex.Lock()
	defer elev.mutex.Unlock()

	for {

		// Determine the next state for the event given the elevators current state
		nextState, err := elev.getNextState(event)
		if err != nil {
			return ErrEventRejected
		}
		// Check if possible to move up or down. Might have to do this elsewhere
		switch nextState {
		case MovingUp:
			if elev.AtTop {
				return ErrIllegalMoveUp
			}
		case MovingDown:
			if elev.AtBottom {
				return ErrIllegalMoveDown
			}
		}

		// TODO: Create functionality to check if the next state is equal the current

		// Identify the state definition for the next state
		state, ok := elev.States[nextState]
		if !ok || state.Action == nil {
			//Configuration error
			// If state.Action == nil this implies there are no actions in the next state
			// This means you have to make one!
		}

		// Transition to the next state
		elev.Previous = elev.Current
		elev.Current = nextState

		// Execute the next states action and loop over if the next state presents a
		// new action that needs to be performed immediately
		nextEvent := state.Action.Execute(elev, eventCtx)
		if nextEvent == NoOp {
			return nil
		}
		event = nextEvent
	}
}

// NewElevatorMachine initializes and returns an elevator object to be used
func NewElevatorMachine() *ElevatorMachine {
	return &ElevatorMachine{

		// This maps all the possible states and events appropriately
		States: States{

			Uninitialized: State{
				Events: Events{
					Initialize: Initializing,
				},
			},

			Initializing: State{
				Action: &InitAction{},
				Events: Events{
					Obstruct:       Obstructed,
					SetInitialized: AtFloorDoorsClosed,
				},
			},

			Obstructed: State{
				Action: &OnObstrAction{},
				Events: Events{
					MoveDown: MovingDown,
					MoveUp:   MovingUp,
					// TODO: Implement unobstruction routine
				},
			},

			MovingUp: State{
				Action: &MovingUpAction{},
				Events: Events{
					Obstruct:      Obstructed,
					ArriveAtFloor: AtFloorDoorsClosed,

					// Continue in current state if passed
					MoveUp: MovingUp,
				},
			},

			MovingDown: State{
				Action: &MovingDownAction{},
				Events: Events{
					Obstruct:      Obstructed,
					ArriveAtFloor: AtFloorDoorsClosed,

					// Continue in current state if passed
					MoveDown: MovingDown,
				},
			},

			AtFloorDoorsClosed: State{
				Action: &AtFloorClosedAction{},
				Events: Events{
					OpenDoors: AtFloorDoorsOpen,
					Obstruct:  Obstructed,
					MoveUp:    MovingUp,
					MoveDown:  MovingDown,
				},
			},

			AtFloorDoorsOpen: State{
				Action: &AtFloorOpenAction{},
				Events: Events{
					CloseDoors: AtFloorDoorsClosed,
					Obstruct:   Obstructed,
				},
			},
		},
		Current: Uninitialized,
	}
}

/*NOW WE HAVE TO DEFINE CONTEXT TYPES FOR PASSING DATA TO EVENT ACTIONS*/
// Note that not all event actions need to be passed data.

// InitContext contains information that needs to passed when initializing elevator
type InitContext struct {
	// What do we need for initializing the elevator??
	Channels ElevChannels
	err      error
}

// MoveContext contains information that needs to passed when moving up or down
type MoveContext struct {
	TargetFloor int
	err         error
}

/*NOTE: I don't quite understand why the context functions are of String()*/

func (ctx *InitContext) String() string {
	return fmt.Sprintf("InitContext")
}

func (ctx *MoveContext) String() string {
	return "Move context called"
}

func (a *InitAction) Execute(elev *ElevatorMachine, eventCtx EventContext) EventType {
	fmt.Println("Initializing elevator...")

	fmt.Println("Stopping motor")
	SetMotorDirection(MD_Stop)
	fmt.Println("Setting door lamp false")
	SetDoorOpenLamp(false)

	// TODO: Set the top floor to be config.NumFloors-1 at some point
	elev.TopFloor = 3
	elev.BottomFloor = 0

	// Give the elevator all the necessary channels
	elev.Channels = eventCtx.(*InitContext).Channels

	switch GetFloor() {

	// If not at floor -> move down to floor
	case -1:
		for GetFloor() == -1 {
			SetMotorDirection(MD_Down)
			time.Sleep(time.Millisecond * 20) //SHOULD BE POLL_RATE!!
		}
		SetMotorDirection(MD_Stop)
		return CloseDoors

	case elev.BottomFloor:
		elev.AtBottom = true
		elev.AtTop = false
	case elev.TopFloor:
		elev.AtBottom = false
		elev.AtTop = true
	default:
		elev.AtBottom = false
		elev.AtTop = false
		// TODO: what to do then?
	}

	switch GetObstruction() { //getObstruction is probably not exported
	case true:
		// Set current state to be obstructed here?
		fmt.Println("Elevator initialized, but obstructed!")
		return Obstruct
	case false:
		// Do nothing?
	}

	fmt.Println("Elevator initialized")
	return SetInitialized
}

func (a *OnObstrAction) Execute(elev *ElevatorMachine, eventCtx EventContext) EventType {
	fmt.Println("Elevator obstructed action called")

	// TODO: Create the functionality for this event

	return NoOp
}

func (a *MovingDownAction) Execute(elev *ElevatorMachine, eventCtx EventContext) EventType {
	fmt.Println("Elevator moving down action called")

	// TODO: Implement functionality for this event
	targetFloor := eventCtx.(*MoveContext).TargetFloor
	SetMotorDirection(MD_Down)

	for {
		select {
		// If new floor detected
		case a := <-elev.Channels.drv_floors:
			fmt.Printf("%+v\n", a)
			SetFloorIndicator(a)

			switch a {
			case elev.BottomFloor:
				elev.AtBottom = true
				SetMotorDirection(MD_Stop)
				fmt.Println("Arrived at ground floor")
				return ArriveAtFloor
			/* Might not need this case
			case elev.BottomFloor:
				elev.AtBottom = true
			*/
			case targetFloor:
				SetMotorDirection(MD_Stop)
				fmt.Printf("Arrived at floor %v", a)
				return ArriveAtFloor
			default:
				elev.AtTop = false
				elev.AtBottom = false
			}

		// If obstructed
		case a := <-elev.Channels.drv_obstr:
			fmt.Printf("%+v\n", a)
			return Obstruct
		}
	}
}

func (a *MovingUpAction) Execute(elev *ElevatorMachine, eventCtx EventContext) EventType {
	fmt.Println("Elevator moving up action called")

	targetFloor := eventCtx.(*MoveContext).TargetFloor
	SetMotorDirection(MD_Up)

	for {
		select {
		// If new floor detected
		case a := <-elev.Channels.drv_floors:
			fmt.Printf("%+v\n", a)
			SetFloorIndicator(a)

			switch a {
			case elev.TopFloor:
				elev.AtTop = true
				SetMotorDirection(MD_Stop)
				fmt.Println("Arrived at top floor")
				return ArriveAtFloor
			/* Might not need this case
			case elev.BottomFloor:
				elev.AtBottom = true
			*/
			case targetFloor:
				SetMotorDirection(MD_Stop)
				fmt.Printf("Arrived at floor %v", a)
				return ArriveAtFloor
			default:
				elev.AtTop = false
				elev.AtBottom = false
			}

		// If obstructed
		case a := <-elev.Channels.drv_obstr:
			fmt.Printf("%+v\n", a)
			return Obstruct
		}
	}
}

func (a *AtFloorClosedAction) Execute(elev *ElevatorMachine, eventCtx EventContext) EventType {
	SetDoorOpenLamp(false)
	return NoOp
}

func (a *AtFloorOpenAction) Execute(elev *ElevatorMachine, eventCtx EventContext) EventType {
	SetDoorOpenLamp(true)
	return NoOp
}

// StopAction represents the action executed when entering the stop state
type StopAction struct{}

func (a *StopAction) Execute(eventCtx EventContext) EventType {
	fmt.Println("The elevator has stopped")
	// TODO: What does it do when it stops?
	return NoOp
}

/*COMMENT ON ACTION STRUCTS:
All of these might not be necessary. It's possible we don't actually have to perform
an action when in certain states. However it is likely the action will only be to
continuously monitor the elevators physical state and act upon changes to it.*/

// InitAction should be used when initializing the elevator
type InitAction struct{}

// OnObstrAction should wait for the elevator to not be obstructed anymore
type OnObstrAction struct{}

// MovingDownAction represents the action executed when entering the MovingDown state
type MovingDownAction struct{}

// MovingUpAction represents the action executed when entering the MovingUp state
type MovingUpAction struct{}

// AtFloorClosedAction represents the action executed when in this state
type AtFloorClosedAction struct{}

// AtFloorOpenAction represents the action executed when in this state
type AtFloorOpenAction struct{}
