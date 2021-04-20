## General

Our solution to the elevator challenge is based on an "everyone knows, everyone shares"-approach. All elevators continuously broadcasts their current status while listening to every other elevators status in a peer-to-peer network. 


*** 

## Modules 

- elevControl
- orderDistribution
- network
- messages

## ElevControl

This module consists of an abstract elevator state machine, an elevator controller and driver (provided by the teaching staff). These together implement all the logic needed to represent the elevator as an entity that can be performed operations on. 

The state machine can be thought of as an abstraction of the 'physical' elevator while the controller is used to perform operations on it. More generally, the controller acts like a broker between the 'button-panel'-interface, the state machine and the order module. It redirects button actions to the order system, and the order system sorts out which commands should be performed first. It then finally tells the controller to execute an order.

The controller also provides channels that continuously sends updates on the elevator status (door status, current floor etc.) to external listeners - in this case the order module. The order module then has the option of using this information to predict what actions will be allowed to perform and decide when to perform them.

The framwork applied on the elevator_state_machine code may seem unfamiliar to many, as it is to the most common approach. Therefore, a Readme with more extensive explanations on how it works, code explanations and where it's fetched from is included in the elevController package [here](./Project/elevController/Readme.md).

## OrderDistribution

The orderDistribution module is responnsibile for synchronzing all infomration regarding lights, buttons and order between the elevators. Every elevator stores all the information about itself and the other elevators in the orderDistribution module. The moudle utlizes the assignOrder function to calculate which order every elevator should respond to and makes a decision based on that about which order it should assign to itself. This is poissible due to every node knowing everything about the other nodes at almost all times. 

The premises for the order module are:

- Every elevator broadcasts its own status about its own elevator and information about recent orders. 

- The elevators listen to each other continuously to always be updated about their status, which floor their elevators are on, which version of the orderlist they have etc. 

- Due to every elevator broadcasting their version of the order list all the time, this also serves as a backup system for hall orders. When a node shuts down and restarts, it fetches the order list from one of the other elevators broadcasting it. However, to not lose any cab-orders (which are not broadcasted), it also saves the caborders in a local file that it reads upon initialization. 

- The order module knows "everything" it needs to know about its own elevator, but also about every other elevator, in order to continuously calculate which orders it is best suited to take relative to all other elevators. It will then decide which order to execute on the basis that it 'knows' that it is in the best suited position to do so.  

## Network
The network module is responsible for the communication between the elevators. The communication is done using broadcasting and all elevators listening to the same port on the network. It only consists of the code handed out by the staff and is relatively small. Two channels are created when initializing the program and these are passed to the broadcast and receiver functions respectively. These two channels are allso passed to the orderDistibution module and the module uses these channels to communicate out its own status and receive other elevator's statuses.


## Messages

In order to work around the fact that we have to avoid circular dependecies, while trying to send structs with specific data through channels, we found that using a third-party message system solved our challenge. By defining messages, we can use these to send more complicated datatypes and reduce the need for connecting too many channels between modules. It must be duly noted that GO probably provides a simpler solution, but this is just something we found worked for our application. 