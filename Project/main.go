package main

import (
	"os"

	"./config"
	"./elevio"
)

func main() {

	port int = 10000

	numFloors := config.NumFloors //Dette bør tilhøre en heis. Ikke selvfølge at alle heiser har like mange etasjer
	//communication.SetupServer()

	port := os.Args[1]
	elevio.RunElevator("localhost:"+port, numFloors)

	newOrder := make(chan elevio.ButtonEvent)
	assignedOrder := make(chan elevio.ButtonEvent)
}
