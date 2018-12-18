package contract_test

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


func TestTransactionSequence(t *testing.T) {
	c := finite.TicTacToeContract()
	// make commits and test for expected error outcome
	commit := func (action string, key PrivateKey, expectError bool) (*ptnet.Event, error) {
		pub := PublicKey{}
		copy(pub[:], x.PrivateKeyToPub(key[:]))

		event, err := contract.Commit(contract.Command{
			ChainID:    c.ChainID,
			ContractID: c.ContractID, // FIXME this should be unique per usage instead is set to chainID
			Schema:     ptnet.OctoeV1, // state machine version
			Action:     action,        // state machine action
			Mult:       1,             // triggers input action 'n' times
			Payload:    nil,           // arbitrary data optionally included
			Pubkey:     pub,
		}, key)

		if expectError {
			assert.NotNil(t, err, fmt.Sprintf("expected action %v to return an error ", action))
		} else {
			assert.Nil(t, err, fmt.Sprintf("unexpected from action %v ", action))
		}
		//println(event.String())
		return event, err
	}

	t.Run("publish offer", func(t *testing.T) {
		event, err := contract.Create(c.Declaration, c.ChainID, Private[DEPOSITOR])
		assert.Nil(t, err)
		assert.NotEqual(t, Private[PLAYERX], Private[PLAYERO], "test keys should not be the same")
		assert.Equal(t, true, contract.Exists(ptnet.OctoeV1, c.ContractID), "Failed to retrieve contract declaration")
		assert.Equal(t, event.Oid, c.ContractID)
		assert.Equal(t, event.Action, ptnet.BEGIN, "")

		assert.Equal(t, StateVector{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0}, event.InputState)
		assert.Equal(t, StateVector{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1}, event.OutputState)
		assert.False(t, contract.IsHalted(c.Declaration))

		t.Run("execute transactions", func(t *testing.T) {
			commit("X11", Private[PLAYERX], expectValid) // valid
			t.Run("reject invalid transactions", func(t *testing.T) {
				commit("X22", Private[PLAYERX], expectError) // out of turn
				commit("O11", Private[PLAYERO], expectError) // move taken already
				commit("O01", Private[PLAYERX], expectError) // sign with wrong key
			})

			// more valid moves
			commit("O01", Private[PLAYERO], expectValid)
			commit("X00", Private[PLAYERX], expectValid)
			commit("O02", Private[PLAYERO], expectValid)
			commit("X22", Private[PLAYERX], expectValid)

			// depositor completes the contract with a winner judgement
			commit("WIN_X", Private[DEPOSITOR], expectValid)
		})


		t.Run("redeem completed contract", func(t *testing.T) {
			assert.True(t, contract.IsHalted(c.Declaration))
			assert.True(t, contract.CanRedeem(c.Declaration, Public[PLAYERX])) // redeemable by winner only
			assert.False(t, contract.CanRedeem(c.Declaration, Public[PLAYERO]))
			assert.False(t, contract.CanRedeem(c.Declaration, Public[DEPOSITOR]))
		})
	})

}
