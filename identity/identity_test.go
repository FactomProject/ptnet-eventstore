package identity_test

import (
	. "github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeyValidation(t *testing.T) {
	pub := Public[DEPOSITOR]
	addr := Address[DEPOSITOR]
	//fmt.Printf("addr: %x pub: %x", addr, pub)
	assert.True(t, pub.MatchesAddress(addr))
}