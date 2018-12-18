package finite

import (
	"github.com/FactomProject/ptnet-eventstore/contract"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
)

var FiniteChain string = x.NewChainID(x.Ext("Merged", ptnet.OctoeV1, ptnet.OptionV1))
var MetaChain string = x.NewChainID(x.Ext(ptnet.Meta, ptnet.FiniteV1))

// FIXME turn these hardcoded fixtures into factories

func OptionContract() Offer {

	return Offer{
		Declaration: contract.OptionContract(),
		ChainID:     FiniteChain,
	}
}

func TicTacToeContract() Offer {

	return Offer{
		Declaration: contract.TicTacToeContract(),
		ChainID:     FiniteChain,
	}
}

func Registry() Offer {

	return Offer{
		Declaration: contract.RegistryTemplate(),
		ChainID:     MetaChain,
	}
}
