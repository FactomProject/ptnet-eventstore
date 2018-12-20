package contract

import (
	"github.com/FactomProject/ptnet-eventstore/gen"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
)
import . "github.com/FactomProject/ptnet-eventstore/identity"

// altering this definition will result in a new chain signature
// REVIEW extend new Contract based on inputs
func TicTacToeContract() Declaration {
	d := TicTacToeTemplate()

	d.Inputs = []AddressAmountMap{ // array of input depositors
		AddressAmountMap{Address[DEPOSITOR], 1, 0},
	}

	d.Outputs = []AddressAmountMap{
		AddressAmountMap{Address[DEPOSITOR], 1, 0},
		AddressAmountMap{Address[PLAYERX], 1, 0},
		AddressAmountMap{Address[PLAYERO], 1, 0},
	}

	d.BlockHeight = 60221409 // deadline for halting state

	return d
}

func TicTacToeTemplate() Declaration {

	return Declaration{ // array of inputs/outputs also referenced by guards and conditions
		Variables: Variables{
			ContractID:  x.NewContractID(x.Ext(ptnet.OctoeV1)),
			BlockHeight: 0,                    // deadline for halting state 0 = never
			Inputs:      []AddressAmountMap{}, // array of input depositors
			Outputs:     []AddressAmountMap{}, // array of possible redeemers

		},
		Invariants: Invariants{
			Schema:     ptnet.OctoeV1, // versioned contract schema
			Parameters: gen.OctoeV1.Places,
			Capacity:   gen.OctoeV1.GetCapacityVector(), // capacity for each place
			State:      gen.OctoeV1.GetInitialState(),   // initial state
			Actions:    gen.OctoeV1.Transitions,         // state machine defined transitions
			Guards: []Condition{ // guard clause restricts actions
				Role(gen.OctoeV1, []string{}, 1),         // depositor is unrestricted
				Role(gen.OctoeV1, []string{"turn_x"}, 1), // players must
				Role(gen.OctoeV1, []string{"turn_o"}, 1), // take turns
			},
			Conditions: []Condition{ // contract conditions specify additional redeem conditions
				Check(gen.OctoeV1, []string{"REFUND"}, 1), // no winner take 1 token prize back
				Check(gen.OctoeV1, []string{"OUT_X"}, 1),  // pay player X 1 token
				Check(gen.OctoeV1, []string{"OUT_O"}, 1),  // pay player O 1 token
			},
		},
	}
}
