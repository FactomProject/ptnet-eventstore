package sim_test

import (
	"fmt"
	"github.com/FactomProject/ptnet-eventstore/sim"
	"github.com/FactomProject/ptnet-eventstore/x"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestSendingCommitAndReveal(t *testing.T) {

	id := "92475004e70f41b94750f4a77bf7b430551113b25d3d57169eadca5692bb043d"
	extids := [][]byte{x.Encode("foo"), x.Encode("bar")}
	a := x.AccountFromFctSecret("Fs2zQ3egq2j99j37aYzaCddPq9AF3mgh64uG9gRaDAnrkjRx3eHs")
	b := x.GetBankAccount()
	numEntries := 21 // including head entry

	// TODO refactor to use
	_ = sim.BlockChain{}

	t.Run("generate accounts", func(t *testing.T) {
		println(b.String())
		println(a.String())
	})

	t.Run("Run sim to create entries", func(t *testing.T) {
		state0 := sim.Setup(t, 20)


		t.Run("Fund EC Address", func(t *testing.T) {
			x.FundECWallet(state0, b.FctPrivHash(), a.EcAddr(), 444*state0.GetFactoshisPerEC())
			bal := x.WaitForEcBalance(state0, a.EcPub())
			assert.Equal(t, bal, int64(444))
		})

		t.Run("Create Chain", func(t *testing.T) {
			e := x.Entry(id, extids, x.Encode("Hello World!"))
			c := x.NewChain(&e)
			commit, _ := x.ComposeChainCommit(a.Priv, c)
			reveal, _ := x.ComposeRevealEntryMsg(a.Priv, c.FirstEntry)
			sim.Dispatch(commit)
			sim.Dispatch(reveal)
		})

		t.Run("Create Entries", func(t *testing.T) {
			publish := func(i int) {
				e := x.Entry(id, extids, x.Encode(fmt.Sprintf("hello@%v", i)))
				commit, _ := x.ComposeCommitEntryMsg(a.Priv, e)
				reveal, _ := x.ComposeRevealEntryMsg(a.Priv, &e)
				sim.Dispatch(commit)
				sim.Dispatch(reveal)
			}

			for x := 1; x < numEntries; x++ {
				publish(x)
			}
		})

		t.Run("End simulation", func(t *testing.T) {
			x.WaitBlocks(state0, 2)
			x.WaitForAllNodes(state0)
			x.ShutDownEverything(t)
		})

		t.Run("Verify Entries", func(t *testing.T) {

			bal := x.GetBalanceEC(state0, a.EcPub())
			assert.Equal(t, int64(444-numEntries-10), bal, "EC spend mismatch")

			for _, v := range state0.Holding {
				s, _ := v.JSONString()
				println(s)
			}

			// TODO: actually check for confirmed entries
			assert.Equal(t, 0, len(state0.Holding), "messages stuck in holding")
		})

	})
}
