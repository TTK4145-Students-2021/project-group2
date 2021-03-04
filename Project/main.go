package main

import (
	"os"

	"./config"
	"./elevio"
)

func main() {

	port int = 10000

	//communication.SetupServer()
	//

	//elevio.SetupElevator()
	//elevio.InitElevatorController()


	//elevio.StopElevator(elevatorID, port)
	numFloors := config.NumFloors

	port := os.Args[1]
	elevio.RunElevator("localhost:"+port, numFloors)

	newOrder := make(chan elevio.ButtonEvent)
	assignedOrder := make(chan elevio.ButtonEvent)
}
