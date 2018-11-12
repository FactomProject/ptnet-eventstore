package finite

import (
	"fmt"
	"github.com/FactomProject/ptnet-eventstore/contract"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
)

type Offer struct{
	contract.Declaration
	ChainID string
}
type Execution struct{
	contract.Command
}
type Transaction struct{
	*ptnet.Event
}

func Depositor(o Offer) string {
	return o.Declaration.Inputs[0].Address
}

func Executor(x Execution) string {
	return x.Pubkey
}

func OfferTransaction(o Offer, privkey string) Transaction {
	event, err := contract.Create(o.Declaration, o.ChainID, func(evt *ptnet.Event) error {
		sig := fmt.Sprintf("signed with: %v", privkey)
		contract.SignEvent(evt, Depositor(o), sig)
		return nil
	})

	if err != nil {
		panic("failed to create contract")
	}

	return Transaction{Event: event}
}

func ExecuteTransaction(x Execution, privkey string) Transaction {
	event, err := contract.Transform(x.Command, func(evt *ptnet.Event) error {
		sig := fmt.Sprintf("signed with: %v", privkey)
		contract.SignEvent(evt, Executor(x), sig)
		return nil
	})

	if err != nil {
		panic("failed to create contract")
	}

	return Transaction{Event: event}
}

// FIXME actually do signing
func fooSignMe() {

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
}
