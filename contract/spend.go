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
