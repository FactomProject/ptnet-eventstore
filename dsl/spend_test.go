package dsl_test

import (
	. "github.com/FactomProject/ptnet-eventstore/dsl"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSpendDSL(t *testing.T) {

	b := NewBlockchain("SpendChain")

	t.Run("Prepare Chains", func(t *testing.T) {

		actor := b.GetAccount("BANK")

		t.Run("Create Chain for SmartContract Registry", func(t *testing.T) {
			registry, _ := DeployRegistry(actor)
			println(registry.String())
		})

		t.Run("Publish Blockchain Specification to Registry", func(t *testing.T) {
			schema, _ := b.Publish(actor)
			println(schema.String())
		})

		t.Run("Create New Chain", func(t *testing.T) {
			chain, _ := b.Deploy(actor)
			println(chain.String())
		})
	})

	t.Run("Use Contract", func(t *testing.T) {

		c := Spend.Contract(ID("Foo", "bar"), 0,
			Inputs(
				FROM(Address("USER1"), 1, Token1),
			),
			Outputs(
				TO(Address("USER2"), 1, Token1),
			),
		)
		println(c.String())

		t.Run("Publish Contract Offer", func(t *testing.T) {
			assert.False(t, Spend.Exists(c.ContractID), "contract should not exist yet")

			oc, _ := Spend.Offer(b, c, b.GetAccount("USER1"))

			assert.True(t, Spend.Exists(c.ContractID), "contract should now exist")
			println(oc.String())
		})


		t.Run("Evaluate Contract State", func(t *testing.T) {
			assert.True(t, IsHalted(c), "contract is only considered complete when it is halted")
			assert.True(t, CanRedeem(c, b.GetAccount("USER2")))
			assert.False(t, CanRedeem(c, b.GetAccount("USER1")), "Only USER2 can redeem output")
		})

	})
}
