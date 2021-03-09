package config

import (
	"os"
	"strconv"

	"./elevio"
	"./orders"
)

//Package for common variables to be defined across the application
var Port string
var ID int
var NumFloors int = 4
var NumElevators int = 3

var HallOrder = make(chan elevio.ButtonEvent)
var ReceivedElevatorUpdate = make(chan orders.ElevatorStatus)
var NewFloor = make(chan int)
var isDoorOpen = make(chan bool)

func InitConfig() {
	Port = os.Args[1]
	ID, _ = strconv.Atoi(os.Args[2])
}
