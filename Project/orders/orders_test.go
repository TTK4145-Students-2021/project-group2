package orders

import (
	"fmt"
	"testing"
)

func TestOrders(t *testing.T) {

	dummyStatus = initElevatorStatus()
	fmt.Println(dummyStatus)
}
