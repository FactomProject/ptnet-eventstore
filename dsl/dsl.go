// Package dsl - provides syntactic sugar to interacting with smart contracts
package dsl

import (
	"github.com/FactomProject/ptnet-eventstore/blockchain"
	"github.com/FactomProject/ptnet-eventstore/contract"
	"github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
	"github.com/stackdump/gopetri/statemachine"
)

var ValidContract = blockchain.ValidContract
var ValidSignature = blockchain.ValidSignature
var DeployRegistry = blockchain.DeployRegistry
var NewBlockchain = blockchain.NewBlockchain


var FROM = TO
var Token0 = ptnet.Default
var IsHalted = contract.IsHalted

func ID (extIDs ...string) string {
	ext := [][]byte{}

	for _, id := range extIDs {
		ext = append(ext, x.Encode(id))
	}
	return x.NewContractID(ext)
}

func TO (addr []byte, amt uint64, token uint8) contract.AddressAmountMap {
	return contract.AddressAmountMap{ Address: addr, Amount: amt, Token: uint8(token) }
}

func Inputs (deposit ...contract.AddressAmountMap) []contract.AddressAmountMap {
	return deposit
}

func Outputs (withdraw ...contract.AddressAmountMap) []contract.AddressAmountMap {
	return withdraw
}

func Address(name string) []byte {
	return identity.Address[name]
}

func Action(a string) statemachine.Action {
	return statemachine.Action(a)
}

func CanRedeem(c contract.Declaration, acct *identity.Account) bool {
	return contract.CanRedeem(c, acct.GetPublicKey())
}
