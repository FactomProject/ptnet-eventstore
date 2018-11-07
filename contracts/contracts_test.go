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

// make commits and test for expected error outcome
func commit(t *testing.T, action string, key string, expectError bool) (*ptnet.Event, error) {
	event, err := contracts.Commit(contracts.Command{
		ChainID:    contracts.CHAIN_ID, // test values
		ContractID: contracts.CONTRACT_ID,
		Schema:     ptnet.OctoeV1, // state machine version
		Action:     action, // state machine action
		Amount:     1, // triggers input action 'n' times
		Payload:    nil, // arbitrary data optionally included
		Privkey:    key, // key used to sign event
	})

	var msg string
	if expectError {
		msg = "expected commit to return an error"
	} else {
		msg = "unexpected error"
	}
	AssertEqual(t, expectError, err != nil, msg)
	return event, err
}

// create a new contract instance and validate resulting event
func setUp(t *testing.T) contracts.Declaration {
	contract := contracts.TicTacToeContract()
	event, err := contracts.Create(contract) // FIXME event should be signed by depositor
	AssertNil(t, err)
	AssertEqual(t, true, contracts.Exists(ptnet.OctoeV1, contracts.CONTRACT_ID), "Failed to retrieve contract declaration")

	AssertEqual(t, event.Oid, contracts.CONTRACT_ID, "")
	AssertEqual(t, event.Action, ptnet.BEGIN, "")
	AssertEqual(t, event.InputState, []uint64{1,1,1,1,1,1,1,1,1,1,1,1,1,1}, "")
	AssertEqual(t, event.OutputState, []uint64{1,1,1,1,1,1,1,1,1,0,1,0,0,0}, "")

	Assert(t, !contracts.IsHalted(contract))
	return contract
}

func TestTransactionSequence(t *testing.T) {
	var expectValid bool = false
	var expectError bool = true

	contract := setUp(t)

	// start the game with a valid move
	commit(t, "X11", contracts.PLAYERX_SECRET, expectValid)

    commit(t, "X22", contracts.PLAYERX_SECRET, expectError) // out of turn
	commit(t, "O11", contracts.PLAYERO_SECRET, expectError) // move taken already
	// TODO: also test that events are rejected if signed by the wrong key

	// more valid moves to finish the game
	commit(t, "O01", contracts.PLAYERO_SECRET, expectValid)
	commit(t, "X00", contracts.PLAYERX_SECRET, expectValid)
	commit(t, "O02", contracts.PLAYERO_SECRET, expectValid)
	commit(t, "X22", contracts.PLAYERX_SECRET, expectValid)

	// depositor closes the game with a winner judgement
	commit(t, "WINX", contracts.DEPOSITOR_SECRET, expectValid)


	// test conditions after halting state
	Assert(t, contracts.IsHalted(contract))
	Assert(t, contracts.CanRedeem(contract, contracts.PLAYERX)) // redeemable by winner only
	Assert(t, !contracts.CanRedeem(contract, contracts.PLAYERO))
	Assert(t, !contracts.CanRedeem(contract, contracts.DEPOSITOR))
}