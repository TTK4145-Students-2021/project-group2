package main

import (
	"fmt"

	"./elevio"
)

func main() {

	elevatorID int = 10
	port int = 10000

	//communication.SetupServer()
	//

	//elevio.SetupElevator()
	//elevio.InitElevatorController()
	elevio.initElevator(elevatorID, port)
	elevio.RunElevator(elevatorID, port)




	//elevio.StopElevator(elevatorID, port)
}
