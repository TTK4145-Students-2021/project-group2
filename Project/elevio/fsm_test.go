package elevio

import (
	"fmt"
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

	fmt.Println(&drv_buttons)
	fmt.Println(&drv_floors)
	fmt.Println(&drv_obstr)

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

	/*
		// Move elevator up
		err := elevator.SendEvent(MoveUp, nil)
		if err != nil {
			t.Errorf("Couldn't move elevator, err: %v", err)
		}

		time.Sleep(2)

		err := elevator.SendEvent(Stop, nil)
		if err != nil {
			t.Errorf("Couldn't stop elevator, err: %v", err)
		}
	*/
}
