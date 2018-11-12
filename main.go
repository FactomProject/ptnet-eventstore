package main

import (
	"encoding/json"
	"github.com/FactomProject/ptnet-eventstore/contract"
	. "github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
)

var emptyPayload []byte

func marshal(x interface{}) []byte {
	data, _ := json.MarshalIndent(x, "", "    ")
	return data
}

func runOption() {
	var event *ptnet.Event
	c := contract.OptionContract()
	event, _ = contract.CreateAndSign(c, contract.CHAIN_ID, DEPOSITOR_SECRET)
	println("event:")
	println(string(marshal(event)))

	var actionQueue []string = []string{
		"OPT_1",
		"HALT",
	}

	if false == contract.IsHalted(c) {
		println("Contract is not halted")
	}

	var err error
	var key string

	for _, action := range actionQueue {

		switch action {
		case "OPT_1":
			key = USER1
		case "OPT_2":
			key = USER2
		default:
			key = DEPOSITOR
		}

		event, err = contract.Commit(contract.Command{ // FIXME add signing
			ChainID:    "|ChainID|",
			ContractID: "|ContractID|",
			Schema:     ptnet.OptionV1,
			Action:     action,       // state machine action
			Amount:     1,            // triggers input action 'n' times
			Payload:    emptyPayload, // arbitrary data optionally included
			Pubkey:     key,
		}, Identity[key])

		if err != nil {
			panic(err)
		}
		println("event:")
		println(string(marshal(event)))
	}

	if contract.IsHalted(c) {
		print("Contract is halted")
	}
}

func runTicTacToe() {
	var event *ptnet.Event
	c := contract.TicTacToeContract()
	event, _ = contract.CreateAndSign(c, contract.CHAIN_ID, DEPOSITOR_SECRET)
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

	if false == contract.IsHalted(c) {
		println("Contract is not halted")
	}

	var err error
	var key string

	for _, action := range actionQueue {
		switch action[0] {
		case 'X':
			key = PLAYERX
		case 'O':
			key = PLAYERO
		default:
			key = DEPOSITOR
		}

		event, err = contract.Commit(contract.Command{ // FIXME add signing
			ChainID:    "|ChainID|",
			ContractID: "|ContractID|", // contract uuid
			Schema:     ptnet.OctoeV1,  // state machine version
			Action:     action,         // state machine action
			Amount:     1,              // triggers input action 'n' times
			Payload:    emptyPayload,   // arbitrary data optionally included
			Pubkey:     key,
		}, Identity[key])

		if err != nil {
			panic(err)
		}
		println("event:")
		println(string(marshal(event)))
	}

	if contract.IsHalted(c) {
		print("Contract is halted")
	}

}

// execute test transactions
func main() {
	emptyPayload, _ = json.Marshal(map[string]string{})
	runTicTacToe() // a game of tic-tac-toe
	runOption()    // a simple contract with choice of output addresses
}
