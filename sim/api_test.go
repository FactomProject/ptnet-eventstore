package sim_test

import (
	"encoding/json"
	"fmt"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/sim"
	"github.com/FactomProject/ptnet-eventstore/x"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)


func TestBlockchainApi(t *testing.T) {
	b, _ := sim.NewBlockchain("salt", ptnet.OctoeV1, ptnet.OptionV1)

	d := b.Tokens[0]
	d.Balance()
	assert.Equal(t, d.Color, sim.Default)

	a := b.GetAccount("BANK")
	println(b.String())

	assert.Equal(t, b.ChainID, "c2487b6e6b7ae49aa813700205735dc40f7a2fef314eb5fcc28d564b217c58f3")

	t.Run("Deploy Chain", func(t *testing.T) {
		entry, err := b.Deploy(a)
		assert.Equal(t, entry.ChainID, b.ChainID)
		assert.NotNil(t, err, "expected error when sim is not running")
		println(entry.String())

		t.Run("Integrity Check", func(t *testing.T) {
			expectedSig := "b3d91ebc6494f76c2a368a755a897df27c2e4209c378f2c52cdd7b36d1fde61d"
			assert.Equal(t, x.Decode(entry.Content), expectedSig)
			assert.Equal(t, x.Decode(b.Digest()), expectedSig)
		})
	})

	t.Run("Add Entry", func(t *testing.T) {
		ts := x.Encode(fmt.Sprintf("%v", time.Now().Unix()))
		body :=x.Encode(fmt.Sprintf("hello@%v", time.Now().Unix()))
		extids := [][]byte{ts}
		entry, err := b.Commit(a, extids, body)
		assert.Equal(t, entry.ChainID, b.ChainID)
		assert.NotNil(t, err, "expected error when sim is not running")
	})

	t.Run("Publish Contract", func(t *testing.T) {
		data, _ := json.Marshal(b)
		fmt.Printf("%s", data)
	})

}
