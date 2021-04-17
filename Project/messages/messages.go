package messages

import (
	"time"

	"../config"
)

type ButtonType_msg int

const (
	UNDEFINED   ButtonType_msg = -1
	BT_HallUp   ButtonType_msg = 0
	BT_HallDown ButtonType_msg = 1
	BT_Cab      ButtonType_msg = 2
)

type ButtonEvent_message struct {
	Floor  int
	Button ButtonType_msg
}

type LampUpdate_message struct {
	Floor  int
	Button ButtonType_msg
	Turn   bool
}

type HallOrder struct {
	HasOrder   bool
	Floor      int
	Direction  ButtonType_msg //up = 0, down = 1
	VersionNum int
	Costs      [config.NumElevators]int
	TimeStamp  time.Time
}

type ElevatorStatus struct {
	ID           int
	Pos          int
	OrderList    [config.NumFloors*2 - 2]HallOrder
	IsOnline     bool
	DoorOpen     bool
	CabOrders    [config.NumFloors]bool
	IsAvailable  bool
	IsObstructed bool
	Timestamp    time.Time
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
	Timestamp   time.Time
}
