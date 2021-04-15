package messages

import (
	"time"

	"../config"
)

type ButtonType_message int

const (
	UNDEFINED   ButtonType_message = -1
	BT_HallUp   ButtonType_message = 0
	BT_HallDown ButtonType_message = 1
	BT_Cab      ButtonType_message = 2
)

type ButtonEvent_message struct {
	Floor  int
	Button ButtonType_message
}

type HallOrder struct {
	HasOrder   bool
	Floor      int
	Direction  int //up = 0, down = 1
	VersionNum int
	Costs      [config.NumElevators]int
	TimeStamp  time.Time
}

type ElevatorStatus struct {
	ID          int
	Pos         int
	OrderList   [config.NumFloors*2 - 2]HallOrder //
	Dir         int                               //up = 0, down = 1, stop = 2
	IsOnline    bool
	DoorOpen    bool
	CabOrders   [config.NumFloors]bool
	IsAvailable bool
	IsObstructed bool
}
type ElevatorStatusTest struct {
	ID  int
	Pos int
	//OrderList   [config.NumFloors*2 - 2]HallOrder //
	Dir         int //up = 0, down = 1, stop = 2
	IsOnline    bool
	DoorOpen    bool
	CabOrders   [config.NumFloors]bool
	IsAvailable bool
}

type Inner struct {
	Str  string
	List [config.NumFloors]int
}

type Outer struct {
	ID  int
	Msg [config.NumFloors*2 - 2]HallOrder
	Pos int
	//OrderList   [config.NumFloors*2 - 2]HallOrder //
	Dir         int //up = 0, down = 1, stop = 2
	IsOnline    bool
	DoorOpen    bool
	CabOrders   [config.NumFloors]bool
	IsAvailable bool
}
