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

func OfferTransaction(o Offer, privKey identity.PrivateKey) Transaction {
	event, err := contract.Create(o.Declaration, o.ChainID, privKey)

	if err != nil {
		panic("failed to create contract")
	}

	return Transaction{Event: event}
}

func ExecuteTransaction(x Execution, privKey identity.PrivateKey) (Transaction, error) {
	event, err := contract.Transform(x.Command, func(evt *ptnet.Event) error {
		contract.SignEvent(evt, privKey)
		return nil
	})

	return Transaction{Event: event}, err
}
