package sim_test

import (
	"fmt"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/ptnet-eventstore/blockchain"
	"github.com/FactomProject/ptnet-eventstore/contract"
	"github.com/FactomProject/ptnet-eventstore/finite"
	"github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)


func TestBlockchainApi(t *testing.T) {
	b := blockchain.NewBlockchain("Merged", ptnet.OctoeV1, ptnet.OptionV1)

	d := b.Tokens[0]
	assert.Equal(t, d.Color, blockchain.Default)

	a := b.GetAccount("BANK")
	u1 := b.GetAccount("USER1")

	//assert.Equal(t, b.ChainID, "c2487b6e6b7ae49aa813700205735dc40f7a2fef314eb5fcc28d564b217c58f3")

	t.Run("Deploy Registry chain", func(t *testing.T) {
		println(blockchain.Metachain().String())
		entry, _ := blockchain.DeployRegistry(a)
		assert.NotNil(t, entry)
		println(entry.String())
		entry, _ = blockchain.Metachain().Publish(a)
		println(entry.String())
	})

	t.Run("Publish Schemata", func(t *testing.T) {
		entry, _ := b.Publish(a)
		assert.NotNil(t, entry)
		println(entry.String())
	})

	t.Run("Deploy Chain", func(t *testing.T) {
		entry, _ := b.Deploy(a)
		assert.Equal(t, entry.ChainID, b.ChainID)
		println(entry.String())

		t.Run("Integrity Check", func(t *testing.T) {
			assert.Equal(t, x.Decode(entry.Content), x.Decode(b.Digest()))
		})
	})

	t.Run("Add Entry", func(t *testing.T) {
		ts := x.Encode(fmt.Sprintf("%v", time.Now().Unix()))
		body :=x.Encode(fmt.Sprintf("hello@%v", time.Now().Unix()))
		extids := [][]byte{ts}
		entry, _ := b.Commit(a, extids, body)
		assert.Equal(t, entry.ChainID, b.ChainID)
	})

	t.Run("Add Offer", func(t *testing.T) {
		entry, _ := b.Offer(finite.OptionContract(), a)
		println(entry.String())

		t.Run("Execute Command", func(t *testing.T) {
			c := b.Contracts[ptnet.OptionV1]

			// make commits and test for expected error outcome
			commit := func(action string, a *identity.Account) (*factom.Entry, error) {

				pub := identity.PublicKey{}
				copy(pub[:], x.PrivateKeyToPub(a.Priv.Key[:]))

				cmd := contract.Command{
					ChainID:    b.ChainID,
					ContractID: c.Template.ContractID,
					Schema:     ptnet.OptionV1,
					Action:     action,         // state machine action
					Mult:       1,              // triggers input action 'n' times
					Payload:    nil,            // arbitrary data optionally included
					Pubkey:     pub,
				}

				return b.Execute(cmd, u1)
			}

			e, _ := commit("OPT_1", u1)
			assert.False(t, "" == e.ChainID)
			println(e.String())

		})
	})


}
