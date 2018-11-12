package finite_test

import (
	"encoding/json"
	"fmt"
	"github.com/FactomProject/ptnet-eventstore/contract"
	"github.com/FactomProject/ptnet-eventstore/finite"
	. "github.com/FactomProject/ptnet-eventstore/identity"
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
	// TODO: convert to using ExecuteTransaction
	event, err := contract.Commit(contract.Command{
		ChainID:    contract.CHAIN_ID, // test values
		ContractID: contract.CONTRACT_ID,
		Schema:     ptnet.OptionV1, // state machine version
		Action:     action,         // state machine action
		Amount:     1,              // triggers input action 'n' times
		Payload:    nil,            // arbitrary data optionally included
		Pubkey:     key,
	}, Identity[key])

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
func setUp(t *testing.T) finite.Transaction {
	offer := finite.OptionContract()
	txn := finite.OfferTransaction(offer, Identity[DEPOSITOR])
	//AssertNil(t, err)
	AssertEqual(t, true, contract.Exists(offer.Schema, contract.CONTRACT_ID), "Failed to retrieve contract declaration")

	AssertEqual(t, txn.Oid, contract.CONTRACT_ID, "")
	AssertEqual(t, txn.Action, ptnet.BEGIN, "")
	AssertEqual(t, txn.InputState, []uint64{1, 1, 1, 1, 1}, "")
	AssertEqual(t, txn.OutputState, []uint64{1, 0, 0, 1, 0}, "")

	Assert(t, !contract.IsHalted(offer.Declaration))
	return txn
}

func TestTransactionSequence(t *testing.T) {
	offerTxn := setUp(t)
	_ = offerTxn

	/*
	commit(t, "OPT_1", USER1, expectValid)
	commit(t, "OPT_2", DEPOSITOR, expectError) // only one option can be selected
	commit(t, "HALT", DEPOSITOR, expectValid)
	Assert(t, contract.IsHalted(c.Declaration))
	Assert(t, contract.CanRedeem(c.Declaration, USER1)) // redeemable by winner only
	Assert(t, !contract.CanRedeem(c.Declaration, DEPOSITOR))
	Assert(t, !contract.CanRedeem(c.Declaration, USER2))
	*/
}
