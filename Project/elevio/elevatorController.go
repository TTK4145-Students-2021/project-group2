package elevio

import (
	"fmt"
)

func RunElevator(port string, numFloors int) {

	Init("localhost:"+port, numFloors)

	var dir MotorDirection = MD_Stop
	SetMotorDirection(dir)

	drv_buttons := make(chan ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	//drv_stop := make(chan bool)

	go PollButtons(drv_buttons)
	go PollFloorSensor(drv_floors)
	go PollObstructionSwitch(drv_obstr)
	//go PollStopButton(drv_stop)

	for {
		select {
		case a := <-drv_buttons:
			// Send new order to newOrder channel
			// if new order at current floor and dir, dont send
			handleButtonEvent(a)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			if a == numFloors-1 {
				dir = MD_Down
			} else if a == 0 {
				dir = MD_Up
			}
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

	switch buttonEvent.Button {
	case BT_HallUp:
		fmt.Printf("Button press up")

	case BT_HallDown:
		fmt.Printf("Button press down")

	case BT_Cab:
		fmt.Printf("Button press cab")
	}
}

func gotoFloor(targetFloor int, drv_floors <-chan int, drv_obstr <-chan bool) int {

	//TODO: Check if targetfloor exists
	if targetFloor > NumberOfFloors-1 || targetFloor < 0 {
		return 1 //Some error code that says the floor does not exist.
	}

	//Checks for status updates while operating the motor towards the target destination
	for {
		select {
		// Move to target floor
		case currentFloor := <-drv_floors:
			switch {
			case currentFloor == targetFloor:
				SetMotorDirection(MD_Stop)
				return 0 //SUCCESS
			case currentFloor < targetFloor:
				SetMotorDirection(MD_Up)
			case currentFloor > targetFloor:
				SetMotorDirection(MD_Down)
			}

		// Stop execution if obstructed
		case <-drv_obstr:
			SetMotorDirection(MD_Stop)
			return 1 //SOME ERROR CODE
		}
	}
}
