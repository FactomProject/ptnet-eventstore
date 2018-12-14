package contract_test

import (
	"fmt"
	"github.com/FactomProject/ptnet-eventstore/contract"
	. "github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
	. "github.com/stackdump/gopetri/statemachine"
	"github.com/stretchr/testify/assert"
	"testing"
)

const expectValid bool = false
const expectError bool = true

// make commits and test for expected error outcome
func commit(t *testing.T, action string, key PrivateKey, expectError bool) (*ptnet.Event, error) {
	pub := PublicKey{}
	copy(pub[:], x.PrivateKeyToPub(key[:]))

	event, err := contract.Commit(contract.Command{
		ChainID:    contract.CHAIN_ID, // test values
		ContractID: "|OctoeContractID|",
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

func TestTransactionSequence(t *testing.T) {
	c := contract.TicTacToeContract()

	t.Run("publish offer", func(t *testing.T) {
		event, err := contract.Create(c, contract.CHAIN_ID, Private[DEPOSITOR])
		assert.Nil(t, err)
		assert.NotEqual(t, Private[PLAYERX], Private[PLAYERO], "test keys should not be the same")
		assert.Equal(t, true, contract.Exists(ptnet.OctoeV1, c.ContractID), "Failed to retrieve contract declaration")
		assert.Equal(t, event.Oid, c.ContractID)
		assert.Equal(t, event.Action, ptnet.BEGIN, "")

		assert.Equal(t, StateVector{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0}, event.InputState)
		assert.Equal(t, StateVector{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1}, event.OutputState)
		assert.False(t, contract.IsHalted(c))

		t.Run("execute transactions", func(t *testing.T) {
			commit(t, "X11", Private[PLAYERX], expectValid) // valid
			t.Run("reject invalid transactions", func(t *testing.T) {
				commit(t, "X22", Private[PLAYERX], expectError) // out of turn
				commit(t, "O11", Private[PLAYERO], expectError) // move taken already
				commit(t, "O01", Private[PLAYERX], expectError) // sign with wrong key
			})

			// more valid moves
			commit(t, "O01", Private[PLAYERO], expectValid)
			commit(t, "X00", Private[PLAYERX], expectValid)
			commit(t, "O02", Private[PLAYERO], expectValid)
			commit(t, "X22", Private[PLAYERX], expectValid)

			// depositor completes the contract with a winner judgement
			commit(t, "WIN_X", Private[DEPOSITOR], expectValid)
		})

		t.Run("redeem completed contract", func(t *testing.T) {
			assert.True(t, contract.IsHalted(c))
			assert.True(t, contract.CanRedeem(c, Public[PLAYERX])) // redeemable by winner only
			assert.False(t, contract.CanRedeem(c, Public[PLAYERO]))
			assert.False(t, contract.CanRedeem(c, Public[DEPOSITOR]))
		})
	})

}
