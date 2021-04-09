package elevator

import (
	"fmt"
	"sync"
	"time"
)

// language spec: https://golang.org/ref/spec#Function_literals
// Framework used: https://venilnoronha.io/a-simple-state-machine-framework-in-go

// StateType is an extensible state type in the elevator
type StateType string

// EventType is an extensible event type in the elevator
type EventType string

const (
	// States of the elevator
	Uninitialized      StateType = "Uninitialized"
	Initializing       StateType = "Initializing"
	Moving             StateType = "Moving"
	AtFloorDoorsClosed StateType = "AtFloorDoorsClosed"
	AtFloorDoorsOpen   StateType = "AtFloorDoorsOpen"

	// Events that can be processed
	NoOp           EventType = ""
	Initialize     EventType = "Initialize"
	SetInitialized EventType = "SetInitialized"
	OpenDoors      EventType = "OpenDoors"
	CloseDoors     EventType = "CloseDoors"
	Move           EventType = "Move"
	ArriveAtFloor  EventType = "ArriveAtFloor"
	StopAtFloor    EventType = "StopAtFloor"
)


// ElevChannels contain all the channels needed for the elevator to work
type ElevChannels struct {
	drv_floors <-chan int
	drv_obstr  <-chan bool
}


// Events represents the mapping of events that can be performed in given states
type Events map[EventType]StateType

// States represents a mapping of states and their implementations
type States map[StateType]State

// State binds a state with an action and a set of events it can handle in that given state
type State struct {
	Action Action
	Events Events
}



// ElevatorMachine represents the elevator itself
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
	TotalFloors int
	TopFloor    int
	BottomFloor int
	AtTop       bool
	AtBottom    bool

	// Ensures only one event is processed by the machine at a time
	mutex sync.Mutex
}



// NewElevatorMachine initializes and returns a new elevator object
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
					SetInitialized: AtFloorDoorsClosed,
				},
			},

			Moving: State{
				Action: &MovingAction{},
				Events: Events{
					ArriveAtFloor: AtFloorDoorsClosed,
				},
			},

			AtFloorDoorsClosed: State{
				Action: &AtFloorClosedAction{},
				Events: Events{
					OpenDoors: AtFloorDoorsOpen,
					Move:      Moving,
				},
			},

			AtFloorDoorsOpen: State{
				Action: &AtFloorOpenAction{},
				Events: Events{
					CloseDoors: AtFloorDoorsClosed,
				},
			},
		},

		Current:   Uninitialized,
		Available: false,
	}
}



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
	return elev.Current, ErrEventRejected //This is most likely wrong!! //TODO: CHECK THIS!
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

		if nextState == Moving {
			if elev.AtTop {
				return ErrIllegalMoveUp
			} else if elev.AtBottom {
				return ErrIllegalMoveDown
			}
		}

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

	fmt.Println("Getting floor")
	switch GetFloor() {

		// If not at floor -> move down to floor
		case NoFloor:
			for GetFloor() == NoFloor {
				SetMotorDirection(MD_Down)
				time.Sleep(time.Millisecond * 20) //SHOULD BE POLL_RATE!!
			}
			SetMotorDirection(MD_Stop)
			elev.CurrentFloor = GetFloor()
			// TODO: REPAIR THIS LOGIC
			// It skips the part where it set which floor we're at

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

	if GetObstruction() {
		fmt.Println("Note: Elevator Obstructed!")
	}

	fmt.Println("Elevator initialized")
	return SetInitialized
}

// Execute a move order received by the elevator machine
func (a *MovingAction) Execute(elev *ElevatorMachine, eventCtx EventContext) EventType {
	fmt.Println("Elevator executing order.")

	targetFloor := eventCtx.(*MoveContext).TargetFloor
	elev.CurrentFloor = GetFloor()

	fmt.Printf("Target Floor: %v | ", targetFloor)
	fmt.Printf("Current Floor: %v,", elev.CurrentFloor)

	// Decide direction
	var dir MotorDirection
	switch tf := targetFloor; {
	case tf == elev.CurrentFloor:
		fmt.Println("Elevator already at destination")
		return OpenDoors // Simply open doors
	case tf < elev.CurrentFloor:
		fmt.Println("Elevator moving down")
		dir = MD_Down
	case tf > elev.CurrentFloor:
		fmt.Println("Elevator moving up")
		dir = MD_Up
	}
	SetMotorDirection(dir)

	waitIfObstructed(elev)

	for {
		select {
		// If new floor detected
		case a := <-elev.Channels.drv_floors:
			fmt.Printf("%+v\n", a)
			SetFloorIndicator(a)

			if a == elev.CurrentFloor {
				// Likely read-error. Do nothing
				break
			} else {
				elev.CurrentFloor = a
			}

			switch a {
			case elev.BottomFloor:
				elev.AtTop = false
				elev.AtBottom = true
				SetMotorDirection(MD_Stop)
				fmt.Println("Arrived at ground floor")
				return ArriveAtFloor
			case elev.TopFloor:
				elev.AtTop = true
				elev.AtBottom = false
				SetMotorDirection(MD_Stop)
				fmt.Println("Arrived at top floor")
				return ArriveAtFloor
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
			if a {
				fmt.Println("Elevator obstructed!")
				SetMotorDirection(MD_Stop)
			} else {
				fmt.Println("Elevator not obstructed anymore.")
				SetMotorDirection(dir)
			}
		}
	}
}

// Execute when AtFloorDoorClosed state is set
func (a *AtFloorClosedAction) Execute(elev *ElevatorMachine, eventCtx EventContext) EventType {
	elev.Available = false
	SetDoorOpenLamp(false)
	return NoOp

	// TODO: What if you want to open the door or do something else here???
}

// Execute when AtFloorDoorOpen state is set
func (a *AtFloorOpenAction) Execute(elev *ElevatorMachine, eventCtx EventContext) EventType {
	elev.Available = true
	fmt.Println("Elevator set available")
	SetDoorOpenLamp(true)
	fmt.Println("Elevator doors open")
	return NoOp
}

// Used to halt move operation if elevator is obstructed upon call
func waitIfObstructed(elev *ElevatorMachine) {
	for GetObstruction() {
		fmt.Println("Elevator obstructed!")
		SetMotorDirection(MD_Stop)
		select {
		case a := <-elev.Channels.drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				SetMotorDirection(MD_Stop)
			} else {
				SetMotorDirection(MD_Down)
			}
		}
	}
}


// EventContext represents the context to be passed to the action implementation
type EventContext interface{}

// Action represents the action to be executed in a given state
type Action interface {
	// Pass an eventContext and receive an event type
	Execute(elev *ElevatorMachine, eventCtx EventContext) EventType
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

// MovingDownAction represents the action executed when entering the Moving state
type MovingAction struct{}

// MovingDownAction represents the action executed when entering the MovingDown state
type MovingDownAction struct{}

// MovingUpAction represents the action executed when entering the MovingUp state
type MovingUpAction struct{}

// AtFloorClosedAction represents the action executed when in this state
type AtFloorClosedAction struct{}

// AtFloorOpenAction represents the action executed when in this state
type AtFloorOpenAction struct{}
