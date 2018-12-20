package blockchain_test

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

var IntegrityChecksum string = "670faace0fe63eadc020cd85b81c6254975a03d5b009363ed6558e85aba9cee7"

func TestBlockchainApi(t *testing.T) {
	b := blockchain.NewBlockchain("Merged", ptnet.OctoeV1, ptnet.OptionV1)

	d := b.Tokens[0]
	assert.Equal(t, d.Color, blockchain.Default)

	a := b.GetAccount("BANK")
	u1 := b.GetAccount("USER1")

	// detect when state machine declarations or contract templates are altered
	//assert.Equal(t, x.Decode(b.Digest()), IntegrityChecksum, "blockchain schema has been altered")

	t.Run("Deploy Registry chain", func(t *testing.T) {
		entry, _ := blockchain.DeployRegistry(a)
		assert.NotNil(t, entry)
		entry, _ = blockchain.Metachain().Publish(a)
	})

	t.Run("Publish Schemata", func(t *testing.T) {
		entry, _ := b.Publish(a)
		assert.NotNil(t, entry)
	})

	t.Run("Deploy Chain", func(t *testing.T) {
		entry, _ := b.Deploy(a)
		assert.Equal(t, entry.ChainID, b.ChainID)

		t.Run("Integrity Check", func(t *testing.T) {
			assert.Equal(t, x.Decode(entry.Content), x.Decode(b.Digest()))
		})
	})

	t.Run("Add Entry", func(t *testing.T) {
		ts := x.Encode(fmt.Sprintf("%v", time.Now().Unix()))
		body := x.Encode(fmt.Sprintf("hello@%v", time.Now().Unix()))
		extids := [][]byte{ts}
		entry, _ := b.Commit(a, extids, body)
		assert.Equal(t, entry.ChainID, b.ChainID)
		assert.True(t, blockchain.ValidSignature(entry))
	})

	t.Run("Add Offer", func(t *testing.T) {

		declaration := finite.OptionContract()

		entry, _ := b.Offer(declaration, a)
		assert.True(t, blockchain.ValidSignature(entry))
		assert.True(t, blockchain.ValidContract(entry))

		t.Run("Execute Command", func(t *testing.T) {

			// make commits and test for expected error outcome
			commit := func(action string, a *identity.Account) (*factom.Entry, error) {

				pub := identity.PublicKey{}
				copy(pub[:], x.PrivateKeyToPub(a.Priv.Key[:]))

				cmd := contract.Command{
					ChainID:    b.ChainID,
					ContractID: declaration.ContractID,
					Schema:     ptnet.OptionV1,
					Action:     action, // state machine action
					Mult:       1,      // triggers input action 'n' times
					Payload:    nil,    // arbitrary data optionally included
					Pubkey:     pub,
				}

				e, err := b.Execute(cmd, u1)
				assert.True(t, blockchain.ValidSignature(e))
				return e, err
			}

			e, _ := commit("OPT_1", u1)
			assert.False(t, "" == e.ChainID)

		})
	})

}
