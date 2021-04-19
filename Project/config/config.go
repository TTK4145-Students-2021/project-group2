package config

import (
	"fmt"
	"os"
	"time"
)

//Package for common variables to be defined across the application
var Port string

const ID int = 2 //the order module is written so that this is the first elevator. (i.e it has -1 in all indexes)
const NumFloors int = 4
const BottomFloor int = 0

const NumElevators int = 3

//var HallOrder = make(chan elevio.ButtonEvent)  //*****Commented out beacuse of include problems. Possible Fix: intialize in a higher order file.
//var ReceivedElevatorUpdate = make(chan orders.ElevatorStatus)  //*****Commented out beacuse of include problems. Possible Fix: intialize in a higher order file.
var NewFloor = make(chan int)
var isDoorOpen = make(chan bool)

const BcastPort int = 12345
const PeersPort int = 54321
const BcastIntervall time.Duration = 1000 * time.Millisecond

const CheckOfflineIntervall time.Duration = time.Second
const OfflineTimeout time.Duration = 5 * time.Second
const OrderTimeOut time.Duration = 5 * time.Second

var SimPort string

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func InitConfig() {
	if len(os.Args) < 1 {
		SimPort = "15657"
	} else {
		SimPort = os.Args[1]
	}
	fmt.Println("SimPort:", SimPort, "\tID: ", ID)
}
