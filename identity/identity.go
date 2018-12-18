// Package identity contains public & private keys used for testing
package identity

import (
	"bytes"
	"encoding/hex"
	"github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/ptnet-eventstore/x"
	"text/template"
)

type PrivateKey [64]byte
type PublicKey [32]byte
type FctAddress []byte

func NewPrivateKey(seed uint64) PrivateKey {
	n := x.NewPrivateKey(seed)
	k := PrivateKey{}
	copy(k[:], n)
	return k
}

func PrivateKeyToPub(priv PrivateKey) PublicKey {
	n := x.PrivateKeyToPub(priv[:])
	k := PublicKey{}
	copy(k[:], n)
	return k
}

var DEPOSITOR string = "DEPOSITOR"
var USER1 string = "USER1"
var USER2 string = "USER2"
var PLAYERX string = "PLAYERX"
var PLAYERO string = "PLAYERO"
var BANK string = "BANK"

var Private = map[string]PrivateKey{
	BANK:      PrivateKeyFromFctSecret("Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK"),
	DEPOSITOR: NewPrivateKey(10000),
	USER1:     NewPrivateKey(10001),
	USER2:     NewPrivateKey(10002),
	PLAYERX:   NewPrivateKey(10003),
	PLAYERO:   NewPrivateKey(10004),
}

var Public = map[string]PublicKey{
	BANK: PrivateKeyToPub(Private[BANK]),
	DEPOSITOR: PrivateKeyToPub(Private[DEPOSITOR]),
	USER1:     PrivateKeyToPub(Private[USER1]),
	USER2:     PrivateKeyToPub(Private[USER2]),
	PLAYERX:   PrivateKeyToPub(Private[PLAYERX]),
	PLAYERO:   PrivateKeyToPub(Private[PLAYERO]),
}
var Address = map[string]FctAddress{
	BANK:      PublicKeyToAddress(Public[BANK]),
	DEPOSITOR: PublicKeyToAddress(Public[DEPOSITOR]),
	USER1:     PublicKeyToAddress(Public[USER1]),
	USER2:     PublicKeyToAddress(Public[USER2]),
	PLAYERX:   PublicKeyToAddress(Public[PLAYERX]),
	PLAYERO:   PublicKeyToAddress(Public[PLAYERO]),
}

func PublicKeyToAddress(publicKey PublicKey) []byte {
	data := []byte{1}
	data = append(data, publicKey[:]...)
	return x.Shad(data)
}

func (p PublicKey) MatchesAddress(address FctAddress) bool {
	a := PublicKeyToAddress(p)
	for offset, c := range address {
		if a[offset] != c {
			return false
		}
	}
	return true
}

type Account struct {
	Priv *primitives.PrivateKey
}

func GetAccount(name string) *Account {
	b := []byte{}
	for _, v := range Private[name] {
		b = append(b, v)
	}
	return &Account{primitives.NewPrivateKeyFromHexBytes(b)}
}

func (d *Account) FctPriv() string {
	x, _ := primitives.PrivateKeyStringToHumanReadableFactoidPrivateKey(d.Priv.PrivateKeyString())
	return x
}

func (d *Account) FctPub() string {
	s, _ := factoid.PublicKeyStringToFactoidAddressString(d.Priv.PublicKeyString())
	return s
}

func (d *Account) EcPub() string {
	s, _ := factoid.PublicKeyStringToECAddressString(d.Priv.PublicKeyString())
	return s
}

func (d *Account) EcPriv() string {
	s, _ := primitives.PrivateKeyStringToHumanReadableECPrivateKey(d.Priv.PrivateKeyString())
	return s
}

func (d *Account) FctAddr() interfaces.IHash {
	a := primitives.ConvertUserStrToAddress(d.FctPub())
	x, _ := primitives.HexToHash(hex.EncodeToString(a))
	return x
}

func (d *Account) PrivHash() interfaces.IHash {
	a := primitives.ConvertUserStrToAddress(d.EcPriv())
	x, _ := primitives.HexToHash(hex.EncodeToString(a))
	return x
}

func (d *Account) EcAddr() interfaces.IHash {
	a := primitives.ConvertUserStrToAddress(d.EcPub())
	x, _ := primitives.HexToHash(hex.EncodeToString(a))
	return x
}

func (d *Account) GetPrivateKey() PrivateKey {
	k := PrivateKey{}
	for i, v := range d.Priv.Key {
		k[i] = v
	}
	return k
}

func (d *Account) GetPublicKey () PublicKey {
	pub := PublicKey{}
	copy(pub[:], x.PrivateKeyToPub(d.Priv.Key[:]))
	return pub
}

func PrivateKeyFromFctSecret(s string) PrivateKey {
	h, _ := primitives.HumanReadableFactoidPrivateKeyToPrivateKey(s)
	k := PrivateKey{}
	copy(k[:], h)
	return k
}

var testFormat string = `
HASH
  PrivHash: {{ .PrivHash }}
FCT
  FctPriv: {{ .FctPriv }}
  FctPub: {{ .FctPub }}
  FctAddr: {{ .FctAddr }}
EC
  EcPriv: {{ .EcPriv }}
  EcPub: {{ .EcPub }}
  EcAddr: {{ .EcAddr }}
`

var testTemplate *template.Template = template.Must(
	template.New("").Parse(testFormat),
)

func (d *Account) String() string {
	b := &bytes.Buffer{}
	testTemplate.Execute(b, d)
	return b.String()
}
