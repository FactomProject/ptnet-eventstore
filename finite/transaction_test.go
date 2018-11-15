package finite_test

import (
	"fmt"
	"github.com/FactomProject/ptnet-eventstore/contract"
	"github.com/FactomProject/ptnet-eventstore/finite"
	. "github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/stretchr/testify/assert"
	"testing"
)
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
	assert.Equal(t, expectError, err != nil, msg)
	return txn, err
}

func TestTransactionSequence(t *testing.T) {
	offer := finite.OptionContract()

	t.Run("publish offer", func(t *testing.T) {
		txn := finite.OfferTransaction(offer, Identity[DEPOSITOR])
		assert.Equal(t, true, contract.Exists(offer.Schema, contract.CONTRACT_ID), "missing declaration")
		assert.Equal(t, txn.Oid, contract.CONTRACT_ID, "")
		assert.Equal(t, txn.Action, ptnet.BEGIN, "")
		assert.Equal(t, txn.InputState, ptnet.StateVector{1, 1, 1, 1, 1})
		assert.Equal(t, txn.OutputState, ptnet.StateVector{1, 0, 0, 1, 0})
		assert.False(t, contract.IsHalted(offer.Declaration))
	})

	t.Run("execute transactions to accept offer", func(t *testing.T) {
		commit(t, "OPT_1", USER1, expectValid)
		commit(t, "OPT_2", DEPOSITOR, expectError) // only first executed option is valid
		commit(t, "HALT", DEPOSITOR, expectValid)
	})

	t.Run("redeem completed contract", func(t *testing.T) {
		assert.True(t, contract.IsHalted(offer.Declaration))
		assert.True(t, contract.CanRedeem(offer.Declaration, USER1))
		assert.False(t, contract.CanRedeem(offer.Declaration, DEPOSITOR))
		assert.False(t, contract.CanRedeem(offer.Declaration, USER2))
	})
}
