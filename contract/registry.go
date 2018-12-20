package contract

import (
	"github.com/FactomProject/ptnet-eventstore/gen"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
)

func RegistryTemplate() Declaration {

	return Declaration{ // array of inputs/outputs also referenced by guards and conditions
		Variables: Variables{
			ContractID:  x.NewContractID(x.Ext(ptnet.Meta)),
			BlockHeight: 0, // deadline for halting state 0 = never
			Inputs:      nil,
			Outputs:     nil,
		},
		Invariants: Invariants{
			Schema:     ptnet.FiniteV1,
			Parameters: gen.FiniteV1.Places,
			Capacity:   gen.FiniteV1.GetCapacityVector(),
			State:      gen.FiniteV1.GetInitialState(),
			Actions:    gen.FiniteV1.Transitions,
			Guards:     nil,
			Conditions: nil,
		},
	}
}
