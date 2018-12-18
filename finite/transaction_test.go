package finite_test

import (
	"encoding/json"
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

func TestTransactionSequence(t *testing.T) {
	offer := finite.OptionContract()

	// make commits and test for expected error outcome
	commit := func(action string, key PrivateKey, expectError bool) (finite.Transaction, error) {
		pub := PublicKey{}
		copy(pub[:], x.PrivateKeyToPub(key[:]))
		txn, err := finite.ExecuteTransaction(
			contract.Command{
				ChainID:    offer.ChainID,
				ContractID: offer.ContractID,
				Schema:     ptnet.OptionV1,
				Action:     action, // state machine action
				Mult:       1,      // triggers input action 'n' times
				Payload:    nil,    // arbitrary data optionally included
				Pubkey:     pub,
			}, key)

		d, _ := json.Marshal(txn)
		fmt.Printf("%s", d)

		var msg string
		if expectError {
			msg = fmt.Sprintf("expected action %v to return an error ", action)
		} else {
			msg = fmt.Sprintf("unexpected error %v from action %v", err, action)
		}
		assert.Equal(t, expectError, err != nil, msg)
		return txn, err
	}

	t.Run("publish offer", func(t *testing.T) {
		txn := finite.OfferTransaction(offer, Private[DEPOSITOR])
		assert.Equal(t, true, contract.Exists(offer.Schema, offer.ContractID), "missing declaration")
		assert.Equal(t, txn.Oid, offer.ContractID)
		assert.Equal(t, txn.Action, ptnet.EXEC, "")
		assert.Equal(t, txn.InputState, StateVector{1, 0, 0, 0, 0})
		assert.Equal(t, txn.OutputState, StateVector{0, 1, 0, 0, 0})
		assert.False(t, contract.IsHalted(offer.Declaration))
	})

	t.Run("execute transactions to accept offer", func(t *testing.T) {
		commit("OPT_1", Private[USER1], expectValid)
		commit("OPT_2", Private[DEPOSITOR], expectError) // only first executed option is valid
		commit("HALT", Private[DEPOSITOR], expectError)
	})

	t.Run("redeem completed contract", func(t *testing.T) {
		assert.True(t, contract.IsHalted(offer.Declaration))
		assert.True(t, contract.CanRedeem(offer.Declaration, Public[USER1]))
		assert.False(t, contract.CanRedeem(offer.Declaration, Public[DEPOSITOR]))
		assert.False(t, contract.CanRedeem(offer.Declaration, Public[USER2]))
	})
}
