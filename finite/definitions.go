package finite

import (
	"github.com/FactomProject/ptnet-eventstore/contract"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
)


func OptionContract() Offer {

	return Offer{
		Declaration: contract.OptionContract(),
		ChainID: x.NewChainID( x.Ext("Merged", ptnet.OctoeV1, ptnet.OptionV1 )),
	}
}

func TicTacToeContract() Offer {

	return Offer{
		Declaration: contract.TicTacToeContract(),
		ChainID: x.NewChainID( x.Ext("Merged", ptnet.OctoeV1, ptnet.OptionV1 )),
	}
}

func Registry() Offer {

	return Offer{
		Declaration: contract.RegistryTemplate(),
		ChainID: x.NewChainID(x.Ext(ptnet.Meta, ptnet.FiniteV1)),
	}
}
