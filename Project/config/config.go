package config

import (
	"os"
	"strconv"


)

//Package for common variables to be defined across the application
var Port string
var ID int = 1 //the order module is written so that this is the first elevator. (i.e it has -1 in all indexes)
const NumFloors int = 4

const NumElevators int = 3

//var HallOrder = make(chan elevio.ButtonEvent)  //*****Commented out beacuse of include problems. Possible Fix: intialize in a higher order file. 
//var ReceivedElevatorUpdate = make(chan orders.ElevatorStatus)  //*****Commented out beacuse of include problems. Possible Fix: intialize in a higher order file. 
var NewFloor = make(chan int)
var isDoorOpen = make(chan bool)

func InitConfig() {
	Port = os.Args[1]
	ID, _ = strconv.Atoi(os.Args[2])
}
