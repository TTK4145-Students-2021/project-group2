package main

import (
	"os"

	"./config"
	"./elevio"
)

func main() {

	numFloors := config.NumFloors

	port := os.Args[1]
	elevio.RunElevator("localhost:"+port, numFloors)

	newOrder := make(chan elevio.ButtonEvent)
	assignedOrder := make(chan elevio.ButtonEvent)

}
