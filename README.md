# ptnet-eventstore

Prototype gRPC evenstore service.

Construct Markov chains backed by CockroachDB or PostgreSQL databases.

Uses Petri-Nets as state machines to validate events before appending to a eventstore.

## Status

[![CircleCI](https://circleci.com/gh/FactomProject/ptnet-eventstore.svg?style=svg)](https://circleci.com/gh/FactomProject/ptnet-eventstore)

Alpha - Seems to function in a development environment

Not tested under load.

## Factom Asset Tokens Compatiblity

The Ultimate aim is to develop a datastore
that maps onto the smart contract platform provided by https://github.com/Factom-Asset-Tokens

Read the FATIP - draft specification [./fatip.md](./fatip.md) # <- TODO rewrite to target latest SmartContract design

## Why use this library?

Using an eventstore that ensures only valid events are stored is a distinct style choice
that can simplify the design of many types of applications where ledger-driven audits are desirable.

Petri-nets are well explored data structures that have mathematically verifiable properties.

States and transitions are computed as a [Vector addition System with State](https://en.wikipedia.org/wiki/Vector_addition_system)
This vector format makes machine learning analysis of event logs very trivial.

This library is compatible with `.pflow` files produced with a [visual editor](http://www.pneditor.org/)
Once a user is familiar with the basic semantics of a Petri-Net, new process flows can be developed rapidly.
