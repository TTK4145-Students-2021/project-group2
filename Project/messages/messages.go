package messages

type ButtonType_message int

const (
	BT_HallUp   ButtonType_message = 0
	BT_HallDown ButtonType_message = 1
	BT_Cab      ButtonType_message = 2
)

type ButtonEvent_message struct {
	Floor  int
	Button ButtonType_message
}
