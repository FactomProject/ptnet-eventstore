# ptnet-eventstore

ptnet-eventstore is a Petri-Net Event Aggregator for use with the Factomd API.

Think of this API as a write-ahead-log for stateful blockchain events atop the Factom protocol.

## Status

Currently this is a Proof-Of-Concept only persisting data in memory
- eventually will extend to be a proper eventstore using leveldb + persisting to Factom blockchain

### Examples

Jupyter notebooks in ./example demo contract execution using Golang.
You can execute these examples iteratively using https://github.com/gopherdata/gophernotes

See [./examples/TicTacToe.ipynb](./examples/TicTacToe.ipynb) for sample contract data structures and event sequence.
