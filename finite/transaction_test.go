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
func commit(t *testing.T, action string, key string, expectError bool) (finite.Transaction, error) {
	txn, err := finite.ExecuteTransaction(finite.Execution{
		Command: contract.Command{
			ChainID:    contract.CHAIN_ID, // test values
			ContractID: contract.CONTRACT_ID,
			Schema:     ptnet.OptionV1, // state machine version
			Action:     action,         // state machine action
			Amount:     1,              // triggers input action 'n' times
			Payload:    nil,            // arbitrary data optionally included
			Pubkey:     key,
		},
	}, Identity[key])

	var msg string
	if expectError {
		msg = fmt.Sprintf("expected action %v to return an error ", action)
	} else {
		msg = fmt.Sprintf("unexpected from action %v ", action)
	}
	AssertEqual(t, expectError, err != nil, msg)
	return txn, err
}

func TestTransactionSequence(t *testing.T) {
	offer := finite.OptionContract()

	t.Run("publish offer", func(t *testing.T) {
		txn := finite.OfferTransaction(offer, Identity[DEPOSITOR])
		AssertEqual(t, true, contract.Exists(offer.Schema, contract.CONTRACT_ID), "missing declaration")
		AssertEqual(t, txn.Oid, contract.CONTRACT_ID, "")
		AssertEqual(t, txn.Action, ptnet.BEGIN, "")
		AssertEqual(t, txn.InputState, []uint64{1, 1, 1, 1, 1}, "")
		AssertEqual(t, txn.OutputState, []uint64{1, 0, 0, 1, 0}, "")
		Assert(t, !contract.IsHalted(offer.Declaration))
	})

	t.Run("execute transactions to accept offer", func(t *testing.T) {
		commit(t, "OPT_1", USER1, expectValid)
		commit(t, "OPT_2", DEPOSITOR, expectError) // only first executed option is valid
		commit(t, "HALT", DEPOSITOR, expectValid)
	})

	t.Run("redeem completed contract", func(t *testing.T) {
		Assert(t, contract.IsHalted(offer.Declaration))
		Assert(t, contract.CanRedeem(offer.Declaration, USER1))
		Assert(t, !contract.CanRedeem(offer.Declaration, DEPOSITOR))
		Assert(t, !contract.CanRedeem(offer.Declaration, USER2))
	})
}
