package finite

import (
	"github.com/FactomProject/ptnet-eventstore/contract"
)

/*
 Transaction fixtures (currently for testing) are defined in this file
 eventually these structures will be stored on chain and indexed in memory when in use
*/

func OptionContract() Offer {
	return Offer{
		Declaration: contract.OptionContract(),
	}
}

func TicTacToeContract() Offer {
	return Offer{
		Declaration: contract.TicTacToeContract(),
	}
}
