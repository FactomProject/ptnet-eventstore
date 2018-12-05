package contract

import (
	"github.com/FactomProject/ptnet-eventstore/gen"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
)
import . "github.com/FactomProject/ptnet-eventstore/identity"

/*
 contract fixtures (currently for testing) are defined in this file
 eventually these structures will be stored on chain and indexed in memory when in use
*/

// define a deposit contract that makes an offer to 2 identities
// provides a way to revoke payment by original sender
// will only redeem once
// tokens should be thought of as "pay-to-script" - locked until state machine is halted
func OptionContract() Declaration {
	p := gen.OptionV1

	return Declaration{ // array of inputs also referenced by guards and conditions
		Inputs: []AddressAmountMap{ // array of input depositors
			AddressAmountMap{Address[DEPOSITOR], 1}, // deposit tokens
		},
		Outputs: []AddressAmountMap{
			AddressAmountMap{Address[DEPOSITOR], 1}, // withdraw token
			AddressAmountMap{Address[USER1], 1},     // deposit to user1
			AddressAmountMap{Address[USER2], 1},     // deposit to user2
		},
		BlockHeight: 60221409,       // deadline for halting state
		Salt:        "|RANDOM|",     // added random salt
		ContractID:  "|OptionContractID|",   // unique ID for this contract instance
		Schema:      ptnet.OptionV1,         // versioned contract schema
		Capacity:    p.GetCapacityVector(), // state machine initial state
		State:       p.GetInitialState(),   // state machine initial state
		Actions:     p.Transitions,  // state machine defined transitions
		Guards: []Condition{ // guard clause restricts actions
			Role(p,[]string{"OPEN"}, 1),
			Role(p,[]string{"OPEN"},1),
			Role(p,[]string{"OPEN"},1),
		},
		Conditions: []Condition{ // contract conditions specify additional redeem conditions
			Check(p,[]string{"REFUND"},1),
			Check(p,[]string{"OUT1"},1),
			Check(p,[]string{"OUT2"},1),
		},
	}
}

func TicTacToeContract() Declaration {
	p := gen.OctoeV1

	return Declaration{ // array of inputs/outputs also referenced by guards and conditions
		Inputs: []AddressAmountMap{ // array of input depositors
			AddressAmountMap{Address[DEPOSITOR], 1},
		},
		Outputs: []AddressAmountMap{ // array of possible redeemers
			AddressAmountMap{Address[DEPOSITOR], 1},
			AddressAmountMap{Address[PLAYERX], 1},
			AddressAmountMap{Address[PLAYERO], 1},
		},
		BlockHeight: 60221409,      // deadline for halting state
		Salt:        "|RANDOM|",    // added random salt
		ContractID:  "|OctoeContractID|",   // unique ID for this contract instance
		Schema:      "OctoeV1",    // versioned contract schema
		Capacity:    p.GetCapacityVector(), // capacity for each place
		State:       p.GetInitialState(), // initial state
		Actions:     p.Transitions, // state machine defined transitions
		Guards: []Condition{ // guard clause restricts actions
			Role(p, []string{}, 1), // depositor is unrestricted
			Role(p, []string{"turn_x"}, 1), // players must
			Role(p, []string{"turn_o"}, 1), // take turns
		},
		Conditions: []Condition{ // contract conditions specify additional redeem conditions
			Check(p, []string{"REFUND"}, 1), // no winner take 1 token prize back
			Check(p, []string{"OUT_X"}, 1), // pay player X 1 token
			Check(p, []string{"OUT_O"}, 1),  // pay player O 1 token
		},
	}
}
