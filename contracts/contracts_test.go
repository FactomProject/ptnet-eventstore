package contracts_test

import (
	"encoding/json"
	"github.com/FactomProject/ptnet-eventstore/contracts"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"testing"
)

func AssertEqual(t *testing.T, a interface{}, b interface{}, msg string) {
	x, _ := json.Marshal(a)
	y, _ := json.Marshal(b)
	lhs := string(x)
	rhs := string(y)
	if lhs != rhs {
		t.Fatalf("%v != %v  %s", lhs, rhs, msg)
	}
}

func AssertNil(t *testing.T, a interface{}) {
	if a != nil {
		t.Fatalf("%v != %v", a, nil)
	}
}

var chainID string = "|ChainID|"
var contractID string = "|ContractID|"

func commit(action string) (*ptnet.Event, error) {
	event, err := contracts.Commit(contracts.Command{
		ChainID:    chainID,
		ContractID: contractID, // contract uuid
		Schema:     ptnet.OctoeV1, // state machine version
		Action:     action, // state machine action
		Amount:     1, // triggers input action 'n' times
		Payload:    nil, // arbitrary data optionally included
	})

	return event, err
}

func TestCommit(t *testing.T) {
	var event *ptnet.Event
	var err error

	contract := contracts.TicTacToeContract()
	event, err = contracts.Create(contract)

	AssertNil(t, err)
	AssertEqual(t, event.Oid, contractID, "")
	AssertEqual(t, event.Action, ptnet.BEGIN, "")
	AssertEqual(t, event.InputState, []uint64{1,1,1,1,1,1,1,1,1,1,1,1,1,1}, "")
	AssertEqual(t, event.OutputState, []uint64{1,1,1,1,1,1,1,1,1,0,1,0,0,0}, "")

	AssertEqual(t, false, contracts.IsHalted(contract), "Contract should not be halted")

	// FIXME
	//event, err = commit("X11")
	//AssertNil(t, err)
}