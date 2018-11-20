package contract

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
	"github.com/hashicorp/go-memdb"
)

// KLUDGE mock values used for testing
var CHAIN_ID string = "|ChainID|"

type Contract struct {
	Schema  string        `json:"schema"`
	Machine ptnet.Machine `json:"state_machine""`
	db      *memdb.MemDB
}

type AddressAmountMap struct {
	Address []byte `json:"address""`
	Amount  uint64 `json:"amount""`
}

type Condition ptnet.Transition

type Declaration struct {
	Inputs      []AddressAmountMap          `json:"inputs"`
	Outputs     []AddressAmountMap          `json:"outputs"`
	BlockHeight uint64                      `json:"blockheight"`
	Salt        string                      `json:"salt"`
	ContractID  string                      `json:"contractid"`
	Schema      string                      `json:"schema"`
	State       ptnet.StateVector           `json:"state"`
	Actions     map[string]ptnet.Transition `json:"actions"`
	Guards      []Condition                 `json:"guards"`     // enforces contract roles
	Conditions  []Condition                 `json:"conditions"` // enforce redeem conditions
}

type State struct {
	ChainID   string      `json:"chainid"`
	LastEntry string      `json:"last_entry"`
	ChainHead string      `json:"chainhead"`
	State     ptnet.State `json:"state"`
}

type Command struct {
	ChainID    string
	ContractID string
	Schema     string
	Action     string
	Amount     uint64
	Payload    []byte
	Pubkey     identity.PublicKey //        compare w/ factom identity standard
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

func SignEvent(event *ptnet.Event, privKey identity.PrivateKey) error {
	sig := x.NewSignature(privKey[:], event.GetDigest())
	pubKey := identity.PublicKey{}
	copy(pubKey[:], x.PrivateKeyToPub(privKey[:]))
	ptnet.AddSignature(event, pubKey, sig.Bytes())
	return nil
}

func Create(contract Declaration, chainID string, privkey identity.PrivateKey) (*ptnet.Event, error) {
	return create(contract, chainID, func(evt *ptnet.Event) error {
		return SignEvent(evt, privkey)
	})
}

func create(contract Declaration, chainID string, signfunc func(*ptnet.Event) error) (*ptnet.Event, error) {

	payload, _ := json.Marshal(contract)
	//println("contract:")
	//println(string(payload))

	pubkey := identity.PublicKey{}

	event, err := Transform(
		Command{
			ChainID:    chainID, // test values
			ContractID: contract.ContractID,
			Schema:     contract.Schema,
			Action:     ptnet.BEGIN,     // state machine action
			Amount:     1,               // triggers input action 'n' times
			Payload:    []byte(payload), // arbitrary data optionally included
			Pubkey:     pubkey,          // REVIEW: will there always be a single input?
		}, signfunc)

	if err != nil {
		panic(err)
	}

	txn := Contracts[contract.Schema].db.Txn(true)
	err = txn.Insert(ContractTable, contract)
	if err != nil {
		panic(err)
	}
	txn.Commit()

	return event, err
}

func state(schema string, contractID string) (ptnet.State, error) {
	txn := ptnet.Txn(schema, false)
	raw, err := txn.First(ptnet.StateTable, "id", contractID)
	if raw == nil || err != nil {
		return ptnet.State{}, errors.New("missing state")
	}

	return raw.(ptnet.State), nil
}

// validate event against guard conditions
func evalGuards(event *ptnet.Event) error {
	if event.Action == ptnet.BEGIN {
		// REVIEW: should identity making offer be validated?
		return nil
	}

	txn := Contracts[event.Schema].db.Txn(false)
	raw, _ := txn.First(ContractTable, "id", event.Oid)

	if raw == nil {
		return errors.New("missing contract " + event.Schema + "." + event.Oid)
	}

	c := raw.(Declaration)

	currentState, _ := state(event.Schema, event.Oid)

	for i, g := range c.Guards {
		_, err := ptnet.VectorAdd(currentState.Vector, ptnet.Transition(g), 1)

		if err != nil {
			continue
		}
		if event.SignatureValid(c.Outputs[i].Address) {
			return nil
		}
	}
	return errors.New("failed guard condition")
}

// sign and commit event
func Commit(cmd Command, privKey identity.PrivateKey) (*ptnet.Event, error) {
	return Transform(cmd, func(evt *ptnet.Event) error {
		SignEvent(evt, privKey)
		return nil
	})
}

func compress(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	var b bytes.Buffer
	zw := gzip.NewWriter(&b)
	_, err := zw.Write(data)
	zw.Flush()
	zw.Close()
	if err != nil {
		panic("failed to compress")
	}
	return b.Bytes()
}

// commit event sign with callback
func Transform(cmd Command, signfunc func(*ptnet.Event) error) (*ptnet.Event, error) {
	return ptnet.Transform(cmd.Schema, cmd.ContractID, cmd.Action, cmd.Amount, compress(cmd.Payload), func(evt *ptnet.Event) error {
		if nil != signfunc(evt) {
			panic("failed to sign event")
		}
		return evalGuards(evt)
	})

}

func Exists(schema string, contractID string) bool {
	txn := Contracts[schema].db.Txn(false)
	raw, err := txn.First(ContractTable, "id", contractID)
	if err != nil {
		panic(err)
	}

	return raw != nil
}

func getState(schema string, contractID string) (ptnet.State, error) {
	txn := ptnet.Txn(schema, false)
	raw, _ := txn.First(ptnet.StateTable, "id", contractID)
	if raw == nil {
		return ptnet.State{}, errors.New("State not found")
	}
	return raw.(ptnet.State), nil
}

func canExecute(state ptnet.State, transition ptnet.Transition, multiplier uint64) bool {
	_, err := ptnet.VectorAdd(state.Vector, transition, multiplier)
	if err == nil {
		return true
	}
	return false
}

func IsHalted(contract Declaration) bool {
	state, _ := getState(contract.Schema, contract.ContractID)
	for _, transition := range contract.Actions {
		if canExecute(state, transition, 1) {
			// REVIEW is it better to enforce state machine is halted without guards?
			// alternatively allow guard/action combination to determine halting state
			return false
		}
	}

	return true
}

func CanRedeem(contract Declaration, publicKey identity.PublicKey) bool {
	state, _ := getState(contract.Schema, contract.ContractID)
	for i, condition := range contract.Conditions {
		if !publicKey.MatchesAddress(contract.Outputs[i].Address) {
			continue
		}
		if canExecute(state, ptnet.Transition(condition), 1) {
			return true
		}
	}

	return false
}
