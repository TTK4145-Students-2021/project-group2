package orderDistribution

import(
	"time"
	"../config"
	"../messages"
)
/*
*=============================================================================
 * @Description:
 * Contains the decleratation for all types used in the orderDistribution module
 * 
 *
 * 
/*=============================================================================
*/

type ButtonEvent = messages.ButtonEvent
type ButtonType = messages.ButtonType
type LampUpdate = messages.LampUpdate
// type OrderList = messages.OrderList

const _orderListLength = config.NumFloors*2 - 2
const _numElevators = config.NumElevators
const _numFloors = config.NumFloors
const _thisID = config.ID

const _pollOfflineRate = config.CheckOfflineIntervall
const _lampUpdateRate = config.LampUpdateIntervall
const _offlineTimeout = config.OfflineTimeout
const _bcastRate = config.BcastIntervall
const _checkOfflineIntervall = config.CheckOfflineIntervall
const _orderTimeOut = config.OrderTimeOut

const UNDEFINED = messages.UNDEFINED
const BT_HallUp = messages.BT_HallUp
const BT_HallDown = messages.BT_HallDown
const BT_Cab = messages.BT_Cab

const _engineFailureTimeOut = 10   //seconds
const _selectBreakOutRate  = 500*time.Millisecond

type ElevatorStatus struct {
	ID           		 int
	Pos          		 int
	OrderList    		 [_orderListLength]HallOrder
	IsOnline     		 bool
	DoorOpen       		 bool
	CabOrders    		 [_numFloors]bool
	IsAvailable  		 bool
	IsObstructed 		 bool
	LastAlive    		 time.Time
}

type HallOrder struct {
	HasOrder   bool
	Floor      int
	Direction  ButtonType //up = 0, down = 1
	VersionNum int
	Costs      [_numElevators]int
	TimeStamp  time.Time
}

type OrderChannels struct {
	ButtonPress             <-chan ButtonEvent
	ReceivedElevatorUpdate 	<-chan ElevatorStatus    
	NewFloor                <-chan int                          
	DoorOpen                <-chan bool                        
	SendStatus              chan<- ElevatorStatus      
	GotoFloor              	chan<- int                         
	GetInitPosition		    chan<- bool
	DoorObstructed          <-chan bool
	UpdateLampMessage       chan<- LampUpdate
}
