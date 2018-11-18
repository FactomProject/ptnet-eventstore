package x

import (
	"crypto/sha256"
	"github.com/FactomProject/ed25519"
	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/testHelper"
)

var FactoidPrefix = []byte{0x5f, 0xb1}
var FactoidPrivatePrefix = []byte{0x64, 0x78}

var NewSignature = factoid.NewED25519Signature

var NewPrivateKey = testHelper.NewPrivKey
var PrivateKeyToPub = testHelper.PrivateKeyToEDPub

var Sha = primitives.Sha
var VerifyCanonical = ed25519.VerifyCanonical
var Sign = ed25519.Sign

func Shad(data []byte) []byte {
	h1 := sha256.Sum256(data)
	h2 := sha256.Sum256(h1[:])
	return h2[:]
}

// FIXME: copy FCT address sting funcs into this package
var ConvertAddressToUser       = primitives.ConvertAddressToUser
var ConvertFctAddressToUserStr = primitives.ConvertFctAddressToUserStr
var ConvertFctPrivateToUserStr = primitives.ConvertFctPrivateToUserStr
var ValidateFUserStr           = primitives.ValidateFUserStr
var ValidateFPrivateUserStr    = primitives.ValidateFPrivateUserStr
var ConvertUserStrToAddress    = primitives.ConvertUserStrToAddress
