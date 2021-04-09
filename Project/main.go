// run main with as follows:
// go run main.go [port] [elevatorID]

package main

import (
	"fmt"

	"./config"
	"./elevio"
)

func main() {
	//var HallOrder = make(chan elevio.ButtonEvent)
	//var ReceivedElevatorUpdate = make(chan orders.ElevatorStatus)
	//var NewFloor = make(chan int)
	//var IsDoorOpen = make(chan bool)

	config.InitConfig()

	fmt.Println(config.NumFloors, config.NumElevators, config.Port, config.ID)

	elevio.RunElevator("localhost:"+config.Port, config.NumFloors)
	orders.runOrders() // Send correct channel

	//newOrder := make(chan elevio.ButtonEvent)
	//assignedOrder := make(chan elevio.ButtonEvent)

	// TestElevatorMachine() ??

}
