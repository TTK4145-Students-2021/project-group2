package elevio

import (
	"fmt"
)

// This has yet to be implemented in the RunElevator function
// ElevatorChannels contains the channels from elevio

type ElevChannels struct {
	drv_buttons <-chan ButtonEvent
	drv_floors  <-chan int
	drv_obstr   <-chan bool
}

func RunElevator(port string, numFloors int) {

	Init("localhost:"+port, numFloors)

	/*Initialize the state of the elevator. Wether we create a passable object or just use the
	fsm file (which is possible!) is up to debate. What if we want to create multiple elevators?*/
	//elevator := initElevator()
	//fmt.Printf("%+v\n", elevator.State)

	//Initialize some order list
	//Create functions to update that order list
	//This particular one handles the sequence of executions that has been assigned to this particular elevator

	var dir MotorDirection = MD_Stop
	SetMotorDirection(dir)

	drv_buttons := make(chan ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	//state := make(chan *Elevator) // For checking the current state of the elevator across the system
	// To make this channel bidirectional, don't pass a <- operator!

	go PollButtons(drv_buttons)
	go PollFloorSensor(drv_floors)
	go PollObstructionSwitch(drv_obstr)

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			// Send new order to newOrder channel
			// if new order at current floor and dir, dont send
			//handleButtonEvent(a)

		case a := <-drv_floors:

			fmt.Printf("%+v\n", a)
			/*
				if a == numFloors-1 {
					dir = MD_Down
				} else if a == 0 {
					dir = MD_Up
				}
			*/
			//elevio.SetMotorDirection(dir)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				SetMotorDirection(MD_Stop)
			} else {
				SetMotorDirection(dir)
			}
		}
	}
}

func handleButtonEvent(buttonEvent ButtonEvent) {

	fmt.Printf("%v\n", buttonEvent.Button)
	fmt.Printf("%v\n", buttonEvent.Floor)

	button := buttonEvent.Button
	//floor := buttonEvent.Floor

	switch button {
	case BT_HallUp:
		fmt.Printf("Button press up")

	case BT_HallDown:
		fmt.Printf("Button press down")

	case BT_Cab:
		fmt.Printf("Button press cab")
		//handleCabOrder(floor)
	}
}

/*

func gotoFloor(targetFloor int, drv_floors <-chan int, drv_obstr <-chan bool) int {

	//Checks for status updates while operating the motor towards the target destination
	for {
		select {
		// Move to target floor
		case currentFloor := <-drv_floors:
			SetFloorIndicator(currentFloor) //Set current floor lamp
			switch {
			case currentFloor == targetFloor:
				SetMotorDirection(MD_Stop)
				SetElevState(IDLE)
				return 0 //SUCCESS
			case currentFloor < targetFloor:
				SetMotorDirection(MD_Up)
				SetElevState(GOING_UP)
			case currentFloor > targetFloor:
				SetMotorDirection(MD_Down)
				SetElevState(GOING_DOWN)
			}

		// Stop execution if obstructed
		case <-drv_obstr:
			SetMotorDirection(MD_Stop)

			//Handle obstruction??

			return 1 //SOME ERROR CODE
		}
	}
}

func handleCabOrder(targetFloor int) {

	err := validateCabOrder(targetFloor)
	if err != 0 { //If err != OK is the intended code. But i need to find information about this
		//handle and display error code
		//return something
	}

	// TODO: Turn on light
	SetButtonLamp(BT_Cab, targetFloor, true)
	closeDoors()

	gotoFloor(targetFloor)
	// What happens if this does not execute properly? Write in a recovery routine!

	SetButtonLamp(BT_Cab, targetFloor, false)
	SetDoorOpenLamp(true)

}

func validateCabOrder(targetFloor int) int {

	if targetFloor > 4-1 || targetFloor < 0 {
		return 1 //Some error code that says the floor does not exist.
	} else {
		return 0
	}
}

*/
