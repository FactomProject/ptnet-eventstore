| FATIP | Title           | Status | Category | Author                               | Created    |
| ----- | --------------- | ------ | -------- | ------------------------------------ | ---------- |
| -     | Finite Protocol | WIP    | Protocol | Matt York \<<matt.york@factom.com>\> | 12-20-2018 |



# Summary

This document describes a protocol for publishing state-driven contracts and transactions to the blockchain.
Rather than trying to define a specific token standard, this is an attempt to provide an executible
specification that can be contained wholly within the Factom blockchain.

# Motivation

Much time and effort as been spent around tokenized value,
this is an attempt to view tokens first and foremost as a meaningful expression of data.

The intent is to introduce a standard for publishing information that is:
1. Verifiable while only having access to on-chain data
2. Constructed using well explored math and computer science theory
3. Expressive without being Turing complete
4. Extensible allowing a systematic method to define expected behavior

# Specification

## Terms

* spec: may refer this document or a specific chain definition composed of contracts
* blockchain schema: a specific datastucture that conforms to this document
* registry: a chain that records chain definitions
* state vector: the current state of a given contract
* state machine: place-transition system allowing for deterministic or non-deterministic behavior
* Petri-Net: visual form of the above type of state machine
* action: the human friendly name of a transition, also part of a state machine
* witness sequence: an ordered set of transaction entries that 'execute' a contract on-chain

### Reference

* State Machines are written to the chain in this form: https://en.wikipedia.org/wiki/Vector_addition_system
* States can be expressed visually as a Petri-Net: https://en.wikipedia.org/wiki/Petri_net
* Contract behavior is modeled using Dual Petri Nets: https://en.wikipedia.org/wiki/Dualistic_Petri_nets
* Formal theory provides an avenue to verify correctness of a contract: https://en.wikipedia.org/wiki/Curry%E2%80%93Howard_correspondence

## Chain Registry 

A Registry allows for initial creation & publication of a new blockchain spec.

    A Registry is a chain that tracks the lifecycle of a blockchain schema
   
Using this spec a given chain schemata can also be updated or disabled by adding entries to a Registry chain.

NOTE: This is not meant to be a global registry, though it could be used to cooperate with other chain specs.
 
## Token ID

Tokens are defined per-chain, and a given chain specify up to `max(uint8)`

## Addresses

Addresses are FCT human-readable addresses converted to hash form.

Essentially this just drops the prefix + base58 encoding.

## Issuance

Using this spec it is possible to declare a chain that has no tokens, or a chain that has many.

The default token is `Color(uint8(0))` and by design it can be arbitrarily created or destroyed during any transaction.

A user may define smart contracts to issue new types of tokens - we adopt the term 'Color' to differentiate between token types.

## State Machines

    State machines are the basic building blocks of a contract.

If a state machine is used in a transaction that does not specify inputs or outputs,
it effectively functions as a method to broadcast structured oracle data rather than an exchange of value.

## Contracts

Contracts are built upon state machines, and define rules and actions a transaction sequence.

A contract is open when:

      A "Offer" entry is created on a supporting chain with a unique Contract ID
    
A contract is considered halted:

    If current state has no valid actions available
    
A contract is considered invalid:

    predefined blockcheight has been exceeded without the contract in the halted state
    
An inductive definition for contract state:
    
    state is first initialized when an contract is offered
    
    Current state is calculated by adding the transaction state vector
    to the output state from the most recent valid transition 
    
    A state transition that contains a negative output value is considered invalid.
   
For a contract that has defined input and output addresses, it is required for the public key signature to match the input address.
* This could be extended to provide for multi-sig inputs, but currently that is not defined behavior.

    No more transactions an be applied to a contract when it is halted
    
    If a contract does not halt before its set blockheight limit, it is considered invalid. 
    
### Roles and Checks

A contract can provide additional Guards and Conditions that are compared to the State Vector.
These basic rules can be leveraged to exchange, transfer, create and destroy tokens.

Roles and Checks are evaluated with the same rules as State Machine actions, but do not alter state.

    A Role or Check is valid if it results in a valid output vector when added to current contract state
    
#### Guards

Guards enforce Roles by comparing with input state of an open contact.

    Guards are state vectors used to check pre-conditions of the input contract state.

This allows for a contract author to limit interaction with the contract to a specific set of users.

#### Conditions

Conditions add checks against the output state of a halted contract

    Conditions are state vectors used to determine which Output address should recieve tokens
    
## Transactions

In this context 'contract' could mean 'Smart Contract' - an execution dealing with the exchange of tokenized value.

Or, it could mean 'Api Contract' - a more broad term referring to the expected behavior of a stateful data API.

    Transactions are individual executions of a defined contract.

Contracts are not required to allow input or output tokens, and so may never halt by design, but for those that do:

    A transaction terminates when it's underlying state machine is in a halted state.
    
One part of the Transaction data structure includes the Input & Output state around the Action being executed:

    A specified state vector input must match previous vector output,
    
    If a transaction entry arrives out of order, execute using First-In-First-Out priority blockhain to specify order
    
    A valid transaction entry is applied to calculate the new state if it results in valid output
    
### Actions

    An action is a human readable string that maps to a vector addtion operation in a contract. 

This standard reduces all rules, conditions, and calculations to vector math that can be easily audited.

    The state of a given contract is represented as a vector of unsigned integers.
    
An Action is executed when it is applied to calculate the current state of an open contract.

### Witness Sequence

The blockchain is a public witness that determines the order in which input transactions should be executed.
By design, the first transaction that results in a hating state completes the contract.

An additional consequence of this choice is that contract IDs must not be reused.

Providing a witness sequence of related on chain entries is the only way to verify a transaction was executed.

### Entry Validation

There are 4 types of entries in this system:

* Blockchain Specification - added to the Registry Chain
* Contract Offer - a transaction constructed to transfer tokens to other users.
* Contract Execution - a transaction constructed to act upon a valid offer.
* Commit - a transaction that does not require further interaction to be valid


     All entries must be signed.
     
     The signature on an offer Entry must match the input address of the Offer.
     
     The signature of an execution Entry must match one of the output addresses
 

### Content

Only Offer and Spec entries require attached of content - the json definition of the contract variables.

Then intent behind most of this protocol type metadata in ExtIDs is to allow additional supporting data to be appended as part of a valid transaction.

# Implementation

Runs against an in-memory Factomd simulation.

Proof-of-Concept (POC): https://github.com/FactomProject/ptnet-eventstore

# Copyright

Copyright and related rights waived via
[CC0](https://creativecommons.org/publicdomain/zero/1.0/).
