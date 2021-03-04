package elevio

import "fmt"

func RunElevator(id int, port int) {

	numFloors := 4

	Init("localhost:15657", numFloors)

	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go PollButtons(drv_buttons)
	go PollFloorSensor(drv_floors)
	go PollObstructionSwitch(drv_obstr)
	go PollStopButton(drv_stop)

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)

			//HandleButtonEvent()

			//TEST:
			//AssignOrder(a.Button, a.Floor, true)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				SetMotorDirection(MD_Stop)
			} else {
				SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := ButtonType(0); b < 3; b++ {
					SetButtonLamp(b, f, false)
				}
			}
		}
	}
}

func gotoFloor(elev *Elevator, targetFloor int, drv_floors chan<-, drv_obstr chan<-) int {

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