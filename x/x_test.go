package x_test

import (
	"github.com/FactomProject/ptnet-eventstore/x"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFactomdSigning(t *testing.T) {
	priv := x.NewPrivateKey(0)
	pub := x.PrivateKeyToPub(priv)
	testData := x.Sha([]byte("sig first half  one")).Bytes()

	pubFixed := [32]byte{}
	copy(pubFixed[:], pub)

	privFixed := [64]byte{}
	copy(privFixed[:], append(priv, pub...)[:])

	t.Run("Sign and Validate", func(t *testing.T) {
		sig := x.NewSignature(priv, testData)
		assert.True(t, x.VerifyCanonical(&pubFixed, testData, &sig.Signature))

		sig2 := x.Sign(&privFixed, testData)
		assert.True(t, x.VerifyCanonical(&pubFixed, testData, sig2))
	})
}
