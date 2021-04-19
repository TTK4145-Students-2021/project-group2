package messages

import (
	"time"

	"../config"
)

type ButtonType int

const (
	UNDEFINED   ButtonType = -1
	BT_HallUp   ButtonType = 0
	BT_HallDown ButtonType = 1
	BT_Cab      ButtonType = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

type LampUpdate_message struct {
	Floor  int
	Button ButtonType
	Turn   bool
}

type HallOrder struct {
	HasOrder   bool
	Floor      int
	Direction  ButtonType //up = 0, down = 1
	VersionNum int
	Costs      [config.NumElevators]int
	TimeStamp  time.Time
}

type OrderList [config.OrderListLength] HallOrder


type ElevatorStatus struct {
	ID           int
	Pos          int
	OrderList    OrderList
	IsOnline     bool
	DoorOpen     bool
	CabOrders    [config.NumFloors]bool
	IsAvailable  bool
	IsObstructed bool
	Timestamp    time.Time
}
