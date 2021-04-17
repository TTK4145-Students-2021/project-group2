package orders

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
func assignOrder(allElevatorSatuses *[numElevators]ElevatorStatus) ButtonEvent_message {
	
	calculateAllCosts(allElevatorSatuses)

	cabOrders := allElevatorSatuses[curID].CabOrders
	curFloor := allElevatorSatuses[curID].Pos
	
	pickedOrder := assignCabOrder(cabOrders, curFloor)
	// if there is an cab order return that as the assignedOrder
	if pickedOrder.Button != messages.UNDEFINED {
		return pickedOrder
	}
	// Pick the order where if there are any possible orders and you have the lowest cost
	pickedOrder = assignHallOrder(allElevatorSatuses)

	return pickedOrder //return the floor that the elevator should go
}

func assignHallOrder(allElevatorSatuses *[numElevators]ElevatorStatus) ButtonEvent_message{

	pickedOrder := ButtonEvent_message{-1, -1} 
	orderList := allElevatorSatuses[curID].OrderList

	// turning the costs into a matrix of [numElevators][orderListIdx]
	costMatrix := make([][numElevators]int, orderListLength)
	for i := range costMatrix {
		costMatrix[i] = orderList[i].Costs
	}
	lowestCost := 500

	/*Findes the lowet cost, if this cost is yours then retrun that order, if not, make that order and the costs for the assigned elevator equal to maxCost,
	 and repeat for numElevators */

	for i := 0; i < numElevators; i++ {
		orderIdx, ID := minumumRowCol(costMatrix)
		cost := costMatrix[orderIdx][ID]
		if ID == curID && 
		cost < lowestCost && 
		orderList[orderIdx].HasOrder {
			pickedOrder.Floor = orderListIdxToFloor(orderIdx)
			pickedOrder.Button = ButtonType_msg(orderList[orderIdx].Direction)
			return pickedOrder
		} else {
			for i := range orderList {
				costMatrix[i][ID] = maxCost
			}
			for i := 0; i < numElevators; i++ {
				costMatrix[orderIdx][i] = maxCost
			}
		}
	}
	return pickedOrder
}

func assignCabOrder(cabOrders [numFloors]bool, curFloor int) ButtonEvent_message {
	//default pickedOrder
	pickedOrder := ButtonEvent_message{-1, -1}

	shortestCabDistance := numFloors + 1
	for orderFloor,curOrder := range cabOrders{
		if curOrder == true {
			curCabDistance := distanceFromOrder(orderFloor, curFloor)
			if curCabDistance < shortestCabDistance {
				pickedOrder.Floor = orderFloor
				pickedOrder.Button = messages.BT_Cab
			}
		}
	}
	return pickedOrder
}



func calculateAllCosts(allElevatorSatuses *[numElevators]ElevatorStatus)  {
	orderList := &allElevatorSatuses[curID].OrderList
	for idx,_ := range orderList{
		curOrder := &orderList[idx]
		if orderList[idx].HasOrder{
			for _, curElevator := range allElevatorSatuses{
				cost := costFunction(curElevator, curOrder.Floor, curOrder.TimeStamp)
				curOrder.Costs[curElevator.ID] = cost
			}
		}else {
			for idx, _ := range curOrder.Costs {
				curOrder.Costs[idx] = maxCost
			}
		}
		
	}
}

func costFunction(curElevator ElevatorStatus, orderFloor int, orderTimeStamp time.Time) int {
	cost := 0
	cost += curElevator.ID
	if !curElevator.IsAvailable || !curElevator.IsOnline {
		cost += unAvailableCost
	}
	if curElevator.IsObstructed {
		cost += ObstructedCost
	}
	// for loop checking if the curElevator has a CabOrders
	for _, cabOrder := range curElevator.CabOrders { //organize into function
		if cabOrder == true {
			cost += unAvailableCost
			break
		}
	}
	cost += waitingTimePenalty(orderTimeStamp)
	cost += oneFloorAwayCost* distanceFromOrder(curElevator.Pos, orderFloor)
	return cost
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

