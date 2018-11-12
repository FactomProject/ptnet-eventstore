package x

import (
	"github.com/FactomProject/ed25519"
	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/testHelper"
)

var NewED25519Signature = factoid.NewED25519Signature
var Sha = primitives.Sha
var VerifyCanonical = ed25519.VerifyCanonical
var Sign = ed25519.Sign
var NewPrivKey = testHelper.NewPrivKey
var PrivateKeyToEDPub = testHelper.PrivateKeyToEDPub
