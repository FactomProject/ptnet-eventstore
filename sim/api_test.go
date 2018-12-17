package sim_test

import (
	"fmt"
	"github.com/FactomProject/ptnet-eventstore/blockchain"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)


func TestBlockchainApi(t *testing.T) {
	b := blockchain.NewBlockchain("Merged", ptnet.OctoeV1, ptnet.OptionV1)

	d := b.Tokens[0]
	d.Balance()
	assert.Equal(t, d.Color, blockchain.Default)

	a := b.GetAccount("BANK")
	println(b.String())

	//assert.Equal(t, b.ChainID, "c2487b6e6b7ae49aa813700205735dc40f7a2fef314eb5fcc28d564b217c58f3")

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

}
