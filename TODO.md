### WIP

Valid events are queued for persistence to factom blockchain.
Allow existing chains to be validated against state machines replay.
Allow aggregation of contracts to calculate wallet balances by account.

### BACKLOG

- [ ] Define Factom Asset Token FATIP spec for petri-nets
- [ ] Demo how to calculate balance for tokens
- [ ] add import/export of graphml test compatibility with gephi & networkx(python)

### COMPLETE
 
- [x] work with factom connect or factomd to publish transactions
- [x] run w/ fnode sim from factomd project
- [x] refactor to use gopetri instead of ptnet everwhere - missing capacity check without this
- [x] generates state machine code from pnml
- [x] improve string output when printing data structures for examples in this project
- [x] add gzip to payload data
- [x] demo smart contract protocol creation
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

- [ ] target chain - adopt FAT standard for choosing a chainid to store a given event
- [ ] refactor 'finite' transactions to avoid adding redundant sigs since that's what blockchains are for ;)
- [ ] demo smart contract protocol playback/validation
- [ ] 'cache miss' should result when a state record is not found in memory when working w/ actual blockchain storage
- [ ] signing/externalIDs - contracts should output an event as a valid factom entry
- [ ] demo composing multiple state machines in a single contract
- [ ] demo up-converting a contract from v1 -> v2 by extending the length of the statevector
- [ ] use testify/mocks rather than hardcoded literals in contract/definitions.go  finite/definitions.go
- [ ] use testify/mocks rather than literals in identity package 
- [ ] add leveldb to persist events and state to disk
- [ ] add pagination to event stream
- [ ] demo generic petri-net validation
- [ ] read state machine definitions from json files
- [ ] allow extra param to allow API users to specify level of persistence MEM -> Disk -> Blockchain
- [ ] demonstrate a single game of tic-tac-toe by OID using 2 agents + arbiter over blockchain
- [ ] tap into journal messages (for testing) - not needed b/c memdb is used as a write-through cache
