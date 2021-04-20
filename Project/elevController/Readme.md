# Elevator package 

This Readme is included in order to clarify any additional information regarding the elevator-package and its content files. 


## Elevator_controller

The ElevatorController works like a broker between the button-panel interface, the order-sorting system and the actions performed on the elevator. Any button-event is redirected to whatever listener that needs this information, and any request to perform an operation on the elevator is performed by the controller (and fails if it is not possible to perform). 

## Elevator_fsm

If we concider the physical elevator as a state machine, this submodule is an abstract implementation of its states, logic and possible actions. To implement the state pattern, a perticular pattern framework has been used (see: https://venilnoronha.io/a-simple-state-machine-framework-in-go). This framework might seem unfamiliar to some programmers as it is not the most 'common approach'. Therefore, we will put some effort in this section to clarify how it works. 

Note from the author:

> *Programmers may or may not have a tendency to prefer their own approaches and pre-assess new architectural paradigms with skepticism. I want to leave a careful reminder that even though something seems unfamiliar and weird at first, it should not be judged neither inferior or superior too early. The approach used might in some cases expand the reades perspective on programming tactics and unveil some new and interesting possibilities. That was at least my own experience implementing this particular framework.*

### Framework explanation 

There are three main elements to the framework: States, Events and Actions. The purpose of the framework is to handle the separation and transition of states in a well defined manner. 

First, we define states and events in the form of EventType and StateType strings. In the following example from the code we define the state "Uninitialized" and the event "Move". 

    
    const (
	Uninitialized StateType = "Uninitialized"

    Move EventType = "Move"
    )

Next, we define 'which events are possible in a which state'. We can't 'Move' in the 'Unitialized' state, so that becomes an illegal transition. We define the legal transitions like this: 

    return &ElevatorMachine{
		States: States{
			Uninitialized: State{
				Events: Events{
					Initialize: Initializing,
				},
			},
            Initializing: State{
				Action: &InitAction{},
				Events: Events{
					SetInitialized: AtFloorDoorsOpen,
				},
			},

In this case we return a new ElevatorMachine with two states: 'Uninitialized' and 'Initializing'. The states have predefined Events and these now become the legal events that can be performed in that State. E.g. 'SetInitialized" represents an event. Then 'AtFloorDoorsOpen' is the State that event transitions into!

In this particular section of the code, we therefore establish our elevator states, the events that are possible in that state and the new state that follows a particular event. 

There is one last piece of the puzzle - how actions are implemented. In the code example above, we see the line `Action: &InitAction{}`. This actually represents the action that is to be performed whenever we transition into the 'Initializing' state. The flow of the program is therefore decided by these Action-calls (which are basically functions). 

Lastly, we have to be able to send events to the machine externally. E.g. when we want it to move to a target floor. We then use the function `SendEvent` that takes in the type of event you want to perform and a context. The context is a wrapper for relevant information that needs to be passed to the event. E.g. when you want to tell the elevator to 'move' you also have to tell it 'which floor'. The SendEvent-function works like a loop that processes a chain of events and state transitions until a 'NoOp'-event is returned, indicating there are no more operations. 

I hope some readers interested in cool pattern implementations found this section useful and enlightening. Personally, I found this framework fun to work with, easy to use once understood and extremely powerful. I would definetly recommend trying it out yourself. Credits to the original author. 

Source: https://venilnoronha.io/a-simple-state-machine-framework-in-go).

***