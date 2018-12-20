package dsl

import (
	"github.com/FactomProject/factom"
	"github.com/FactomProject/ptnet-eventstore/blockchain"
	"github.com/FactomProject/ptnet-eventstore/contract"
	"github.com/FactomProject/ptnet-eventstore/finite"
	"github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
)

// syntactic sugar to group contract-specific operations
type spend struct{}

var Spend = spend{}

func ( _ spend) Contract(id string, ht uint64, in []contract.AddressAmountMap, out []contract.AddressAmountMap) contract.Declaration {
	// TODO: add extra checks to make sure input/output count matches template

	c := contract.SpendContract()
	c.Variables = contract.Variables{
		ContractID: id,
		BlockHeight: ht,
		Inputs: in,
		Outputs: out,
	}
	return c
}

func ( _ spend) Exists( id string) bool {
	return contract.Exists(ptnet.Spend, id)
}

func( _ spend) Offer(b *blockchain.Blockchain, d contract.Declaration, acct *identity.Account) (*factom.Entry, error) {
	return b.Offer(finite.Offer{ Declaration: d, ChainID: b.ChainID}, acct)
}
