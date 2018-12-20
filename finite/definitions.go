package finite

import (
	"github.com/FactomProject/ptnet-eventstore/contract"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
)

var FiniteChain string = x.NewChainID(x.Ext("Merged", ptnet.Spend, ptnet.Tip, ptnet.OctoeV1, ptnet.OptionV1, ptnet.AuctionV1))
var MetaChain string = x.NewChainID(x.Ext(ptnet.Meta, ptnet.FiniteV1))

func SpendContract() Offer {

	return Offer{
		Declaration: contract.SpendContract(),
		ChainID:     FiniteChain,
	}
}


func TipContract() Offer {

	return Offer{
		Declaration: contract.TipContract(),
		ChainID:     FiniteChain,
	}
}

func AuctionContract() Offer {

	return Offer{
		Declaration: contract.AuctionContract(),
		ChainID:     FiniteChain,
	}
}

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
