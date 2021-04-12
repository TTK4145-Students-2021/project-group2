package elevator

import (
	"testing"
)

// TestElevatorMachine used to test the elevator!
func TestElevatorController(t *testing.T) {

	// Setup all the channels we need
	elev_available := make(chan bool) // Door_open*
	current_floor := make(chan int)
	button_action := make(chan ButtonEvent)
	order := make(chan int)
	order_response := make(chan error) //TODO: Order response can be a const int instead

	// Channels to be passed controller
	ctrChans := ControllerChannels{
		elev_available:         elev_available,
		current_floor:          current_floor,
		redirect_button_action: button_action,
		receive_order:          order,
		respond_order:          order_response,
	}

	orderChans := &OrderChannels{
		elev_available: elev_available,
		current_floor:  current_floor,
		button_event:   button_action,
		send_order:     order,
		order_response: order_response,
	}

	controller := NewController(ctrChans)
	go controller.Run()

	elev := controller.GetElevator()

	elevStatus := &ElevatorStatus{
		ID:           1,
		Available:    elev.Available,
		CurrentFloor: elev.CurrentFloor,
		AtTop:        elev.AtTop,
		AtBottom:     elev.AtBottom,
	}

	go OrderSimulator(orderChans, elevStatus)

	for {

	}

}
