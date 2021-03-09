package elevio

/*
import (
	"fmt"
	"sync"
)

//USE SLICES FOR LIST!!!
//https://blog.golang.org/slices-intro

var _lsmtx sync.Mutex

type ExecutionOrder struct {
	OrderType   ButtonType
	TargetFloor int
}

type ExecutionList []ExecutionOrder

func CreateExecutionList(len int, cap int) []ExecutionOrder {
	var list []ExecutionOrder = make([]ExecutionOrder, len, cap)
	return list
}

func InsertOrder(orderList *ExecutionList, order ExecutionOrder) int {
	_lsmtx.Lock()
	defer _lsmtx.Unlock()

	elevState = GetElevState() // Use channel??
	currentFloor = getFloor()  // Use channel??

	Here I want the elevator to wait until it is in an IDLE state -> not moving
	When it is idle, it should perform the calculations and prioritize orders that are
	in its MOVING DIRECTION (going up or going down)

	switch elevState {
	case UNAVAILABLE:
		fmt.Printf("Elevator is unavailable. Can't give order!")
		return 1
	case GOING_UP:

	}
	//Check if possible?
	//TODO: Write the algortihm that sorts the orders
	return 0
}
*/
