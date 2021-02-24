package elevio

import "fmt"

func TakeOrder(button ButtonType, floor int, value bool) {

	fmt.Printf("%v\n", button)
	fmt.Printf("%v\n", floor)

	switch button {
	case BT_HallUp:
		fmt.Printf("Button press up")

	case BT_HallDown:
		fmt.Printf("Button press down")

	case BT_Cab:
		fmt.Printf("Button press cab")
	}
}
