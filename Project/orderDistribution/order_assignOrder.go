package orderDistribution

// this module chooses wich order the elevator should execute
import (
	"time"
	"../messages"
)

// this is an overview of the different costs utilized by the cost function
const maxCost = 1000000
const ObstructedCost = 100
const oneFloorAwayCost = 10
const unAvailableCost = 4000
const timePenaltyCost = -15

// runs through the costFunction for all true orders
func assignOrder(allElevatorStatuses *AllElevatorStatuses) ButtonEvent {
	
	allElevatorStatuses[_thisID].OrderList.calculateAllCosts(*allElevatorStatuses)

	cabOrders := allElevatorStatuses[_thisID].CabOrders
	curFloor := allElevatorStatuses[_thisID].Pos
	
	pickedOrder := assignCabOrder(cabOrders, curFloor)
	// if there is an cab order return that as the assignedOrder
	if pickedOrder.Button != messages.UNDEFINED {
		return pickedOrder
	}
	// Pick the order where if there are any possible orders and you have the lowest cost
	pickedOrder = assignHallOrder(*allElevatorStatuses)

	return pickedOrder //return the floor that the elevator should go
}

func assignHallOrder(allElevatorSatuses AllElevatorStatuses) ButtonEvent{

	pickedOrder := ButtonEvent{-1, -1} 
	orderList := allElevatorSatuses[_thisID].OrderList

	// turning the costs into a matrix of [_numElevators][orderListIdx]
	costMatrix := make([][_numElevators]int, _orderListLength)
	for i := range costMatrix {
		costMatrix[i] = orderList[i].Costs
	}
	lowestCost := 500

	/*Findes the lowet cost, if this cost is yours then retrun that order, if not, make that order and the costs for the assigned elevator equal to maxCost,
	 and repeat for _numElevators */

	for i := 0; i < _numElevators; i++ {
		orderIdx, ID := minumumRowCol(costMatrix)
		cost := costMatrix[orderIdx][ID]
		if ID == _thisID && 
		cost < lowestCost && 
		orderList[orderIdx].HasOrder {
			pickedOrder.Floor = orderListIdxToFloor(orderIdx)
			pickedOrder.Button = ButtonType(orderList[orderIdx].Direction)
			return pickedOrder
		} else {
			for i := range orderList {
				costMatrix[i][ID] = maxCost
			}
			for i := 0; i < _numElevators; i++ {
				costMatrix[orderIdx][i] = maxCost
			}
		}
	}
	return pickedOrder
}

func assignCabOrder(cabOrders [_numFloors]bool, curFloor int) ButtonEvent {
	//default pickedOrder
	pickedOrder := ButtonEvent{-1, -1}

	shortestDistance := _numFloors + 1
	for orderFloor, hasOrder := range cabOrders{
		if hasOrder == true {
			curDistance := distanceFromOrder(curFloor, orderFloor)
			if curDistance < shortestDistance {
				pickedOrder.Floor = orderFloor
				pickedOrder.Button = messages.BT_Cab
				shortestDistance = curDistance
			}
		}
	}
	return pickedOrder
}

func waitingTimePenalty(timeStamp time.Time) int {
	extraCost := 0
	timeWaited := time.Now().Sub(timeStamp) //organize into function
	//larger than 1 min
	if timeWaited > time.Duration(time.Minute*1) {
		extraCost += timePenaltyCost
	}
	//larger than 3 min
	if timeWaited > time.Duration(time.Minute*3) {
		extraCost += timePenaltyCost*3
	}
	return extraCost
}

func distanceFromOrder(elevatorFloor int, goingToFloor int) int {
	floorsAway := elevatorFloor - goingToFloor
	if floorsAway < 0 { // takeing abs of floorsAway
		floorsAway = -floorsAway
	}
	return floorsAway
}

