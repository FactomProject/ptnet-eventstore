package contracts_test

import (
	"encoding/json"
	"fmt"
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

const expectValid bool = false
const expectError bool = true

// make commits and test for expected error outcome
func commit(t *testing.T, action string, key string, expectError bool) (*ptnet.Event, error) {
	event, err := contracts.Commit(contracts.Command{
		ChainID:    contracts.CHAIN_ID, // test values
		ContractID: contracts.CONTRACT_ID,
		Schema:     ptnet.OctoeV1,   // state machine version
		Action:     action,          // state machine action
		Amount:     1,               // triggers input action 'n' times
		Payload:    nil,             // arbitrary data optionally included
		Privkey:    contracts.Identity[key], // key used to sign event
		Pubkey:     key,
	})

	var msg string
	if expectError {
		msg = fmt.Sprintf("expected action %v to return an error ", action)
	} else {
		msg = fmt.Sprintf("unexpected from action %v ", action)
	}
	AssertEqual(t, expectError, err != nil, msg)
	return event, err
}

// create a new contract instance and validate resulting event
func setUp(t *testing.T) contracts.Declaration {
	contract := contracts.TicTacToeContract()
	event, err := contracts.Create(contract, contracts.Identity[contracts.DEPOSITOR]) // FIXME event should be signed by depositor
	AssertNil(t, err)
	AssertEqual(t, true, contracts.Exists(ptnet.OctoeV1, contracts.CONTRACT_ID), "Failed to retrieve contract declaration")

	AssertEqual(t, event.Oid, contracts.CONTRACT_ID, "")
	AssertEqual(t, event.Action, ptnet.BEGIN, "")
	AssertEqual(t, event.InputState, []uint64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, "")
	AssertEqual(t, event.OutputState, []uint64{1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0}, "")

	Assert(t, !contracts.IsHalted(contract))
	return contract
}

func TestTransactionSequence(t *testing.T) {
	contract := setUp(t)

	// start the game with a valid move
	commit(t, "X11", contracts.PLAYERX, expectValid)

	commit(t, "X22", contracts.PLAYERX, expectError) // out of turn
	commit(t, "O11", contracts.PLAYERO, expectError) // move taken already
	commit(t, "O01", contracts.PLAYERX, expectError) // sign with wrong key

	// more valid moves to finish the game
	commit(t, "O01", contracts.PLAYERO, expectValid)
	commit(t, "X00", contracts.PLAYERX, expectValid)
	commit(t, "O02", contracts.PLAYERO, expectValid)
	commit(t, "X22", contracts.PLAYERX, expectValid)

	// depositor closes the game with a winner judgement
	commit(t, "WINX", contracts.DEPOSITOR, expectValid)

	// test conditions after halting state
	Assert(t, contracts.IsHalted(contract))
	Assert(t, contracts.CanRedeem(contract, contracts.PLAYERX)) // redeemable by winner only
	Assert(t, !contracts.CanRedeem(contract, contracts.PLAYERO))
	Assert(t, !contracts.CanRedeem(contract, contracts.DEPOSITOR))
}
