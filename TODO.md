### WIP

Demonstrate an event driven API
Where on-insert events are applied to state using an in-memory table

Valid events are queued for persistence to disk and to factom blockchain

### BACKLOG
 
- [x] create a POST api for dispatching events
- [x] implement in-memory state machine using go-memdb
- [x] implement schema for states and events
- [x] add ability to list events 
- [x] add value to apply n(action) event
- [x] finish adding a post for dispatch to actually commit an event
- [x] run memdb benchmark > 10k event commits /sec (in memory)
- [x] add ability to query state
- [x] add query for events
- [x] Added first pass data structure
- [x] Q: can we conform to some standards set here: https://github.com/Factom-Asset-Tokens/FAT/blob/master/fatips/101.md
      A: not really but can get close - add new fields to content
- [ ] demo smart contract protocol
- [ ] Define Factom Asset Token FATIP spec for petri-nets
      assets are considered locked when state machine begins and awarded on halt 

### ICEBOX
- [ ] push events to factomd using : git@github.com:FactomProject/factom.git
- [ ] add leveldb to persist events and state to disk
- [ ] add pagination to event stream
- [ ] demo generic petri-net validation
- [ ] read state machine definitions from json files
- [ ] allow extra param to allow API users to specify level of persistence MEM -> Disk -> Blockchain
- [ ] 'cache miss' should result when a state record is not found in memory
- [ ] push events to harmony API
- [ ] demonstrate a single game of tic-tac-toe by OID using 2 agents + arbiter over blockchain
