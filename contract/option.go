package contract

import (
	"encoding/json"
	"github.com/FactomProject/ptnet-eventstore/gen"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
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
	d := OptionTemplate()

	d.Inputs = []AddressAmountMap{ // array of input depositors
		AddressAmountMap{Address[DEPOSITOR], 1, 0}, // deposit tokens
	}

	d.Outputs = []AddressAmountMap{
		AddressAmountMap{Address[DEPOSITOR], 1, 0}, // withdraw token
		AddressAmountMap{Address[USER1], 1, 0},     // deposit to user1
		AddressAmountMap{Address[USER2], 1, 0},     // deposit to user2
	}

	d.BlockHeight = 60221409 // deadline for halting state

	sig, _ := json.Marshal(d.Variables)
	extids := append(x.Ext(ptnet.OptionV1), sig)
	d.ContractID = x.NewContractID(extids)

	return d
}

// templates are  used to generate the chain digest

// altering this definition will result in a new chain signature
func OptionTemplate() Declaration {
	return Declaration{
		Variables: Variables{
			ContractID:  x.NewContractID(x.Ext(ptnet.OptionV1, "|SALT|")),
			BlockHeight: 0, // deadline for halting state 0 = never
			Inputs:      []AddressAmountMap{},
			Outputs:     []AddressAmountMap{},
		},
		Invariants: Invariants{
			Schema:     ptnet.OptionV1,
			Parameters: gen.OptionV1.Places,
			Capacity:   gen.OptionV1.GetCapacityVector(),
			State:      gen.OptionV1.GetInitialState(),
			Actions:    gen.OptionV1.Transitions,
			Guards: []Condition{ // guard clause restricts actions
				Role(gen.OptionV1, []string{"OPEN"}, 1),
				Role(gen.OptionV1, []string{"OPEN"}, 1),
				Role(gen.OptionV1, []string{"OPEN"}, 1),
			},
			Conditions: []Condition{ // contract conditions specify additional redeem conditions
				Check(gen.OptionV1, []string{"REFUND"}, 1),
				Check(gen.OptionV1, []string{"OUT1"}, 1),
				Check(gen.OptionV1, []string{"OUT2"}, 1),
			},
		},
	}
}
