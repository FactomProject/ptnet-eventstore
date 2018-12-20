package contract

import (
	"encoding/json"
	"github.com/FactomProject/ptnet-eventstore/gen"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
)
import . "github.com/FactomProject/ptnet-eventstore/identity"

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
