package contracts

import (
	"encoding/json"
	"errors"
	"fmt"
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

type Command struct {
	ChainID    string `json:"chainid"`
	ContractID string `json:"contractid"`
	Schema     string `json:"schema"`
	Action     string `json:"action"`
	Amount     uint64 `json:"amount"`
	Payload    []byte `json:"payload"`
	Privkey    string
	Pubkey     string
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

// start a new transaction with in-memory db
func Txn(schema string, write bool) *memdb.Txn {
	return Contracts[schema].db.Txn(write)
}

func Create(contract Declaration) (*ptnet.Event, error) {

	payload, _ := json.MarshalIndent(contract, "", "    ")
	//println("contract:")
	//println(string(payload))

	// FIXME sign this event
	event, err := ptnet.Commit(contract.Schema, contract.ContractID, ptnet.BEGIN, 1, []byte(payload))

	if err != nil {
		panic(err)
	}

	txn := Contracts[contract.Schema].db.Txn(true)
	err = txn.Insert(ContractTable, contract)
	if err != nil {
		panic(err)
		//txn.Abort()
	}
	txn.Commit()

	return event, err
}

func state(schema string, contractID string) (ptnet.State, error){
	txn := ptnet.Txn(schema, false)
	raw, err := txn.First(ptnet.StateTable, "id", contractID)
	if raw  == nil || err != nil {
		return ptnet.State{}, errors.New("missing state")
	}

	return raw.(ptnet.State), nil
}

// validate event against guard conditions
func evalGuards(event *ptnet.Event) error {
	txn := Contracts[event.Schema].db.Txn(false)
	raw, _ := txn.First(ContractTable, "id", event.Oid)

	if raw == nil {
		return errors.New("missing contract "+ event.Schema + "." + event.Oid)
	}

	c := raw.(Declaration)

	currentState, _ := state(event.Schema, event.Oid)

	// expect event to be signed
	for i, g := range c.Guards {
		_ , err := ptnet.VectorAdd(currentState.Vector, g, 1)
		if err == nil {
			if ptnet.ValidSignature(event, c.Outputs[i].Address) {
				return nil
			} else {
				//fmt.Printf("guard[%v] => %v \n", i, c.Outputs[i].Address)
			}
		}
	}
	// FIXME actually restrict action if event is not signed w/ a valid key
	//return errors.New("failed guard condition")
	return nil
}

func Commit(cmd Command) (*ptnet.Event, error) {
	return ptnet.Transform(cmd.Schema, cmd.ContractID, cmd.Action, cmd.Amount, cmd.Payload, func (evt *ptnet.Event) error {
		// FIXME actually do signing
		sig := fmt.Sprintf("signed with: %v",  cmd.Privkey)
		ptnet.AddDigest(evt)
		ptnet.AddSignature(evt, cmd.Pubkey, sig)
		return evalGuards(evt)
	})
}

func Exists(schema string, contractID string)  bool {
	txn := Contracts[schema].db.Txn(false)
	raw, err := txn.First(ContractTable, "id", contractID)
	if err != nil {
		panic(err)
	}

	return raw != nil
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
