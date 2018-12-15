package sim_test

import (
	"fmt"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/sim"
	"github.com/FactomProject/ptnet-eventstore/x"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)


func TestSendingCommitAndReveal(t *testing.T) {
	numEntries := 21 // including head entry

	chain, _ := sim.NewBlockchain("salt", ptnet.OctoeV1, ptnet.OptionV1)
	b := chain.GetAccount("BANK")

	t.Run("generate accounts", func(t *testing.T) {
		println(b.String())
	})

	t.Run("Run sim to create entries", func(t *testing.T) {
		state0 := sim.Setup(t, 20)

		t.Run("Fund EC Address", func(t *testing.T) {
			x.FundECWallet(state0, b.PrivHash(), b.EcAddr(), 444*state0.GetFactoshisPerEC())
			bal := x.WaitForEcBalance(state0, b.EcPub())
			assert.Equal(t, bal, int64(444))
		})

		t.Run("Create Chain", func(t *testing.T) {
			_, err := chain.Deploy(b)
			assert.Nil(t, err)
		})

		t.Run("Create Entries", func(t *testing.T) {
			for i := 1; i < numEntries; i++ {
				ts := x.Encode(fmt.Sprintf("%v", time.Now().UnixNano()))
				body :=x.Encode(fmt.Sprintf("hello@%v", i))
				extids := [][]byte{ts}
				_, err := chain.Commit(b, extids, body)
				assert.Nil(t, err)
			}
		})

		t.Run("End simulation", func(t *testing.T) {
			x.WaitBlocks(state0, 2)
			x.WaitForAllNodes(state0)
			x.ShutDownEverything(t)
		})

		t.Run("Verify EC Balance", func(t *testing.T) {
			bal := x.GetBalanceEC(state0, b.EcPub())
			assert.Equal(t, int64(444-numEntries-10), bal, "EC spend mismatch")
		})

	})
}
