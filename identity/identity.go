// Package identity contains public & private keys used for testing
package identity

import (
	"github.com/FactomProject/ptnet-eventstore/x"
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

var Private = map[string]PrivateKey{
	DEPOSITOR: NewPrivateKey(10000),
	USER1: NewPrivateKey(10001),
	USER2: NewPrivateKey(10002),
	PLAYERX: NewPrivateKey(10003),
	PLAYERO: NewPrivateKey(10004),
}

var Public = map[string]PublicKey{
	DEPOSITOR: PrivateKeyToPub(Private[DEPOSITOR]),
	USER1: PrivateKeyToPub(Private[USER1]),
	USER2: PrivateKeyToPub(Private[USER2]),
	PLAYERX: PrivateKeyToPub(Private[PLAYERX]),
	PLAYERO: PrivateKeyToPub(Private[PLAYERO]),
}
var Address = map[string]FctAddress{
	DEPOSITOR: PublicKeyToAddress(Public[DEPOSITOR]),
	USER1: PublicKeyToAddress(Public[USER1]),
	USER2: PublicKeyToAddress(Public[USER2]),
	PLAYERX: PublicKeyToAddress(Public[PLAYERX]),
	PLAYERO: PublicKeyToAddress(Public[PLAYERO]),
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
