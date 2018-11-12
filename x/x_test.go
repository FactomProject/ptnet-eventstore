package x_test

import (
	"github.com/FactomProject/ptnet-eventstore/x"
	"testing"
)

func TestFactomdSigning(t *testing.T) {
	t.Run("Validation", func(t *testing.T) {
		testData := x.Sha([]byte("sig first half  one")).Bytes()
		priv := x.NewPrivKey(0)

		sig := x.NewED25519Signature(priv, testData)

		pub := x.PrivateKeyToEDPub(priv)
		pub2 := [32]byte{}
		copy(pub2[:], pub)

		s := sig.Signature
		valid := x.VerifyCanonical(&pub2, testData, &s)
		if valid == false {
			panic("Signature is invalid")
		}

		priv2 := [64]byte{}
		copy(priv2[:], append(priv, pub...)[:])

		sig2 := x.Sign(&priv2, testData)

		valid = x.VerifyCanonical(&pub2, testData, sig2)
		if valid == false {
			panic("Test signature is invalid")
		}
	})
}
