# ptnet-eventstore

ptnet-eventstore is a Petri-Net Event Aggregator for use with the factomd API.

Think of this API as a write-ahead-log for stateful blockchain events atop the factom protocol.

## Status

Currently this is a Proof-Of-Concept only persisting data in memory
- eventually will extend to be a proper eventstore using leveldb + persisting to factom blockchain

### Examples

See [./example.md](./example.md) for sample contract data structures and event sequence.
