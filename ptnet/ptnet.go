package ptnet

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/hashicorp/go-memdb"
	. "github.com/stackdump/gopetri/statemachine"
	"text/template"
	"time"
)

// Reserved Actions
const EXEC string = "EXEC"

// Token colors
const (
	Default = uint8(iota)
	Coin
	Karma
	AntiKarma
)

type State struct {
	Oid       string      `json:"oid"`
	Vector    StateVector `json:"vector"`
	Timestamp uint64      `json:"timestamp"`
}

type Machine struct {
	StateMachine
	db *memdb.MemDB
}

type Event struct {
	Timestamp   uint64               `json:"timestamp"`
	Schema      string               `json:"schema"`
	Action      string               `json:"action"`
	Oid         string               `json:"oid"`
	Mult        uint64               `json:"value"`
	InputState  StateVector          `json:"input"`
	OutputState StateVector          `json:"output"`
	Payload     []byte               `json:"payload"`
	pubkeys     []identity.PublicKey // TODO
	signatures  [][]byte             // signatures
	digest      []byte
}

// start a new transaction with in-memory db
func Txn(schema string, write bool) *memdb.Txn {
	return StateMachines[schema].db.Txn(write)
}

// change state by appending a valid event to the events table
func Commit(schema string, oid string, action string, value uint64, payload []byte) (*Event, error) {

	event := Event{
		Timestamp:   uint64(time.Now().UnixNano()),
		Schema:      schema,
		Action:      action,
		Oid:         oid,
		Mult:        value,
		InputState:  nil,
		OutputState: nil,
		Payload:     payload,
	}

	err := applyTransform(StateMachines[schema], &event, true, beforeCommit, afterCommit)
	return &event, err
}

func Transform(schema string, oid string, action string, mult uint64, payload []byte, beforeCommitCallback func(*Event) error) (*Event, error) {

	event := Event{
		Timestamp:   uint64(time.Now().UnixNano()),
		Schema:      schema,
		Action:      action,
		Oid:         oid,
		Mult:        mult,
		InputState:  nil,
		OutputState: nil,
		Payload:     payload,
	}

	err := applyTransform(StateMachines[schema], &event, true, beforeCommitCallback, afterCommit)
	return &event, err
}

func AddSignature(event *Event, publicKey identity.PublicKey, sig []byte) {
	if event.digest == nil {
		panic("must add digest before affixing signature")
	}

	event.pubkeys = append(event.pubkeys, publicKey)
	event.signatures = append(event.signatures, sig)
}

func (event *Event) SignatureValid(address []byte) bool {
	// TODO: also validate sig against pubkey
	// or consider doing this validation in the contract layer
	// FIXME
	for _, key := range event.pubkeys {
		if key.MatchesAddress(address) {
			return true
		}
	}
	return false
}

func (event *Event) AddDigest() {
	data, _ := json.Marshal(event)
	h := sha256.New()
	h.Write(data)
	event.digest = h.Sum(nil)
}

func (event *Event) GetDigest() []byte {
	return event.digest
}

func (event *Event) GetPubkeys() []identity.PublicKey {
	return event.pubkeys
}

func (event *Event) GetSignatures() [][]byte {
	return event.signatures
}

func encodeEvent(event *Event) *bytes.Buffer {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(event)
	if err != nil {
		panic(err)
	}
	return &b
}

func decodeEvent(b *bytes.Buffer) *Event {
	var a *Event
	dec := gob.NewDecoder(b)
	err := dec.Decode(&a)
	if err != nil {
		panic(err)
	}
	return a
}

func afterCommit(evt *Event) {
	// REVIEW: consider storing event in DB
	/*
		data := encodeEvent(evt)
		fmt.Printf("storagePersist => %v\n", decodeEvent(data))
	*/
}

func beforeCommit(evt *Event) error {
	return nil
}

// update apply vector addition and update output State
// FIXME refactor to use StateMachine methods
func VectorAdd(vectorIn []uint64, transform []int64, multiplier uint64) ([]uint64, error) {
	var vectorOut []uint64
	var err error = nil

	for offset, delta := range transform {
		val := int64(vectorIn[offset]) + delta*int64(multiplier)
		// FIXME check against capacity
		if val >= 0 {
			vectorOut = append(vectorOut, uint64(val))
		} else {
			err = errors.New(fmt.Sprintf("invalid output: %v offset: %v", val, offset))
			break
		}
	}
	return vectorOut, err

}

func vectorApply(vectorIn []uint64, transform []int64, multiplier uint64, capacity StateVector, stateOut *State) error {
	out, err := VectorAdd(vectorIn, transform, multiplier)

	if err != nil {
		return err
	}

	for i, v := range capacity {
		if v > 0 && out[i] > v {
			return errors.New(fmt.Sprintf("Exceeded Capacity[%v]: %v", i, v))
		}
	}

	stateOut.Vector = out
	return err
}

// add transform*multiplier to input vector and save valid output to state
// optionally persist to the eventstore
func applyTransform(machine Machine, event *Event, persistEvent bool, precondition func(*Event) error, callback func(*Event)) error {

	if machine.Initial == nil {
		return errors.New(fmt.Sprintf("Unknown schema: %v", event.Schema))
	}

	actionVector := machine.Transitions[Action(event.Action)]
	if actionVector == nil {
		return errors.New(fmt.Sprintf("Unknown action: %v", event.Action))
	}

	txn := machine.db.Txn(true)
	raw, err := txn.First(StateTable, "id", event.Oid)
	if err != nil {
		return err
	}

	var inputVector []uint64

	if raw == nil {
		// REVIEW: eventually refactor to make this a cache-miss
		// use initial iff a given OID not found in leveldb
		inputVector = machine.StateMachine.Initial
	} else {
		inputVector = raw.(State).Vector
	}

	event.InputState = inputVector
	outState := State{Oid: event.Oid, Timestamp: event.Timestamp}

	err = vectorApply(inputVector, actionVector, event.Mult, machine.Capacity, &outState)
	if err != nil {
		txn.Abort()
		return err
	}

	event.AddDigest()
	err = precondition(event)
	if err != nil {
		txn.Abort()
		return err
	}

	err = txn.Insert(StateTable, outState)
	if err != nil {
		panic(err)
	}

	if persistEvent {
		event.OutputState = outState.Vector
		err = txn.Insert(EventTable, event)
		if err != nil {
			panic(err)
		}
	}

	callback(event)

	txn.Commit()
	return err
}

var eventFormat string = `
Timestamp:   {{ .Timestamp }}
Schema:      {{ .Schema }}
Action:      {{ .Action }}
Oid:         {{ .Oid }}
Mult:        {{ .Mult }}
InputState:  {{ .VectorIn }}
OutputState: {{ .VectorOut }}
Payload:
	{{ printf "%x" .Payload }}
digest:
	{{ printf "%x" .GetDigest }}
pubkeys: {{ range $_, $pubkey := .GetPubkeys }}
	{{ printf "%x" $pubkey }}{{ end }}
signatures: {{ range $_, $sig := .GetSignatures }}
	{{ printf "%x" $sig }}{{ end }}
`
var eventTemplate *template.Template = template.Must(
	template.New("").Parse(eventFormat),
)

type eventSource struct {
	*Event
}

func ToVector(vector StateVector) []uint64 {
	x := []uint64{}
	for _, v := range vector {
		x = append(x, v)
	}
	return x
}

func (s eventSource) VectorIn() []uint64 {
	return ToVector(s.InputState)
}

func (s eventSource) VectorOut() []uint64 {
	return ToVector(s.OutputState)
}

func (e *Event) String() string {
	b := &bytes.Buffer{}
	eventTemplate.Execute(b, eventSource{e})
	return b.String()
}
