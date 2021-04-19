package orderDistribution

import(
	"time"
	"../messages"
	"errors"
	"../config"
	// "fmt"
)

type ElevatorStatus struct {
	ID           int
	Pos          int
	OrderList    OrderList
	IsOnline     bool
	DoorOpen     bool
	CabOrders    [_numFloors]bool
	IsAvailable  bool
	IsObstructed bool
	TimeOfLastCompletion time.Time 
	LastAlive    time.Time
}

func (es *ElevatorStatus) updateDoor(value bool) {
	es.DoorOpen = value
}

func (es *ElevatorStatus) updateFloor(pos int) {
	es.Pos = pos
}


func (es *ElevatorStatus) sendOutLightsUpdate(sendLampUpdate chan<- messages.LampUpdate) {
	for {
		for _, order := range es.OrderList {
			light := messages.LampUpdate{
				Floor:  order.Floor,
				Button: messages.ButtonType(order.Direction),
				Turn:   order.HasOrder,
			}
			sendLampUpdate <- light
		}
		for i, order := range es.CabOrders {

			light := messages.LampUpdate{
				Floor:  i,
				Button: messages.BT_Cab,
				Turn:   order,
			}
			sendLampUpdate <- light
		}
		time.Sleep(_lampUpdateRate)
	}
}
/*
func (es *ElevatorStatus) sendOutStatus(channel chan<- ElevatorStatus) {
	es.Timestamp = time.Now()
	channel <- es
}
*/

func (es *ElevatorStatus) checkEngineFailure(assignedFloor int, send_status chan<- ElevatorStatus) {
	timeSinceMovment := 0
	lastFloor := -1
	for {
		if es.Pos == assignedFloor {
			return
		}
		if es.Pos != lastFloor {
			lastFloor = es.Pos 
			 // TODO: RETT DENNE
		} else {
			timeSinceMovment += 1
			if timeSinceMovment > engineFailureTimeOut {
				err := errors.New("engine failure")
				panic(err)
			}
		}
		time.Sleep(time.Second)
	}
}

func (es *ElevatorStatus) broadcast(channel chan<- ElevatorStatus) {
	for {
		time.Sleep(_bcastRate)
		es.LastAlive = time.Now()
		channel <- *es
		printElevatorStatus(*es)
		// fmt.Println("BCAST")
	}
}

func (es *ElevatorStatus) broadcastOneTime(channel chan<- ElevatorStatus) {

	es.LastAlive = time.Now()
	channel <- *es
}

//updateOrderListButton
func (es *ElevatorStatus) updateOrders(btnEvent messages.ButtonEvent) {

	if btnEvent.Button == messages.BT_Cab {
		es.CabOrders[btnEvent.Floor] = true

	} else {
		idx := floorToOrderListIdx(btnEvent.Floor, btnEvent.Button)
		if es.OrderList[idx].HasOrder == false {
			es.OrderList[idx].HasOrder = true
			es.OrderList[idx].VersionNum += 1
		}
	}
}

func (es *ElevatorStatus) canAcceptOrder(newAssignedOrder ButtonEvent) bool {

	return (newAssignedOrder.Floor != -1 && 
		es.DoorOpen && !es.IsObstructed && 
		time.Since(es.TimeOfLastCompletion) > 
		config.OrderTimeOut && 
		es.Pos != -1)
}

func (es *ElevatorStatus) init(chans OrderChannels){
	
	es.IsOnline = true
	es.IsAvailable = true

	chans.RequestElevatorStatus <- true
	initFloor := <-chans.NewFloor
	initDoor := <-chans.DoorOpen

	es.Pos = initFloor
	es.DoorOpen = initDoor
}

func (es *ElevatorStatus)updateOrderCompleted(assignedOrder ButtonEvent) {
	curFloor := es.Pos

	//TODO CLEAN CODE ANNNND LOGIC

	// checks if the elevator has a cab call for the current floor
	if es.CabOrders[curFloor] == true && es.DoorOpen == true && assignedOrder.Floor==curFloor {
		es.CabOrders[curFloor] = false
		es.TimeOfLastCompletion = time.Now()
		return
	}

	orderIdx := floorToOrderListIdx(assignedOrder.Floor, assignedOrder.Button)
	if orderIdx == -1{
		return
	}

	if es.OrderList[orderIdx].HasOrder && es.OrderList[orderIdx].Floor == curFloor {
		es.OrderList[orderIdx].HasOrder = false
		es.OrderList[orderIdx].VersionNum += 1
		es.TimeOfLastCompletion = time.Now()
		// fmt.Println("Hall order compleete:", orderListIdxToFloor(orderIdx))
	}
}

func (es *ElevatorStatus) costFunction(orderFloor int, orderTimeStamp time.Time) int {
	cost := 0
	cost += es.ID
	if !es.IsAvailable || !es.IsOnline {
		cost += unAvailableCost
	}
	if es.IsObstructed {
		cost += ObstructedCost
	}
	// for loop checking if the curElevator has a CabOrders
	for _, cabOrder := range es.CabOrders { //organize into function
		if cabOrder == true {
			cost += unAvailableCost
			break
		}
	}
	cost += waitingTimePenalty(orderTimeStamp)
	cost += oneFloorAwayCost* distanceFromOrder(es.Pos, orderFloor)
	return cost
}