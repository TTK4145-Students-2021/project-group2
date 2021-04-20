package config

import (
	"fmt"
	"os"
	"time"
)

/*
*=============================================================================
 * @Description:
 * Contains globale public variables used throughout the system
 * 
 *
 * 
/*=============================================================================
*/

var Port string

const ID int = 0 // Elevator ID is zero-indexed
const NumFloors int = 4
const BottomFloor int = 0

const NumElevators int = 3

const BcastPort int = 12345
const BcastIntervall time.Duration = 100 * time.Millisecond

const CheckOfflineIntervall time.Duration = 250 * time.Millisecond
const OfflineTimeout time.Duration = 5 * time.Second
const OrderTimeOut time.Duration = 4 * time.Second
const LampUpdateIntervall time.Duration = 500 * time.Millisecond

const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

var SimPort string

func InitConfig() {
	if len(os.Args) < 1 {
		SimPort = "15657"
	} else {
		SimPort = os.Args[1]
	}
	fmt.Println("SimPort:", SimPort, "\tID: ", ID)
}
