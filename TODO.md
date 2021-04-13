# TODO:
- Working with 4 elevators 4 floors. Not with 9 floors 4 elevators
- Logic for `isOnline` must be implemented

## Errors
- When elevator arrives at floor it can only handle one order.

## Brukerfeil
- Heis 1 får hall order, på vei får den en cab order. Heis 1 vil fullføre hall orderen den har begynt på, men vil fullføre cab order etterpå. Når cab order kommer inn vil heis 2 begynne på hall orderen heis 1 fikk i utgangspunktet.

## Refactoring
- Rename config.ID to something like "thisID"
- Change for loops to run on `for i, element := range list` and not `for i=0; i<0; i++`