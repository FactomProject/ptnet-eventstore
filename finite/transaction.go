package finite

import (
	"github.com/FactomProject/ptnet-eventstore/contract"
	"github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
)

type Offer struct{
	contract.Declaration
	ChainID string
}
type Execution struct{
	contract.Command
}
type Transaction struct{
	*ptnet.Event
}

func Depositor(o Offer) identity.PrivateKey {
	// FIXME actually return the proper key
	return identity.PrivateKey{}
}

func Executor(x Execution) identity.PrivateKey {
	// FIXME actually return the proper key
	return identity.PrivateKey{}
}

func OfferTransaction(o Offer, privkey identity.PrivateKey) Transaction {
	event, err := contract.Create(o.Declaration, o.ChainID, func(evt *ptnet.Event) error {
		contract.SignEvent(evt, Depositor(o))
		return nil
	})

	if err != nil {
		panic("failed to create contract")
	}

	return Transaction{Event: event}
}

func ExecuteTransaction(x Execution, privkey identity.PrivateKey) (Transaction, error) {
	event, err := contract.Transform(x.Command, func(evt *ptnet.Event) error {
		contract.SignEvent(evt, Executor(x))
		return nil
	})

	return Transaction{Event: event}, err
}
