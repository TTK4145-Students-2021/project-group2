package elevator

import (
	"fmt"
	"sync"
	"time"

	"../config"
)

/*=============================================================================
 * @Description:
 * Contains a struct that represents the physical state of the elevator as a
 * State Machine and operations that can be performed on that struct. The file also contains
 * functions to handle transitions between states.
 *
 * For more information and a detailed explanation on the state machine framework
 * used in this code, see the package Readme.md
 *
 * @Author: Eldar Sandanger
/*=============================================================================
*/

const _floorPollRate = time.Millisecond * 20
const _obstrPollRate = time.Millisecond * 20
const NoFloor int = -1

// StateType is an extensible state type in the elevator
type StateType string

// EventType is an extensible event type in the elevator
type EventType string

const (
	// States the elevator can be set to
	Uninitialized      StateType = "Uninitialized"
	Initializing       StateType = "Initializing"
	Moving             StateType = "Moving"
	AtFloorDoorsClosed StateType = "AtFloorDoorsClosed"
	AtFloorDoorsOpen   StateType = "AtFloorDoorsOpen"

	// Events that can be processed
	NoOp           EventType = ""
	Initialize     EventType = "Initialize"
	SetInitialized EventType = "SetInitialized"
	Move           EventType = "Move"
	ArriveAtFloor  EventType = "ArriveAtFloor"
)

// Action represents an action to be executed upon a state transition
type Action interface {
	Execute(elev *ElevatorMachine, eventCtx EventContext) EventType
}

// State binds a state with an action and a set of events it can handle in that given state
type State struct {
	Action Action
	Events Events
}

// Events represents the mapping of events that can be performed in given states
type Events map[EventType]StateType

// States represents a mapping of states and their implementations
type States map[StateType]State

// ElevChannels contain all the channels needed for the elevator to work
type ElevChannels struct {
	drv_floors <-chan int
	drv_obstr  <-chan bool
}

// ElevatorMachine contains all information meeded to represent the physical state
type ElevatorMachine struct {

	// Track if elevator available
	Available bool

	// Previous state
	Previous StateType

	// Current state
	Current StateType

	// Holds the configuration of states and events handled by the state machine
	States States

	// Contains all channels the elevator needs
	Channels ElevChannels

	// Track current floor
	CurrentFloor int

	// Other relevant floor data
	TopFloor    int
	BottomFloor int

	// Ensures exclusive access to the data contained within the struct
	mutex sync.Mutex
}

// NewElevatorMachine initializes and returns a new elevator object
func NewElevatorMachine() *ElevatorMachine {
	return &ElevatorMachine{

		// This maps all the possible states, actions and events appropriately.
		// The pattern might seem unfamiliar, but a detailed explanation is provided in this
		// packages Readme as it is vital for understanding how the code works.
		States: States{

			Uninitialized: State{
				Events: Events{
					Initialize: Initializing,
				},
			},

			Initializing: State{
				Action: &InitAction{},
				Events: Events{
					SetInitialized: AtFloorDoorsOpen,
				},
			},

			Moving: State{
				Action: &MovingAction{},
				Events: Events{
					ArriveAtFloor: AtFloorDoorsOpen,
				},
			},

			AtFloorDoorsOpen: State{
				Action: &AtFloorOpenAction{},
				Events: Events{
					Move: Moving,
				},
			},
		},

		// Other default parameters
		Current:     Uninitialized,
		Available:   false,
		TopFloor:    config.NumFloors - 1,
		BottomFloor: config.BottomFloor,
	}
}

// SendEvent sends an event to the state machine
// Triggering an event can cause a chain of events. This function ensures that every event in the chain
// is processed until a "NoOp"-event is returned, indicating that there are No Operations left to perform.
func (elev *ElevatorMachine) SendEvent(event EventType, eventCtx EventContext) error {

	for {

		// Determine the next state for the event given the elevators current state
		elev.mutex.Lock()
		nextState, err := elev.getNextState(event)
		if err != nil {
			return ErrEventRejected
		}

		// Identify the state definition for the next state
		state, ok := elev.States[nextState]
		if !ok || state.Action == nil {
			// Configuration error
			// If state.Action == nil, this implies there are no actions in the next state
			// ... this means you have to make one!
			// If not 'ok', this implies the action returned is not a legal action in the next state
		}

		// Transition to the next state
		elev.Previous = elev.Current
		elev.Current = nextState
		elev.mutex.Unlock()

		// Execute the next states action
		nextEvent := state.Action.Execute(elev, eventCtx)
		if nextEvent == NoOp {
			return nil
		}

		event = nextEvent

	}
}

// Return next state for given event, or an error if the event is illegal in the current state
func (elev *ElevatorMachine) getNextState(event EventType) (StateType, error) {
	if state, ok := elev.States[elev.Current]; ok {
		if state.Events != nil {
			if next, ok := state.Events[event]; ok {
				return next, nil
			}
		}
	}
	return elev.Current, ErrEventRejected
}

// EventContext represents the context to be passed to the action implementation
type EventContext interface{}

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

// Execute the initializing action for the elevator
func (a *InitAction) Execute(elev *ElevatorMachine, eventCtx EventContext) EventType {
	elev.mutex.Lock()
	defer elev.mutex.Unlock()

	fmt.Println("Initializing elevator...")
	SetMotorDirection(MD_Stop)
	SetDoorOpenLamp(false)

	// Give the elevator all the necessary channels
	elev.Channels = eventCtx.(*InitContext).Channels

	// If not initially at floor -> move down to closest floor
	for GetFloor() == NoFloor {
		SetMotorDirection(MD_Down)
		time.Sleep(_floorPollRate)
	}
	SetMotorDirection(MD_Stop)
	elev.CurrentFloor = GetFloor()

	fmt.Println("Elevator initialized")
	return SetInitialized
}

// Execute a move order received by the elevator machine
func (a *MovingAction) Execute(elev *ElevatorMachine, eventCtx EventContext) EventType {
	fmt.Println("Elevator executing order.")
	elev.mutex.Lock()
	elev.Available = false
	targetFloor := eventCtx.(*MoveContext).TargetFloor
	elev.CurrentFloor = GetFloor()
	elev.mutex.Unlock()

	fmt.Printf("Target Floor: %v | ", targetFloor)
	fmt.Printf("Current Floor: %v, \n", elev.CurrentFloor)

	// Decide direction or stop operation if already at target floor
	var dir MotorDirection
	switch tf := targetFloor; {
	case tf == elev.CurrentFloor:
		fmt.Println("Elevator already at destination")
		return ArriveAtFloor
	case tf < elev.CurrentFloor:
		fmt.Println("Elevator moving down")
		dir = MD_Down
	case tf > elev.CurrentFloor:
		fmt.Println("Elevator moving up")
		dir = MD_Up
	}
	SetMotorDirection(dir)

	// Wait to close door if obstructed
	for GetObstruction() {
		fmt.Println("Elevator obstructed! Waiting...")
		time.Sleep(_obstrPollRate)
	}
	SetDoorOpenLamp(false)

	for {
		select {
		// If new floor detected
		case newFloor := <-elev.Channels.drv_floors:

			SetFloorIndicator(newFloor)

			if newFloor == elev.CurrentFloor {
				// Likely read-error. Do nothing
				break
			} else {
				elev.mutex.Lock()
				elev.CurrentFloor = newFloor
				elev.mutex.Unlock()
			}

			switch newFloor {
			case elev.BottomFloor:
				SetMotorDirection(MD_Stop)
				fmt.Println("Arrived at ground floor")
				return ArriveAtFloor
			case elev.TopFloor:
				SetMotorDirection(MD_Stop)
				fmt.Println("Arrived at top floor")
				return ArriveAtFloor
			case targetFloor:
				SetMotorDirection(MD_Stop)
				fmt.Printf("Arrived at target floor\n", a)
				return ArriveAtFloor
			default:
				// Do nothing
			}
		}
		// POTENTIAL IMPROVEMENT: Implement a channel than enables breaking the elevators operation midway
	}
}

// Execute when AtFloorDoorOpen state is set
func (a *AtFloorOpenAction) Execute(elev *ElevatorMachine, eventCtx EventContext) EventType {
	elev.mutex.Lock()
	defer elev.mutex.Unlock()
	SetDoorOpenLamp(true)
	fmt.Println("Elevator doors open")
	elev.Available = true
	return NoOp
}

/*COMMENT ON ACTION STRUCTS:
These action structs are pretty much only used to differentiate between different
execution-functions.

All of these might not be necessary. It's possible we don't actually have to perform
an action when in certain states. However it is likely the action will only be to
continuously monitor the elevators physical state and act upon changes to it.*/

// InitAction should be used when initializing the elevator
type InitAction struct{}

// OnObstrAction should wait for the elevator to not be obstructed anymore
type OnObstrAction struct{}

// MovingAction represents the action executed when entering the Moving state
type MovingAction struct{}

// AtFloorClosedAction represents the action executed when in this state
type AtFloorClosedAction struct{}

// AtFloorOpenAction represents the action executed when in this state
type AtFloorOpenAction struct{}
