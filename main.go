package main

import (
	"encoding/json"
	"github.com/FactomProject/ptnet-eventstore/contracts"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
)

var emptyPayload []byte

// execute test transactions
func main() {
	emptyPayload, _ = json.Marshal(map[string]string{})
	runTicTacToe() // a game of tic-tac-toe
	runOption()    // a simple contract with choice of output addresses
}

func marshal(x interface{}) []byte {
	data, _ := json.MarshalIndent(x, "", "    ")
	return data
}

func runOption() {
	var event *ptnet.Event
	contract := contracts.OptionContract()
	event, _ = contracts.Create(contract)
	println("event:")
	println(string(marshal(event)))

	var actionQueue []string = []string{
		"OPT_0",
		"HALT",
	}

	if false == contracts.IsHalted(contract) {
		println("Contract is not halted")
	}

	var err error

	for _, action := range actionQueue {
		event, err = contracts.Commit(contracts.Command{ // FIXME add signing
			ChainID:    "|ChainID|",
			ContractID: "|ContractID|",
			Schema:     ptnet.OptionV1,
			Action:     action, // state machine action
			Amount:     1, // triggers input action 'n' times
			Payload:    emptyPayload, // arbitrary data optionally included
		})

		if err != nil {
			panic(err)
		}
		println("event:")
		println(string(marshal(event)))
	}

	if contracts.IsHalted(contract) {
		print("Contract is halted")
	}
}

func runTicTacToe() {
	var event *ptnet.Event
	contract := contracts.TicTacToeContract()
	event, _ = contracts.Create(contract)
	println("event:")
	println(string(marshal(event)))

	var actionQueue []string = []string{
		"X11",
		"O01",
		"X00",
		"O02",
		"X22",
		"WINX",
	}

	if false == contracts.IsHalted(contract) {
		println("Contract is not halted")
	}

	var err error
	for _, action := range actionQueue {
		event, err = contracts.Commit(contracts.Command{ // FIXME add signing
			ChainID:    "|ChainID|",
			ContractID: "|ContractID|", // contract uuid
			Schema:     ptnet.OctoeV1, // state machine version
			Action:     action, // state machine action
			Amount:     1, // triggers input action 'n' times
			Payload:    emptyPayload, // arbitrary data optionally included
		})

		if err != nil {
			panic(err)
		}
		println("event:")
		println(string(marshal(event)))
	}

	if contracts.IsHalted(contract) {
		print("Contract is halted")
	}

}
