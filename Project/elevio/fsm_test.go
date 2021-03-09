package elevio

import "testing"

// TestElevatorMachine used to test the elevator!
func TestElevator(t *testing.T) {

	chans := Channels{
		drv_buttons: make(chan ButtonEvent),
		drv_floors:  make(chan int),
		drv_obstr:   make(chan bool),
	}

	elevator := NewElevatorMachine()

	// Initialize the elevator
	err := elevator.SendEvent(Initialize, chans)
	if err != nil {
		t.Errorf("Couldn't initialize elevator, err: %v", err)
	}

	/*
		// Move elevator up
		err := elevator.SendEvent(MoveUp, nil)
		if err != nil {
			t.Errorf("Couldn't move elevator, err: %v", err)
		}

		time.sleep(2)

		err := elevator.SendEvent(Stop, nil)
		if err != nil {
			t.Errorf("Couldn't stop elevator, err: %v", err)
		}
	*/
}
