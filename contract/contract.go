package contract

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FactomProject/ptnet-eventstore/gen"
	"github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
	"github.com/hashicorp/go-memdb"
	. "github.com/stackdump/gopetri/statemachine"
	"text/template"
)

type Contract struct {
	Schema   string       `json:"schema"`
	Machine  StateMachine `json:"-"`
	Template Declaration  `json:"template"`
	db       *memdb.MemDB
}

type AddressAmountMap struct {
	Address []byte `json:"address""`
	Amount  uint64 `json:"amount""`
	Token   uint8  `json:"color""`
}

type Condition Transition

type Variables struct {
	ContractID  string             `json:"contractid"`
	BlockHeight uint64             `json:"blockheight"`
	Inputs      []AddressAmountMap `json:"inputs"`
	Outputs     []AddressAmountMap `json:"outputs"`
}

type Invariants struct {
	Schema     string                `json:"schema"`
	Parameters map[string]Place      `json:"parameters"`
	Capacity   StateVector           `json:"capacity"`
	State      StateVector           `json:"state"`
	Actions    map[Action]Transition `json:"actions"`
	Guards     []Condition           `json:"guards"`     // enforces contract roles
	Conditions []Condition           `json:"conditions"` // enforce redeem conditions
}

type Declaration struct {
	Variables
	Invariants
}

type State struct {
	ChainID   string      `json:"chainid"`
	LastEntry string      `json:"last_entry"`
	ChainHead string      `json:"chainhead"`
	State     StateVector `json:"state"`
}

type Command struct {
	ChainID    string
	ContractID string
	Schema     string
	Action     string
	Mult       uint64
	Payload    []byte
	Pubkey     identity.PublicKey //        compare w/ factom identity standard
}

var Contracts map[string]Contract = map[string]Contract{
	ptnet.FiniteV1: Contract{
		Schema:   ptnet.FiniteV1,
		Machine:  gen.FiniteV1.StateMachine(),
		Template: RegistryTemplate(),
		db:       ContractStore(),
	},
	ptnet.Spend: Contract{
		Schema:   ptnet.Spend,
		Machine:  gen.Spend.StateMachine(),
		Template: SpendTemplate(),
		db:       ContractStore(),
	},
	ptnet.Tip: Contract{
		Schema:   ptnet.Tip,
		Machine:  gen.Tip.StateMachine(),
		Template: SpendTemplate(),
		db:       ContractStore(),
	},
	ptnet.AuctionV1: Contract{
		Schema:   ptnet.AuctionV1,
		Machine:  gen.AuctionV1.StateMachine(),
		Template: SpendTemplate(),
		db:       ContractStore(),
	},
	ptnet.OptionV1: Contract{
		Schema:   ptnet.OptionV1,
		Machine:  gen.OptionV1.StateMachine(),
		Template: OptionTemplate(),
		db:       ContractStore(),
	},
	ptnet.OctoeV1: Contract{
		Schema:   ptnet.OctoeV1,
		Machine:  gen.OctoeV1.StateMachine(),
		Template: TicTacToeTemplate(),
		db:       ContractStore(),
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

	// FIXME convert private to pub
	pubkey := identity.PublicKey{}

	event, err := Transform(
		Command{
			ChainID:    chainID, // test values
			ContractID: contract.ContractID,
			Schema:     contract.Schema,
			Action:     ptnet.EXEC, // state machine action
			Mult:       1,          // triggers input action 'n' times
			Payload:    payload,    // arbitrary data optionally included
			Pubkey:     pubkey,     // REVIEW: will there always be a single input?
		}, signfunc)

	if err != nil {
		panic(err)
	}

	c, ok := Contracts[contract.Schema]
	if !ok {
		panic(fmt.Sprintf("Unknown Schema %v\n", contract.Schema))
	}

	txn := c.db.Txn(true)
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
	if event.Action == ptnet.EXEC {
		// REVIEW: should identity making offer be validated?
		return nil
	}

	txn := Contracts[event.Schema].db.Txn(false)
	raw, _ := txn.First(ContractTable, "id", event.Oid)

	if raw == nil {
		return errors.New("missing contract " + event.Schema + "." + event.Oid)
	}


	currentState, _ := state(event.Schema, event.Oid)
	c := raw.(Declaration)

	for i, g := range c.Guards {
		_, err := ptnet.VectorAdd(currentState.Vector, Transition(g), 1)

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

func Compress(data []byte) []byte {
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
	return ptnet.Transform(cmd.Schema, cmd.ContractID, cmd.Action, cmd.Mult, cmd.Payload, func(evt *ptnet.Event) error {
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

func canExecute(state ptnet.State, transition Transition, multiplier uint64) bool {
	_, err := ptnet.VectorAdd(state.Vector, transition, multiplier)
	if err == nil {
		return true
	}
	return false
}

func IsHalted(contract Declaration) bool {
	state, _ := getState(contract.Schema, contract.ContractID)
	for _, transition := range contract.Actions {
		_, err := ptnet.VectorAdd(state.Vector, transition, 1)

		if err != nil {
			continue
		}

		for _, g := range contract.Guards {
			if canExecute(state, Transition(g), 1) {
				return false
			}
		}
	}

	return true
}

func CanRedeem(contract Declaration, publicKey identity.PublicKey) bool {
	state, _ := getState(contract.Schema, contract.ContractID)

	// in this case contract is invalid
	if len(contract.Conditions) > len(contract.Outputs) {
		return false
	}

	for i, condition := range contract.Conditions {
		if !publicKey.MatchesAddress(contract.Outputs[i].Address) {
			continue
		}
		if canExecute(state, Transition(condition), 1) {
			return true
		}
	}

	return false
}

var contractFormat string = `
Inputs: {{ range $_, $input := .Inputs}}
	Address: {{ printf "%x" $input.Address }} Amount: {{ $input.Amount }} Token: {{ $input.Token }} {{ end }}
Outputs: {{ range $_, $output := .Outputs}}
	Address: {{ printf "%x" $output.Address }} Amount: {{ $output.Amount }} Token: {{ $output.Token }} {{ end }}
BlockHeight: {{ .BlockHeight }}
ContractID: {{ .ContractID }}
Schema: {{ .Schema }}
State: {{ .GetState }}
Actions: {{ range $key, $action := .Actions }}
	{{$key}}: {{ $action }}{{ end }}
Guards: {{ range $i, $guard := .Guards }}
	{{$i}}: {{ $guard }}{{ end }}
Conditions: {{ range $i, $condition := .Conditions }}
	{{$i}}: {{ $condition }}{{ end }}
`
var contractTemplate *template.Template = template.Must(
	template.New("").Parse(contractFormat),
)

type contractSource struct {
	Declaration
}

func (c contractSource) GetState() (s []uint64) {
	return ptnet.ToVector(c.State)
}

func (contract Declaration) String() string {
	b := &bytes.Buffer{}
	contractTemplate.Execute(b, contractSource{contract})
	return b.String()
}

// Dual is a vector that won't be used as a transformation
// instead these vectors provide a way to simulate additional actions
// to test the current state vector by way of subtraction
func Dual(p PetriNet, placeNames []string, mult int64) []int64 {
	role := p.GetEmptyVector()
	for _, k := range placeNames {
		attr, ok := p.Places[k]
		if ok {
			role[attr.Offset] = mult * -1 // test by subtraction
		} else {
			panic(fmt.Sprintf("unknown place: %v", k))
		}
	}
	return role
}

// alias for syntactic sugar
var Role = Dual
var Check = Role

func (c Contract) Version() []byte {
	data, _ := json.Marshal(c.Template)
	//fmt.Printf("%s", data)
	return x.Shad(data)
}
