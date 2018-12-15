package x

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/FactomProject/ed25519"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/engine"
	"github.com/FactomProject/factomd/state"
	"github.com/FactomProject/factomd/testHelper"
	"time"
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
var ConvertAddressToUser = primitives.ConvertAddressToUser
var ConvertFctAddressToUserStr = primitives.ConvertFctAddressToUserStr
var ConvertFctPrivateToUserStr = primitives.ConvertFctPrivateToUserStr
var ValidateFUserStr = primitives.ValidateFUserStr
var ValidateFPrivateUserStr = primitives.ValidateFPrivateUserStr

var ComposeCommitEntryMsg = testHelper.ComposeCommitEntryMsg
var ComposeRevealEntryMsg = testHelper.ComposeRevealEntryMsg
var ComposeChainCommit = testHelper.ComposeChainCommit

var FundWallet = engine.FundWallet
var FundECWallet = engine.FundECWallet

var AccountFromFctSecret = testHelper.AccountFromFctSecret
var GetRandomAccount = testHelper.GetRandomAccount()
var GetBankAccount = testHelper.GetBankAccount // funded by genesis block

var GetBalanceEC = engine.GetBalanceEC
var GetBalance = engine.GetBalance

var ShutDownEverything = testHelper.ShutDownEverything
var WaitForAllNodes = testHelper.WaitForAllNodes
var WaitBlocks = testHelper.WaitBlocks
var SetupSim = testHelper.SetupSim


func WaitForEcBalance(s *state.State, ecPub string) int64 {
	var bal int64 = 0

	for {
		bal = GetBalanceEC(s, ecPub)
		time.Sleep(time.Millisecond * 200)
		//fmt.Printf("WaitForBalance: %v => %v\n", ecPub, bal)

		if bal > 0 {
			return bal
		}
	}
}

var NewChain = factom.NewChain

func Entry(chainID string, extIDs  [][]byte, content []byte) factom.Entry {
	return factom.Entry{
		ChainID: chainID,
		ExtIDs:  extIDs,
		Content: content,

	}
}

func Encode(s string) []byte {
	b := bytes.Buffer{}
	b.WriteString(s)
	return b.Bytes()
}

func Decode(b []byte) string {
	return string(b)
}

// create the chainid from a series of hashes of the Entries ExtIDs
func NewChainID(extIDs [][]byte) string {
	hs := sha256.New()
	for _, id := range extIDs {
		h := sha256.Sum256(id)
		hs.Write(h[:])
	}
	return hex.EncodeToString(hs.Sum(nil))
}
