| FATIP | Title           | Status | Category | Author                               | Created    |
| ----- | --------------- | ------ | -------- | ------------------------------------ | ---------- |
| -     | Finite Protocol | WIP    | Protocol | Matt York \<<matt.york@factom.com>\> | 12-20-2018 |


* [Summary](#summary)
* [Motivation](#motivation)
* [Specification](#specification) - introducing the protocol
   * [Terms](#terms) - short dictionary of terms used in this spec
   * [Reference](#reference) - wikipedia links to CompSci topics
   * [Behavior](#state-machine--contract-behavior) - Primitive Operations
     * [Halting State](#halting-state)
     * [Guards and Conditions](#additional-inhibitors-and-checks)
   * Finite Protocol - key aspects
     * [Registry Chain](#chain-registry)
     * [TokenID](#token-id)
     * [Addresses](#addresses)
     * [Issuance](#issuance)
     * [State Machines](#state-machines)
     * [Contracts](#contracts)
     * [Roles and Checks](#roles-and-checks)
     * [Guards](#guards)
     * [Conditions](#conditions)
     * [Transactions](#transactions)
     * [Actions](#actions)
     * [Witness Sequence](#witness-sequence)
     * [Entry Validation](#entry-validation)
     * [Payload](#content)
* [Implementation](#implementation) - Data Structures & examples
  * [Entries](#entry-format)
  * [Offers](#offer-format)
  * [Registry](#registry-chain)
  * Use Cases
    * [Spend](#use-case-token-spend) - simple token transfer
    * [Auction](#use-case-auction) - conditional token transfer
    * [Tip](#use-case-tracking-tips-with-colored-tokens) - accounting with colored tokens
  

# Summary

STATUS: will revisit to account for WASM smart contracts

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
4. Extensible and Updatable using a systematic method of expressing behavior

# Specification

## Terms

* spec: may refer this document or a specific chain definition composed of contracts
* blockchain schema: a specific data stucture that conforms to this document
* registry: a chain that records chain definitions
* state vector: the current state of a given contract
* state machine: place-transition system allowing for deterministic or non-deterministic behavior
* Petri-Net: visual form of a place-transition style state machine
* action: the human friendly name of a transition, also part of a state machine
* witness sequence: an ordered set of transaction entries that are said to 'execute' a contract on-chain

### Reference

* State Machines are written to the chain in this form: https://en.wikipedia.org/wiki/Vector_addition_system
* States can be expressed visually as a Petri-Net: https://en.wikipedia.org/wiki/Petri_net
* Contract behavior is modeled using Dual Petri Nets: https://en.wikipedia.org/wiki/Dualistic_Petri_nets
* Formal theory provides an avenue to verify correctness of a contract: https://en.wikipedia.org/wiki/Curry%E2%80%93Howard_correspondence

### State Machine / Contract Behavior

The type of state machines being used here can be visually explained using Petri-Nets,
but are expressed on-chain as [vector addition systems](https://en.wikipedia.org/wiki/Vector_addition_system).

This basic behavior defines the 'executable' part of this specification.

#### Halting State

State machines can be easily analyzed to determine it it is 'Halted' or not -- simply put
if there are no valid moves to be made this state machine is said to be in a halted state.

As part of this standard any state machine based contract execution must be halted to indicate it has completed.
As a rule, if a given contract instance is not halted before the specified block hight it is considered invalid.

#### Additional Inhibitors and Checks

On top of a standard Petri-net model - this design includes special guard and condition vectors that run 'on top' of an elementary state machine.
These added checks allow for more complex rule sets to be composed from primitive operations.

*Guards*: add a precondition check that is used to control what user is allowed to post the next entry to the contract.
*Conditions*: output check called conditions control the release of tokens from the locked contract.

All these rules combined serve as a generalized state protocol that can be used to extend the overall behavior of an token
system without having to encode specific rules into a client-side/wallet application.

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
    
    Any state transition that results in a vector containing a negative scalar output is considered invalid.
   
For a contract that has defined input and output addresses, it is required for the public key signature to match the input address.
* This could be extended to provide for multi-sig inputs, but currently that is not defined behavior.

    No more transactions an be applied to a contract when it is halted
    
    If a contract does not halt before its set blockheight limit, it is considered invalid. 
    
### Roles and Checks

A contract can provide additional Guards and Conditions that are compared to the State Vector.
These basic rules can be leveraged to exchange, transfer, create and destroy tokens.

Roles and Checks are evaluated with the same rules as State Machine actions, but do not alter state.

    A Role or Check is valid if it results in a valid output vector when added to current contract state
    
* Example Use Cases
    * [Spend](#use-case-token-spend) - unconditional transfer immediately halted
    * [Auction](#use-case-auction)  - multi-step conditional transfer
    * [Tip](#use-case-tracking-tips-with-colored-tokens) - unconditional transfer with multiple token colors
    
#### Guards

Guards enforce Roles by comparing with input state of an open contact.

    Guards are state vectors used to check pre-conditions of the input contract state.

This allows for a contract author to limit interaction with the contract to a specific set of users.

A guard is a vector transformation that works just like a Transition/Action, with the exception that the state of the contract is not mutated.

Basic logic here is - when a guard is applied to contract state, it is considered valid if the output vector contains no negative numbers.

#### Conditions

Conditions add checks against the output state of a halted contract

    Conditions are state vectors used to determine which Output address should recieve tokens
    
Condition checks function just as guards do - except they are evaluated only after the contract is halted.
    
## Transactions

In this context 'contract' could mean 'Smart Contract' - an execution dealing with the exchange of tokenized value.

Or, it could mean 'Api Contract' - a more broad term referring to the expected behavior of a stateful data API.

    Transactions are individual executions of a defined contract.

Contracts are not required to allow input or output tokens, and so may never halt by design, but for those that do:

    A transaction terminates when it's underlying state machine is in a halted state.
    
One part of the Transaction data structure optionally includes the Input & Output state around the Action being executed:

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

Providing a witness sequence of on-chain entries is the only way to verify a transaction was faithfully executed.

NOTE: Wallet like functionality (as in fatd) should be developed to calculate the current state of any given address
that may contain tokens defined by the finite protocol.

### Entry Validation

There are 4 types of entries in this system:

* Blockchain Specification - added to the Registry Chain
* Contract Offer - a transaction constructed to transfer tokens to other users.
* Contract Execution - a transaction constructed to act upon a valid offer.
* Commit - a transaction that does not require further interaction to be valid

     All entries must be signed.
     
     The signature on an offer Entry must match the input address of the Offer.
     
     The signature of an execution Entry must match one of the output addresses
 
 * [Entry Format](#entry-format)
 * [Entry Validation](#entry-validation)

### Content

Only Offer and Spec entries require attached of content - the json definition of the contract variables.

Then intent behind most of this protocol type metadata in ExtIDs is to allow additional supporting data to be appended as part of a valid transaction.

 * [Entry Content](#entry-format)

# Implementation

Code in this repo runs against an in-memory Factomd simulation.
Some example data structures are provided below.

NOTE: this is a draft - any part of this spec that can be amended to conform
to other FAT specs should be amended.

Proof-of-Concept (POC): https://github.com/FactomProject/ptnet-eventstore

### Entry Format
FIXME: add auction example

### Offer Format

FIXME: add auction example

```json
EntryHash: 387e8921ac629e8e35bfba28be6637097596fce628181363f340d8f9ee166b09
ChainID: cd4b3c137f353a9682c9e017f7eae6b205b5478bfd2794a7447bde96b686a7dc
ExtID: Spend // contract schema
ExtID: EXEC // action
ExtID: 395d3030f7f471b2c7f8783efa57b884181493abfc2c52d5f5842b72efb78526 // contractID
ExtID: 3  // multiplier - trigger transfer 3x
ExtID: [0,0] //input state (optional)
ExtID: [3,0] // output state (optional)
ExtID: c04c1cd15ba4581ee677232cd7a2726f9bf73a8f6f5bc7d64970074d76e64f9f
ExtID: da5dd59d59cf320b7a89cd0e474e2e7787903daff93d7372a0a23a9d3733594ccaaace55d6cdf4f03c74481f8408aa55f4159c302e423c5b725964264ae9ad08
Content: // payload - required for 'EXEC' - but other subsequent 'action' entries can optionally include additional data or nothing
{"contractid":"395d3030f7f471b2c7f8783efa57b884181493abfc2c52d5f5842b72efb78526","blockheight":0,
"inputs":[{"address":"Ayb3xCM4BVP3Cw0ob+bBf6KbFc9Ojm5uAp7BzvO98U0=","amount":1,"color":1}],
"outputs":[{"address":"cRQW+RoDrfgKLhpaL4n7r0WlxICDOhbcEJC5Sxg4yMA=","amount":1,"color":1}]}
```

### Registry Chain

Here's an example entry for a new chain containing use cases referenced in this document.

```json
{
    "chainid": "cd4b3c137f353a9682c9e017f7eae6b205b5478bfd2794a7447bde96b686a7dc",
    "contracts": {
        "FiniteV1": {
            "schema": "FiniteV1",
            "template": {
                "actions": {
                    "DISABLE": [ -1, 0, 1 ],
                    "ENABLE": [ 1, 0, -1 ],
                    "EXEC": [ 0, 1, 0 ]
                },
                "blockheight": 0,
                "capacity": [ 0, 0, 0 ],
                "conditions": null,
                "contractid": "46050f96e9197f332c94effdbf921d7ec79cca55f0dae0f28aa0d677aa578b42",
                "guards": null,
                "inputs": null,
                "outputs": null,
                "parameters": {
                    "active": {
                        "capacity": 0,
                        "initial": 1,
                        "offset": 0
                    },
                    "edits": {
                        "capacity": 0,
                        "initial": 0,
                        "offset": 1
                    },
                    "inactive": {
                        "capacity": 0,
                        "initial": 0,
                        "offset": 2
                    }
                },
                "schema": "FiniteV1",
                "state": [ 1, 0, 0 ]
            }
        },
        "Spend": {
            "schema": "Spend",
            "template": {
                "actions": {
                    "EXEC": [ 1, 0 ]
                },
                "blockheight": 0,
                "capacity": [ 0, 1 ],
                "conditions": [
                    [
                        -1,
                        0
                    ]
                ],
                "contractid": "72fdd6e42898ac04f037e161e0c28f47ffd5675718e4f986715cf9232fac16a9",
                "guards": [
                    [ 0, -1 ]
                ],
                "inputs": [],
                "outputs": [],
                "parameters": {
                    "HALTED": {
                        "capacity": 1,
                        "initial": 0,
                        "offset": 1
                    },
                    "PAYMENT": {
                        "capacity": 0,
                        "initial": 0,
                        "offset": 0
                    }
                },
                "schema": "Spend",
                "state": [ 0, 0 ]
            }
        },
        "AuctionV1": {
            "schema": "AuctionV1",
            "template": {
                "actions": {
                    "EXEC": [ 1, 0 ]
                },
                "blockheight": 0,
                "capacity": [ 0, 1 ],
                "conditions": [
                    [ -1, 0 ]
                ],
                "contractid": "72fdd6e42898ac04f037e161e0c28f47ffd5675718e4f986715cf9232fac16a9",
                "guards": [
                    [ 0, -1 ]
                ],
                "inputs": [],
                "outputs": [],
                "parameters": {
                    "HALTED": {
                        "capacity": 1,
                        "initial": 0,
                        "offset": 1
                    },
                    "PAYMENT": {
                        "capacity": 0,
                        "initial": 0,
                        "offset": 0
                    }
                },
                "schema": "Spend",
                "state": [ 0, 0 ]
            }
        },
        "Tip": {
            "schema": "Tip",
            "template": {
                "actions": {
                    "EXEC": [ 1, 0 ]
                },
                "blockheight": 0,
                "capacity": [ 0, 1 ],
                "conditions": [
                    [ -1, 0 ]
                ],
                "contractid": "72fdd6e42898ac04f037e161e0c28f47ffd5675718e4f986715cf9232fac16a9",
                "guards": [
                    [ 0, -1 ]
                ],
                "inputs": [],
                "outputs": [],
                "parameters": {
                    "HALTED": {
                        "capacity": 1,
                        "initial": 0,
                        "offset": 1
                    },
                    "PAYMENT": {
                        "capacity": 0,
                        "initial": 0,
                        "offset": 0
                    }
                },
                "schema": "Spend",
                "state": [ 0, 0 ]
            }
        }
    },
    "extids": [
        "U3BlbmRDaGFpbg=="
    ],
    "tokens": [
        {
            "Color": 0
        },
        {
            "Color": 1
        },
        {
            "Color": 2
        },
        {
            "Color": 3
        }
    ]
}
```

### Use Case: Token Spend
source: https://github.com/FactomProject/ptnet-eventstore/blob/master/contract/spend.go

#### spend state
![spend state machine without checks][spend3]

The essential state of a spend transaction is a single issued offer that is created
in an 'already-halted' state.

#### spend contract offer
![running spend state machine][spend3 running]

To conform with the halting state rules of the protocol - an additional place called 'HALTED' is added.

Notice above that if there were a token in the Halted position - the transaction would still be considered to be running.
#### spend contract execution
![halted spend state machine][spend3 halted]

By initializing the contract offer in the already-halted state, we achive a single-entry contract that requires no interaction by the receiving address to redeem.

[spend3]: http://factomstatus.com/ptnet-eventstore/image/spend3.png "spend state machine without checks in halted position"
[spend3 running]: http://factomstatus.com/ptnet-eventstore/image/spend3-running.png "spend state machine with checks in running position"
[spend3 halted]: http://factomstatus.com/ptnet-eventstore/image/spend3-halted.png "spend state machine with checks in halted position"

#### spend state machine code

The Petri-Nets above are used to generate the code below.

```go
var Spend PetriNet = PetriNet{
	Places: map[string]Place { 
		"HALTED": Place{
				Initial: 0,
				Offset: 1,
				Capacity: 1,
		},
		"PAYMENT": Place{
				Initial: 0,
				Offset: 0,
				Capacity: 0,
		},
	},
	Transitions: map[Action]Transition { 
		"EXEC": Transition{ 1,0 },
	},
}
```

#### spend state machine contract

When a contract is instantiated - variables are provided to create a transaction.

```go
func SpendContract() Declaration {
	d := SpendTemplate()

	d.Inputs = []AddressAmountMap{ // array of input depositors
		AddressAmountMap{Address[USER1], 1, 0}, // deposit tokens
	}

	d.Outputs = []AddressAmountMap{
		AddressAmountMap{Address[USER2], 1, 0}, // deposit to user1
	}

	sig, _ := json.Marshal(d.Variables)
	extids := append(x.Ext(ptnet.Spend), sig)
	d.ContractID = x.NewContractID(extids)

	return d
}

```

#### spend state machine template

The underlying state machine, Guards, and Conditions are all consider 'invariants' and are ultimately defined be the chain's registry entry.

```go
func SpendTemplate() Declaration {
	return Declaration{
		Variables: Variables{
			ContractID:  x.NewContractID(x.Ext(ptnet.Spend, "|SALT|")),
			BlockHeight: 0, // deadline for halting state 0 = never
			Inputs:      []AddressAmountMap{},
			Outputs:     []AddressAmountMap{},
		},
		Invariants: Invariants{
			Schema:     ptnet.Spend,
			Parameters: gen.Spend.Places,
			Capacity:   gen.Spend.GetCapacityVector(),
			State:      gen.Spend.GetInitialState(),
			Actions:    gen.Spend.Transitions,
			Guards: []Condition{ // guard clause restricts actions
				Role(gen.Spend, []string{"HALTED"}, 1),
			},
			Conditions: []Condition{ // contract conditions specify additional redeem conditions
				Check(gen.Spend, []string{"PAYMENT"}, 1),
			},
		},
	}
}
```

### Use Case: Auction
source: https://github.com/FactomProject/ptnet-eventstore/blob/master/contract/auction.go

[auction3]: http://factomstatus.com/ptnet-eventstore/image/auction3.png "auction state machine without checks in halted position"
[auction3 running]: http://factomstatus.com/ptnet-eventstore/image/auction3-running.png "auction state machine with checks in running position"
[auction3 halted]: http://factomstatus.com/ptnet-eventstore/image/auction3-halted.png "auction state machine with checks in halted position"

#### auction state
![auction state machine without checks][auction3]

This state machine allows for a user to increase a bid amount in an effort to get another user to accept the contract offer.

#### auction contract offer
![running auction state machine][auction3 running]

Notice the `?2` label represents a check - this is different from the Spend transaction because it is created in the 'running' state.

This indicates that additional inputs are expected to drive the contract to a halted state.

#### auction contract execution
![halted auction state machine][auction3 halted]

Above we see the completed contract, notice the `?0` label is marking a Condition indicating that the bid has been accepted by the user looking to sell tokens.

#### auction state machine code

The state machine below is generated from the Petri-Net definitions above.
```go
var AuctionV1 PetriNet = PetriNet{
	Places: map[string]Place { 
		"ACCEPTED": Place{
				Initial: 0,
				Offset: 0,
				Capacity: 1,
		},
		"NEW": Place{
				Initial: 1,
				Offset: 1,
				Capacity: 1,
		},
		"OPEN": Place{
				Initial: 0,
				Offset: 2,
				Capacity: 1,
		},
		"PRICE": Place{
				Initial: 0,
				Offset: 3,
				Capacity: 0,
		},
		"REJECTED": Place{
				Initial: 0,
				Offset: 4,
				Capacity: 1,
		},
	},
	Transitions: map[Action]Transition { 
		"BID": Transition{ 0,0,0,1,0 },
		"EXEC": Transition{ 0,-1,1,0,0 },
		"HALT": Transition{ 0,0,-1,0,1 },
		"SOLD": Transition{ 1,0,-1,0,0 },
	},
}
```

#### auction state machine contract

The state machine is reference by the contract when an offer is issued.
Notice that in the example below - the 'Color' of the token is provided as part of input.

```go
func AuctionContract() Declaration {
	d := AuctionTemplate()

	d.Inputs = []AddressAmountMap{ // array of input depositors
		AddressAmountMap{Address[DEPOSITOR], 1, ptnet.Coin}, // deposit tokens quantity
	}

	d.Outputs = []AddressAmountMap{
		AddressAmountMap{Address[DEPOSITOR], 1, ptnet.Coin}, // withdraw token deposit back to original owner (no sale)
		AddressAmountMap{Address[USER1], 1, ptnet.Coin},     // deposit to user1
	}

	d.BlockHeight = 60221409 // deadline for halting state

	sig, _ := json.Marshal(d.Variables)
	extids := append(x.Ext(ptnet.AuctionV1), sig)
	d.ContractID = x.NewContractID(extids)

	return d
}
```

Currently this example isn't very interesting as the user is bidding using ptnet.Coins to win ptent.Coins.
More likely this use case would be used to swap one token type for another.

#### auction state machine template

The underlying state machine, Guards, and Conditions are all consider 'invariants' and are ultimately defined be the chain's registry entry.

```go
func AuctionTemplate() Declaration {
	return Declaration{
		Variables: Variables{
			ContractID:  x.NewContractID(x.Ext(ptnet.AuctionV1, "|SALT|")),
			BlockHeight: 0, // deadline for halting state 0 = always/never
			Inputs:      []AddressAmountMap{},
			Outputs:     []AddressAmountMap{},
		},
		Invariants: Invariants{
			Schema:     ptnet.AuctionV1,
			Parameters: gen.AuctionV1.Places,
			Capacity:   gen.AuctionV1.GetCapacityVector(),
			State:      gen.AuctionV1.GetInitialState(),
			Actions:    gen.AuctionV1.Transitions,
			Guards: []Condition{ // guard clause restricts actions
				Role(gen.AuctionV1, []string{"OPEN", "PRICE"}, 1),
				Role(gen.AuctionV1, []string{"OPEN"}, 1),
			},
			Conditions: []Condition{ // contract conditions specify additional redeem conditions
				Check(gen.AuctionV1, []string{"REJECTED"}, 1),
				Check(gen.AuctionV1, []string{"ACCEPTED"}, 1),
			},
		},
	}
}
```

### Use Case: tracking tips with Colored Tokens
source: https://github.com/FactomProject/ptnet-eventstore/blob/master/contract/tip.go

[tip3]: http://factomstatus.com/ptnet-eventstore/image/tip3.png "tip state machine without checks in halted position"
[tip3 running]: http://factomstatus.com/ptnet-eventstore/image/tip3-running.png "tip state machine with checks in running position"
[tip3 halted]: http://factomstatus.com/ptnet-eventstore/image/tip3-halted.png "tip state machine with checks in halted position"

#### tip state machine
![tip state machine without checks][tip3]

See above: the underyling state works just as in the [spend state](#spend-state)

This is a 1-way transaction that does not require additional interaction beyond the contract offer.

#### tip state machine offer
![running tip state machine][tip3 running]

Comparing the diagrams above and below - notice how the `?0` label indicates that the guard is testing for halting state.
#### tip state machine execution
![halted tip state machine][tip3 halted]

When halted the `?1` and `?2` Conditions are then evalutated to determine which address has unlocked the deposited tokens.

#### tip state machine code
The code below is gerated from the Petri-Net shown above.

```go
var Tip PetriNet = PetriNet{
	Places: map[string]Place { 
		"ANTIKARMA": Place{
				Initial: 0,
				Offset: 0,
				Capacity: 0,
		},
		"CHARITY": Place{
				Initial: 0,
				Offset: 1,
				Capacity: 0,
		},
		"HALTED": Place{
				Initial: 0,
				Offset: 2,
				Capacity: 0,
		},
		"KARMA": Place{
				Initial: 0,
				Offset: 3,
				Capacity: 0,
		},
		"PAYMENT": Place{
				Initial: 0,
				Offset: 4,
				Capacity: 0,
		},
	},
	Transitions: map[Action]Transition { 
		"TIP": Transition{ 0,0,0,1,1 },
		"WARN": Transition{ 1,1,0,0,0 },
	},
}
```


#### tip state machine contract

The state machine code is referenced by the contract offer when it is created.

```go
func TipContract() Declaration {
	d := OptionTemplate()

	d.Inputs = []AddressAmountMap{ // array of input depositors
		AddressAmountMap{Address[DEPOSITOR], 1, ptnet.Coin}, // deposit tokens
	}

	d.Outputs = []AddressAmountMap{ // User1 withdraws all tokens
		AddressAmountMap{Address[USER1], 1, ptnet.Coin},     // deposit to user1
		AddressAmountMap{Address[USER2], 1, ptnet.Coin},     // deposit to another charity (anyone but user1)
		AddressAmountMap{Address[USER1], 1, ptnet.Karma},     // deposit to user1
		AddressAmountMap{Address[USER1], 1, ptnet.AntiKarma},     // deposit to user1
	}

	d.BlockHeight = 60221409 // deadline for halting state

	sig, _ := json.Marshal(d.Variables)
	extids := append(x.Ext(ptnet.Tip), sig)
	d.ContractID = x.NewContractID(extids)

	return d
}
```

#### tip state machine template
The underlying state machine, Guards, and Conditions are all consider 'invariants' and are ultimately defined be the chain's registry entry.

```go
func TipTemplate() Declaration {
	return Declaration{
		Variables: Variables{
			ContractID:  x.NewContractID(x.Ext(ptnet.Tip, "|SALT|")),
			BlockHeight: 0, // deadline for halting state 0 = never
			Inputs:      []AddressAmountMap{},
			Outputs:     []AddressAmountMap{},
		},
		Invariants: Invariants{
			Schema:     ptnet.Tip,
			Parameters: gen.Tip.Places,
			Capacity:   gen.Tip.GetCapacityVector(),
			State:      gen.Tip.GetInitialState(),
			Actions:    gen.Tip.Transitions,
			Guards: []Condition{ // guard clause restricts actions
				Role(gen.Tip, []string{"HALTED"}, 1), // created in a halted state
			},
			Conditions: []Condition{ // contract conditions specify additional redeem conditions
				Check(gen.Tip, []string{"PAYMENT"}, 1),
				Check(gen.Tip, []string{"CHARITY"}, 1),
				Check(gen.Tip, []string{"KARMA"}, 1),
				Check(gen.Tip, []string{"ANTIKARMA"}, 1),
			},
		},
	}
}
```

# Copyright

Copyright and related rights waived via
[CC0](https://creativecommons.org/publicdomain/zero/1.0/).
