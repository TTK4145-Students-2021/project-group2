package elevio

import (
	"strconv"
	"testing"
)

// TestElevatorMachine used to test the elevator!
func TestElevator(t *testing.T) {

	// Init elevator connection
	//config.InitConfig() Doesn't work yet
	//Init("localhost:"+config.Port, config.NumFloors)
	Init("localhost:"+strconv.Itoa(15657), 4)

	// First we declare ALL channels
	drv_buttons := make(chan ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)

	go PollButtons(drv_buttons)
	go PollFloorSensor(drv_floors)
	go PollObstructionSwitch(drv_obstr)

	// Create new elevator
	elevator := NewElevatorMachine()

	// Then we can portion the channels into designated structs
	elevChans := ElevChannels{
		drv_buttons: drv_buttons,
		drv_floors:  drv_floors,
		drv_obstr:   drv_obstr,
	}

	// Define the context for order creation.
	initCtx := &InitContext{
		Channels: elevChans,
	}

	// Initialize the elevator
	err := elevator.SendEvent(Initialize, initCtx)
	if err != nil {
		t.Errorf("Couldn't initialize elevator, err: %v", err)
	}

	moveUpCtx := &MoveContext{
		TargetFloor: 3,
	}

	// Move elevator up
	err = elevator.SendEvent(MoveUp, moveUpCtx)
	if err != nil {
		t.Errorf("Couldn't move elevator, err: %v", err)
	}

	moveDownCtx := &MoveContext{
		TargetFloor: 0,
	}

	err = elevator.SendEvent(MoveDown, moveDownCtx)
	if err != nil {
		t.Errorf("Couldn't move elevator, err: %v", err)
	}

	err = elevator.SendEvent(OpenDoors, nil)
	if err != nil {
		t.Errorf("Couldn't open doors, err: %v", err)
	}

	/*
		err = elevator.SendEvent(CloseDoors, nil)
		if err != nil {
			t.Errorf("Couldn't close doors, err: %v", err)
		}
	*/

	/*
		err = elevator.SendEvent(Stop, nil)
		if err != nil {
			t.Errorf("Couldn't stop elevator, err: %v", err)
		}
	*/
}
