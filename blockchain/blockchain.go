// Package sim runs the Factomd blockchain simulator to test prototype Factom Asset Token implementations
package blockchain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/ptnet-eventstore/contract"
	"github.com/FactomProject/ptnet-eventstore/finite"
	"github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/sim"
	"github.com/FactomProject/ptnet-eventstore/x"
	"text/template"
)

type Color uint8

const (
	Default = Color(iota)
)

type Token struct {
	Color Color
}

type Blockchain struct {
	ChainID   string `json:"chainid"`
	ExtIDs    [][]byte `json:"extids"`
	Tokens    []Token `json:"tokens"`
	Contracts map[string]contract.Contract `json:"contracts"`
}

func NewBlockchain(extids ...string) *Blockchain {
	ext := [][]byte{}
	for _, id := range extids {
		ext = append(ext, x.Encode(id))
	}

	b := Blockchain{
		ChainID:   x.NewChainID(ext),
		ExtIDs:    ext,
		Tokens:    []Token{{Color: Default}},
		Contracts: contract.Contracts,
	}

	return &b
}

var txtFormat string = `
ChainID: {{ .ChainID }}
Digest: {{ printf "%s" .Digest }}
{{ range $_, $i := .ExtIDStr}}ExtID: {{ printf "%v" $i}}
{{ end }} Tokens: {{ range $_, $token := .Tokens}}
    Color: {{ printf "%x" $token.Color }} {{ end }}
Contracts: {{ range $_, $c := .Contracts}}
    {{ $c.Schema }}: {{ printf "%x" $c.Version }} {{ end }}
`

func (b *Blockchain) ExtIDStr() (out []string) {
	for _, v := range b.ExtIDs {
		out = append(out, x.Decode(v))
	}
	return out
}

var txtTemplate *template.Template = template.Must(
	template.New("").Parse(txtFormat),
)

func (b *Blockchain) String() string {
	f := &bytes.Buffer{}
	txtTemplate.Execute(f, b)
	return f.String()
}

func (b *Blockchain) Digest() []byte {
	data, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return x.Encode(fmt.Sprintf("%x", x.Shad(data)))
}

// use blockchain spec to create new factom chain
func (b *Blockchain) Deploy(a *identity.Account) (*factom.Entry, error) {
	e := x.Entry(b.ChainID, b.ExtIDs, b.Digest())
	c := x.NewChain(&e)
	commit, _ := x.ComposeChainCommit(a.Priv, c)
	reveal, _ := x.ComposeRevealEntryMsg(a.Priv, c.FirstEntry)
	err := sim.Dispatch(commit, reveal)
	return &e, err
}

// create a new signed entry on factom
func (b *Blockchain) Commit(a *identity.Account, extids [][]byte, content []byte) (*factom.Entry, error) {
	e := AppendSignature(x.Entry(b.ChainID, extids, content), a)
	commit, _ := x.ComposeCommitEntryMsg(a.Priv, *e)
	reveal, _ := x.ComposeRevealEntryMsg(a.Priv, e)
	err := sim.Dispatch(commit, reveal)
	return e, err
}

// blockchain def used to publish blockchain defs
func Metachain() *Blockchain {
	ext := x.Ext(ptnet.FiniteV1, ptnet.Meta)
	b := Blockchain{
		ChainID: x.NewChainID(ext),
		ExtIDs:  ext,
		Tokens:  []Token{{Color: Default}},
		Contracts: map[string]contract.Contract{
			ptnet.FiniteV1: contract.Contracts[ptnet.FiniteV1],
		},
	}

	return &b
}

// create registry chain
func DeployRegistry(a *identity.Account) (*factom.Entry, error) {
	return Metachain().Deploy(a)
}

// publish blockchain definition to factom 'metachain' registry
func (b *Blockchain) Publish(a *identity.Account) (*factom.Entry, error) {
	offer := finite.Registry()
	offer.ContractID = b.ChainID

	data, _ := json.Marshal(b)

	t := finite.OfferTransaction(offer, a.GetPrivateKey())
	t.Payload = data
	t.AddDigest()

	in, _ := json.Marshal(t.InputState)
	out, _ := json.Marshal(t.OutputState)

	extids := x.Ext(t.Schema, t.Action, t.Oid, fmt.Sprintf("%v", t.Mult))
	extids = append(extids, in)
	extids = append(extids, out)

	return Metachain().Commit(a, extids, t.Payload)
}

func (b Blockchain) GetAccount(name string) *identity.Account {
	return identity.GetAccount(name)
}

func (b *Blockchain) Offer(offer finite.Offer, a *identity.Account) (*factom.Entry, error) {
	t := finite.OfferTransaction(offer, a.GetPrivateKey())
	extids := x.Ext(t.Schema, t.Action, t.Oid, fmt.Sprintf("%v", t.Mult))

	in, _ := json.Marshal(t.InputState)
	extids = append(extids, in)

	out, _ := json.Marshal(t.OutputState)
	extids = append(extids, out)

	t.Payload, _ = json.Marshal(offer.Variables)

	return b.Commit(a, extids, t.Payload)
}

func (b *Blockchain) Execute(cmd contract.Command, a *identity.Account) (*factom.Entry, error) {
	_, ok := b.Contracts[cmd.Schema]

	if !ok {
		msg := fmt.Sprintf("Undefined Contract: %v", cmd.Schema)
		return new(factom.Entry), errors.New(msg)
	}

	t, err := finite.ExecuteTransaction(cmd, a.GetPrivateKey())
	if err != nil {
		return new(factom.Entry), err
	}

	extids := x.Ext(cmd.Schema, cmd.Action, cmd.ContractID, fmt.Sprintf("%v", cmd.Mult))

	in, _ := json.Marshal(t.InputState)
	extids = append(extids, in)

	out, _ := json.Marshal(t.OutputState)
	extids = append(extids, out)

	return b.Commit(a, extids, t.Payload)
}

// add signature to extIDs
func AppendSignature(entry factom.Entry, a *identity.Account) *factom.Entry{
	e := factom.Entry{ entry.ChainID, entry.ExtIDs, entry.Content }
	s := a.Priv.Sign(e.Hash())
	key := a.Priv.Pub[:]
	keyString := x.EncodeToString(key)

	sig := x.Encode(fmt.Sprintf("%x", s.Bytes()))
	e.ExtIDs = append(e.ExtIDs, x.Encode(keyString))
	e.ExtIDs = append(e.ExtIDs, sig)
	return &e
}

// validate appended signatures
func ValidSignature(entry *factom.Entry)  bool {
	l := len(entry.ExtIDs)
	key := entry.ExtIDs[l-2]
	signature := entry.ExtIDs[l-1]
	e := factom.Entry{ entry.ChainID, entry.ExtIDs[:l-2], entry.Content }

	pub := new([32]byte)
	x.HexDecode(pub[:], key)

	sig := new([64]byte)
	x.HexDecode(sig[:], signature)

	return x.VerifyCanonical(pub, e.Hash(), sig)
}

func ValidContract(entry *factom.Entry)  bool {
	_, ok := contract.Contracts[x.Decode(entry.ExtIDs[0])]

	if ! ok ||  ! ValidSignature(entry) {
		return false
	}

	v := new(contract.Variables)
	err := json.Unmarshal(entry.Content, v)

	if err != nil || v.ContractID == "" {
		return false
	}

	return true
}
