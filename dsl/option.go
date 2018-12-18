package dsl

import (
	"github.com/FactomProject/factom"
	"github.com/FactomProject/ptnet-eventstore/blockchain"
	"github.com/FactomProject/ptnet-eventstore/contract"
	"github.com/FactomProject/ptnet-eventstore/finite"
	"github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/stackdump/gopetri/statemachine"
)

// syntactic sugar to group contract-specific operations
type option  struct {}
var Option = option{}

func ( _ option) Exists( id string) bool {
	return contract.Exists(ptnet.OptionV1, id)
}

func ( _ option) Contract(id string, ht uint64, in []contract.AddressAmountMap, out []contract.AddressAmountMap) contract.Declaration {
	c := contract.OptionTemplate()
	c.Variables = contract.Variables{
		ContractID: id,
		BlockHeight: ht,
		Inputs: in,
		Outputs: out,
	}
	return c
}

func( _ option) Offer(b *blockchain.Blockchain, d contract.Declaration, acct *identity.Account) (*factom.Entry, error) {
	return b.Offer(finite.Offer{ Declaration: d, ChainID: b.ChainID}, acct)
}

func( _ option) Execute(b *blockchain.Blockchain, c contract.Declaration, a statemachine.Action, acct *identity.Account) (*factom.Entry, error) {
	return b.Execute(contract.Command{
		b.ChainID,
		c.ContractID,
		c.Schema,
		string(a),
		1, // KLUDGE specific to option contract - no action can be triggered more than 1x
		nil,
		acct.GetPublicKey(),
	}, acct)
}
