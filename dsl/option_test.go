package dsl_test

import (
	. "github.com/FactomProject/ptnet-eventstore/dsl"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptionDSL(t *testing.T) {

	b := NewBlockchain("foo", "option")
	actor := b.GetAccount("BANK")

	t.Run("Prepare Chains", func (t *testing.T) {

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

	t.Run("Use Contract", func (t *testing.T) {

		c := Option.Contract( ID("Foo", "bar"), 0,
			Inputs(
				FROM(Address("DEPOSITOR"), 1, Token0),
			),
			Outputs(
				TO(Address("DEPOSITOR"), 1, Token0),
				TO(Address("USER1"), 1, Token0),
				TO(Address("USER2"), 1, Token0),
			),
		)
		println(c.String())

		t.Run("Publish Contract Offer", func (t *testing.T) {
			assert.False(t, Option.Exists(c.ContractID), "contract should not exist yet")

			oc, _ := Option.Offer(b, c, b.GetAccount("DEPOSITOR"))

			assert.True(t, Option.Exists(c.ContractID), "contract should now exist")
			println(oc.String())
		})

		t.Run("Accept Offer", func (t *testing.T) {
			ex, _ := Option.Execute(b, c, Action("OPT_1"), b.GetAccount("USER1"))
			println(ex.String())
		})

		t.Run("Evaluate Contract State", func (t *testing.T) {
			assert.True(t, IsHalted(c), "contract is only considered complete when it is halted")
			assert.True(t, CanRedeem(c, b.GetAccount("USER1")), "Only USER1 can redeem output")

			assert.False(t, CanRedeem(c, b.GetAccount("USER2")))
			assert.False(t, CanRedeem(c, b.GetAccount("DEPOSITOR")))
		})

	})
}
