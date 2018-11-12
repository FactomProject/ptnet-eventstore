### WIP

Valid events are queued for persistence to disk and to factom blockchain.
Allow existing chains to be validated against state machines replay.
Allow aggregation of contracts to calculate wallet balances by account.

### BACKLOG

- [ ] model w/ protobuff
- [ ] demo smart contract protocol creation and playback
- [ ] target chain - adopt FAT standard for choosing a chainid to store a given event
- [ ] signing/externalIDs - contracts should output an event as a valid factom entry
- [ ] demo composing multiple state machines in a single contract
- [ ] demo up-converting v1 -> v2 by extending the length of the statevector

### COMPLETE
 
- [x] create a POST api for dispatching events
- [x] in-memory state machine using go-memdb
- [x] schema for states and events
- [x] ability to list events 
- [x] value to apply n(action) event
- [x] URL route for dispatch to actually commit an event
- [x] memdb benchmark > 10k event commits /sec (in memory)
- [x] ability to query state
- [x] query for events
- [x] first pass data structure
- [x] investigate Q: can we conform to some standards set here: https://github.com/Factom-Asset-Tokens/FAT/blob/master/fatips/101.md
      A: not really but can get close - add new fields to content

### ICEBOX

- [ ] Define Factom Asset Token FATIP spec for petri-nets
      assets are considered locked when state machine begins and awarded on halt 
- [ ] push events to factomd using : git@github.com:FactomProject/factom.git
- [ ] add leveldb to persist events and state to disk
- [ ] add pagination to event stream
- [ ] demo generic petri-net validation
- [ ] read state machine definitions from json files
- [ ] allow extra param to allow API users to specify level of persistence MEM -> Disk -> Blockchain
- [ ] 'cache miss' should result when a state record is not found in memory
- [ ] push events to harmony API
- [ ] demonstrate a single game of tic-tac-toe by OID using 2 agents + arbiter over blockchain

