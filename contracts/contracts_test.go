package contracts_test

import (
	"encoding/json"
	"github.com/FactomProject/ptnet-eventstore/contracts"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"testing"
)

func Assert(t *testing.T, a interface{}) {
	if a != true {
		t.Fatalf("%v != %v", a, nil)
	}
}

func AssertNil(t *testing.T, a interface{}) {
	if a != nil {
		t.Fatalf("%v != %v", a, nil)
	}
}

func AssertEqual(t *testing.T, a interface{}, b interface{}, msg string) {
	x, _ := json.Marshal(a)
	y, _ := json.Marshal(b)
	lhs := string(x)
	rhs := string(y)
	if lhs != rhs {
		t.Fatalf("%v != %v  %s", lhs, rhs, msg)
	}
}

var chainID string = "|ChainID|"
var contractID string = "|ContractID|"
var depositorKey string = "|DepositorSecret|"
var playerXKey string = "|PlayerXSecret|"
var playerOKey string = "|PlayerOSecret|"

func commit(action string, key string) (*ptnet.Event, error) {
	event, err := contracts.Commit(contracts.Command{
		ChainID:    chainID,
		ContractID: contractID, // contract uuid
		Schema:     ptnet.OctoeV1, // state machine version
		Action:     action, // state machine action
		Amount:     1, // triggers input action 'n' times
		Payload:    nil, // arbitrary data optionally included
		Privkey:    key, // key used to sign event
	})

	return event, err
}

func TestTransactionSequence(t *testing.T) {
	var event *ptnet.Event
	var err error

	contract := contracts.TicTacToeContract()
	event, err = contracts.Create(contract) // FIXME event should be signed by depositor
	AssertNil(t, err)
	AssertEqual(t, true, contracts.Exists(ptnet.OctoeV1, contractID), "Failed to retrieve contract declaration")

	AssertEqual(t, event.Oid, contractID, "")
	AssertEqual(t, event.Action, ptnet.BEGIN, "")
	AssertEqual(t, event.InputState, []uint64{1,1,1,1,1,1,1,1,1,1,1,1,1,1}, "")
	AssertEqual(t, event.OutputState, []uint64{1,1,1,1,1,1,1,1,1,0,1,0,0,0}, "")

	Assert(t, !contracts.IsHalted(contract))

	event, err = commit("X11", playerXKey)
	AssertNil(t, err) // valid

	_, err = commit("X22", playerXKey)
	Assert(t, err != nil) // invalid because move is out of turn

	_, err = commit("O11", playerOKey)
	Assert(t, err != nil) // invalid because move is already taken

	// TODO: also test that events are rejected if signed by the wrong key

	_, err = commit("O01", playerOKey)
	AssertNil(t, err)

	_, err = commit("X00", playerXKey)
	AssertNil(t, err)

	_, err = commit("O02", playerOKey)
	AssertNil(t, err)

	_, err = commit("X22", playerXKey)
	AssertNil(t, err)

	_, err = commit("WINX", depositorKey)
	AssertNil(t, err)

	Assert(t, contracts.IsHalted(contract))

	// TODO: test that contract is redeemable by playerX
}