package contract

import (
	"encoding/json"
	"github.com/FactomProject/ptnet-eventstore/gen"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
)
import . "github.com/FactomProject/ptnet-eventstore/identity"

func SpendContract() Declaration {
	d := SpendTemplate()

	d.Inputs = []AddressAmountMap{ // array of input depositors
		AddressAmountMap{Address[DEPOSITOR], 1, 0}, // deposit tokens
	}

	d.Outputs = []AddressAmountMap{
		AddressAmountMap{Address[USER1], 1, 0}, // deposit to user1
	}

	sig, _ := json.Marshal(d.Variables)
	extids := append(x.Ext(ptnet.OptionV1), sig)
	d.ContractID = x.NewContractID(extids)

	return d
}

func SpendTemplate() Declaration {
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
