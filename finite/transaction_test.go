package finite_test

import (
	"fmt"
	"github.com/FactomProject/ptnet-eventstore/contract"
	"github.com/FactomProject/ptnet-eventstore/finite"
	. "github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
	. "github.com/stackdump/gopetri/statemachine"
	"github.com/stretchr/testify/assert"
	"testing"
)

const expectValid bool = false
const expectError bool = true
const optionContractID string = "|OptionContractID|"

// make commits and test for expected error outcome
func commit(t *testing.T, action string, key PrivateKey, expectError bool) (finite.Transaction, error) {
	pub := PublicKey{}
	copy(pub[:], x.PrivateKeyToPub(key[:]))
	txn, err := finite.ExecuteTransaction(finite.Execution{
		Command: contract.Command{
			ChainID:    contract.CHAIN_ID, // test values
			ContractID: optionContractID,
			Schema:     ptnet.OptionV1, // state machine version
			Action:     action,         // state machine action
			Mult:       1,              // triggers input action 'n' times
			Payload:    nil,            // arbitrary data optionally included
			Pubkey:     pub,
		},
	}, key)

	println(txn.String())
	var msg string
	if expectError {
		msg = fmt.Sprintf("expected action %v to return an error ", action)
	} else {
		msg = fmt.Sprintf("unexpected error %v from action %v", err, action)
	}
	assert.Equal(t, expectError, err != nil, msg)
	return txn, err
}

func TestTransactionSequence(t *testing.T) {
	offer := finite.OptionContract()
	offer.ChainID = contract.CHAIN_ID

	//println(offer.String())
	t.Run("publish offer", func(t *testing.T) {
		txn := finite.OfferTransaction(offer, Private[DEPOSITOR])
		assert.Equal(t, true, contract.Exists(offer.Schema, optionContractID), "missing declaration")
		assert.Equal(t, txn.Oid, optionContractID)
		assert.Equal(t, txn.Action, ptnet.BEGIN, "")
		assert.Equal(t, txn.InputState, StateVector{1, 0, 0, 0, 0})
		assert.Equal(t, txn.OutputState, StateVector{0, 1, 0, 0, 0})
		assert.False(t, contract.IsHalted(offer.Declaration))
	})

	t.Run("execute transactions to accept offer", func(t *testing.T) {
		commit(t, "OPT_1", Private[USER1], expectValid)
		commit(t, "OPT_2", Private[DEPOSITOR], expectError) // only first executed option is valid
		commit(t, "HALT", Private[DEPOSITOR], expectError)
	})

	t.Run("redeem completed contract", func(t *testing.T) {
		assert.True(t, contract.IsHalted(offer.Declaration))
		assert.True(t, contract.CanRedeem(offer.Declaration, Public[USER1]))
		assert.False(t, contract.CanRedeem(offer.Declaration, Public[DEPOSITOR]))
		assert.False(t, contract.CanRedeem(offer.Declaration, Public[USER2]))
	})
}
