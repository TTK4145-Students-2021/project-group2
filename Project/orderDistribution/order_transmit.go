package orderDistribution

import (
	"../messages"
	"time"
)
/*
*=============================================================================
 * @Description:
 * Contains all functions used for sneding information from the orderDsitribution module
 * to the elevtorController and to the rest of the network
 *
 * 
/*=============================================================================
*/

const _lampPollRate = 20 * time.Millisecond

func sendOutStatus(channel chan<- ElevatorStatus, status ElevatorStatus) {
	status.LastAlive = time.Now()
	channel <- status
}

func sendingElevatorToFloor(channel chan<- int, goToFloor int) {
	channel <- goToFloor
}

func sendOutLightsUpdate(sendLampUpdate chan<- messages.LampUpdate, status *ElevatorStatus) {
	for {
		for _, order := range status.OrderList {
			light := messages.LampUpdate{
				Floor:  order.Floor,
				Button: messages.ButtonType(order.Direction), // = 0 up og 1 down
				Turn:   order.HasOrder,
			}
			sendLampUpdate <- light
		}
		for i, order := range status.CabOrders {

			light := messages.LampUpdate{
				Floor:  i,
				Button: messages.BT_Cab, // = 0 up og 1 down
				Turn:   order,
			}
			sendLampUpdate <- light
		}
		time.Sleep(_lampPollRate)
	}
}

func broadcast(channel chan<- ElevatorStatus, status *ElevatorStatus) {
	for {
		go sendOutStatus(channel, *status)
		time.Sleep(_bcastRate)
	}
}