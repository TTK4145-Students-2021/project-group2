package elevio

import (
	"fmt"

	"./elevio"
)

func RunElevator(port string, numFloors int) {

	elevio.Init("localhost:"+port, numFloors)

	var dir elevio.MotorDirection = elevio.MD_Stop
	elevio.SetMotorDirection(dir)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)

	for {
		select {
		case a := <-drv_buttons:
			// Send new order to newOrder channel
			// if new order at current floor and dir, dont send

			fmt.Printf("%+v\n", a)
			elevio.SetButtonLamp(a.Button, a.Floor, true)

			//TEST:
			elevio.TakeOrder(a.Button, a.Floor, true)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			if a == numFloors-1 {
				dir = elevio.MD_Down
			} else if a == 0 {
				dir = elevio.MD_Up
			}
			elevio.SetMotorDirection(dir)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(dir)
			}
		}
	}
}
