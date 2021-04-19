

- Når slices returneres, enten:
    1. Lag en predefined const som erstatter uttrykk somm config.NumFloors*2-2
    2. Gjør slicene tomme brackets []


- I initOrderList() står det kommentarer //initializing HallUp og HallDown - hva betyr dette egentlig? Er det initializing HallUp order status eller noe sånn? 

- "Direction" bør kanskje være en egen type? 

- Bra funksjonsnavn på initOrderList og initAllElevatorStatuses

- De fleste stedene det står f.eks. messages.HallOrder bør det ikke være 'messages'. Det kan være den beste metoden er å ikke kalle messages-mappa vår for messages i det hele tatt, men heller noe mer generelt. Grunnen er at det i de fleste tilfeller ikke faktisk er en "message" egentlig.

- Preferansesak kanskje, men burde ikke RunOrders vært øverst i fila siden det er den "overordnede" funksjonen for ordresystemet? 

- Hva betyr "idx"??

- Kan det stemme at en del funksjoner er matriseoperasjoner? F. eks. 
"minimumRowCol" (hva betyr egentlig det? Hva gjør den?). Isåfall bør dette trekkes ut i egen fil siden det er "støttefunksjoner"

- ifXinSliceInt er definitivt en støttefunksjon

- Finner Idx igjen veldig forvirrende

- messages.LampUpdate_message er på en måte et eksempel på en beskrivende type message ettersom det faktisk handler om å sende en "message". 

- I sendOutLightsUpdate er det en Sleep med time.Millisecon * 100. 100 er et "magic number" så bør lage en const over som er _updateLightRate eller noe sånn. 

- hva gjør func sendOutStatus? Kan være tvetydig kanskje? 

- hva betyr contSend? 

- I RunOrders brukes allElevators[config.ID] veldig ofte, men ettersom det ikke er helt klart at 'config.ID' refererer til 'denne heisen', kan det være greit å helt i starten sette thisElevator = allElevators[config.ID]!. Da er det mye mer tydelig gjennom hele koden når vi referer til denne spesifikke heisen. 

- case <-time.After(500*time.Millisecons) slår meg som en råkul måte å gjøre dette på faktisk, men gjerne legg inn en kommentar rett under på hvorfor dette gjøres 


*Kommentar: Vi må nesten ta noen endringer om gangen og så se på det inkrementelt over flere runder, men dette er noen forslag i første omgang*