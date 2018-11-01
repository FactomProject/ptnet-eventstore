package contracts

import (
	"encoding/json"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/hashicorp/go-memdb"
)

type Contract struct {
	Schema  string        `json:"schema"`
	Machine ptnet.Machine `json:"state_machine""`
	db      *memdb.MemDB
}

type Transaction struct {
	Address string `json:"address""`
	Amount  uint64 `json:"amount""`
}

type Condition ptnet.Transition

type Declaration struct {
	Inputs      []Transaction               `json:"inputs"`
	Outputs     []Transaction               `json:"outputs"`
	BlockHeight uint64                      `json:"blockheight"`
	Salt        string                      `json:"salt"`
	ContractID  string                      `json:"contractid"`
	Schema      string                      `json:"schema"`
	State       ptnet.StateVector           `json:"state"`
	Actions     map[string]ptnet.Transition `json:"actions"`
	Guards      []Condition                 `json:"guards"` // this enforces contract roles
	Conditions  []Condition                 `json:"conditions"`
}

type ContractState struct {
	ChainID   string      `json:"chainid"`
	LastEntry string      `json:"last_entry"`
	ChainHead string      `json:"chainhead"`
	State     ptnet.State `json:"state"`
}

// TODO add signing
type Command struct {
	ChainID    string `json:"chainid"`
	ContractID string `json:"contractid"`
	Schema     string `json:"schema"`
	Action     string `json:"action"`
	Amount     uint64 `json:"amount"`
	Payload    []byte
}

var Contracts map[string]Contract = map[string]Contract{
	ptnet.OptionV1: Contract{
		Schema:  ptnet.OptionV1,
		Machine: ptnet.StateMachines[ptnet.OptionV1],
		db:      ContractStore(),
	},
	ptnet.OctoeV1: Contract{
		Schema:  ptnet.OctoeV1,
		Machine: ptnet.StateMachines[ptnet.OctoeV1],
		db:      ContractStore(),
	},
}

/*
NOTE: a contract is considered complete if state machine is halted
and one or more redeem conditions are met

// TODO include all the other inputs needed to conform with  https://github.com/Factom-Asset-Tokens/FAT/blob/master/fatips/0.md
*/

// TODO Rewrite to accept a declaration as input
func Create(contract Declaration) (*ptnet.Event, error) {

	payload, _ := json.MarshalIndent(contract, "", "    ")
	println("contract:")
	println(string(payload))

	// FIXME sign this event
	event, err := ptnet.Commit(contract.Schema, contract.ContractID, ptnet.BEGIN, 1, []byte(payload))

	if err != nil {
		txn := Contracts[contract.Schema].db.Txn(true)
		err = txn.Insert(ContractTable, contract)
		if err != nil {
			txn.Commit()
		} else {
			txn.Abort()
		}
	}

	return event, err
}

// FIXME add signing
func Commit(command Command) (*ptnet.Event, error) {
	// TODO: check guards and conditions
	event, err := ptnet.Commit(command.Schema, command.ContractID, command.Action, command.Amount, command.Payload)
	return event, err
}

func IsHalted(contract Declaration) bool {
	txn := ptnet.Txn(contract.Schema, false)
	raw, _ := txn.First(ptnet.StateTable, "id", contract.ContractID)
	if raw  == nil {
		return false
	}

	// TODO also test that if there are open state machine actions the
	// Guard roles should be tested to make sure action is available to contract
	vectorIn := raw.(ptnet.State).Vector
	for _, transition := range contract.Actions {
		_, err := ptnet.VectorAdd(vectorIn, transition, 1)
		if err == nil {
			return false
		}
	}

	return true
}

// TODO: add Redeem Condition Checking
