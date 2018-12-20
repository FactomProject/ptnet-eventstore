package ptnet

import (
	"github.com/FactomProject/ptnet-eventstore/gen"
)

/*
declare state machines
REVIEW: consider loading definitions from file on demand

These state machines are re-usable components that can be
referenced and integrated w/ smart contract definitions

This may mean that part of this protocol is to publish versioned
state machines to a distinct revision-control chain on Factom
*/

const Meta string = "Meta" // registry schema
const Tip string = "Tip"
const Spend string = "Spend"
const FiniteV1 string = "FiniteV1"   // meta protocol for publishing blockchain definitions
const OptionV1 string = "OptionV1"   // version contract definitions by using schema name
const AuctionV1 string = "AuctionV1" // protocol for running an Auction between tokesn
const OctoeV1 string = "OctoeV1"     // this allows for future mechanism to 'upconvert' v1 -> v2

var StateMachines map[string]Machine = map[string]Machine{
	"counter": counterMachine,
	OctoeV1:   octoeMachine,
	OptionV1:  optionMachine,
	FiniteV1:  finiteMachine,
	Meta:      finiteMachine,
}

var optionMachine Machine = Machine{
	StateMachine: gen.OptionV1.StateMachine(),
	db:           EventStore(),
}

var counterMachine Machine = Machine{
	StateMachine: gen.CounterV1.StateMachine(),
	db:           EventStore(),
}

var octoeMachine Machine = Machine{
	StateMachine: gen.OctoeV1.StateMachine(),
	db:           EventStore(),
}

var finiteMachine Machine = Machine{
	StateMachine: gen.FiniteV1.StateMachine(),
	db:           EventStore(),
}
