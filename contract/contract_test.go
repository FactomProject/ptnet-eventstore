package contract_test

import (
	"fmt"
	"github.com/FactomProject/ptnet-eventstore/contract"
	. "github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/stretchr/testify/assert"
	"testing"
)

const expectValid bool = false
const expectError bool = true

// make commits and test for expected error outcome
func commit(t *testing.T, action string, key string, expectError bool) (*ptnet.Event, error) {
	event, err := contract.Commit(contract.Command{
		ChainID:    contract.CHAIN_ID, // test values
		ContractID: contract.CONTRACT_ID,
		Schema:     ptnet.OctoeV1, // state machine version
		Action:     action,        // state machine action
		Amount:     1,             // triggers input action 'n' times
		Payload:    nil,           // arbitrary data optionally included
		Pubkey:     key,
	}, Identity[key])

	var msg string
	if expectError {
		msg = fmt.Sprintf("expected action %v to return an error ", action)
	} else {
		msg = fmt.Sprintf("unexpected from action %v ", action)
	}
	assert.Equal(t, expectError, err != nil, msg)
	return event, err
}

func TestTransactionSequence(t *testing.T) {
	c := contract.TicTacToeContract()

	t.Run("publish offer", func(t *testing.T) {
		event, err := contract.CreateAndSign(c, contract.CHAIN_ID, Identity[DEPOSITOR])
		assert.Nil(t, err)
		assert.Equal(t, true, contract.Exists(ptnet.OctoeV1, contract.CONTRACT_ID), "Failed to retrieve contract declaration")
		assert.Equal(t, event.Oid, contract.CONTRACT_ID, "")
		assert.Equal(t, event.Action, ptnet.BEGIN, "")
		assert.Equal(t, event.InputState, ptnet.StateVector{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
		assert.Equal(t, event.OutputState, ptnet.StateVector{1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0})
		assert.False(t, contract.IsHalted(c))
	})


	t.Run("execute transactions to complete game", func(t *testing.T) {
		// start the game with a valid move
		commit(t, "X11", PLAYERX, expectValid)

		commit(t, "X22", PLAYERX, expectError) // out of turn
		commit(t, "O11", PLAYERO, expectError) // move taken already
		commit(t, "O01", PLAYERX, expectError) // sign with wrong key

		// more valid moves to finish the game
		commit(t, "O01", PLAYERO, expectValid)
		commit(t, "X00", PLAYERX, expectValid)
		commit(t, "O02", PLAYERO, expectValid)
		commit(t, "X22", PLAYERX, expectValid)

		// depositor closes the game with a winner judgement
		commit(t, "WINX", DEPOSITOR, expectValid)
	})

	t.Run("redeem completed contract", func(t *testing.T) {
		// test conditions after halting state
		assert.True(t, contract.IsHalted(c))
		assert.True(t, contract.CanRedeem(c, PLAYERX)) // redeemable by winner only
		assert.False(t, contract.CanRedeem(c, PLAYERO))
		assert.False(t, contract.CanRedeem(c, DEPOSITOR))
	})
}
