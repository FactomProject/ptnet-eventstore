package contract

import (
	"encoding/json"
	"github.com/FactomProject/ptnet-eventstore/gen"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
)
import . "github.com/FactomProject/ptnet-eventstore/identity"

func AuctionContract() Declaration {
	d := AuctionTemplate()

	d.Inputs = []AddressAmountMap{ // array of input depositors
		AddressAmountMap{Address[DEPOSITOR], 1, ptnet.Coin}, // deposit tokens
	}

	d.Outputs = []AddressAmountMap{
		AddressAmountMap{Address[DEPOSITOR], 1, ptnet.Coin}, // withdraw token
		AddressAmountMap{Address[USER1], 1, ptnet.Coin},     // deposit to user1
	}

	d.BlockHeight = 60221409 // deadline for halting state

	sig, _ := json.Marshal(d.Variables)
	extids := append(x.Ext(ptnet.AuctionV1), sig)
	d.ContractID = x.NewContractID(extids)

	return d
}

func AuctionTemplate() Declaration {
	return Declaration{
		Variables: Variables{
			ContractID:  x.NewContractID(x.Ext(ptnet.AuctionV1, "|SALT|")),
			BlockHeight: 0, // deadline for halting state 0 = never
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
