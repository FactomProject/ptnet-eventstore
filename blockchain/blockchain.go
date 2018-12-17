// Package sim runs the Factomd blockchain simulator to test prototype Factom Asset Token implementations
package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/ptnet-eventstore/contract"
	"github.com/FactomProject/ptnet-eventstore/finite"
	"github.com/FactomProject/ptnet-eventstore/gen"
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
	ChainID string
	ExtIDs  [][]byte
	Tokens  []Token
	Contracts map[string]contract.Contract
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

func (b *Blockchain) Deploy(a *identity.Account) (*factom.Entry, error) {
	e := x.Entry(b.ChainID, b.ExtIDs, b.Digest())
	c := x.NewChain(&e)
	commit, _ := x.ComposeChainCommit(a.Priv, c)
	reveal, _ := x.ComposeRevealEntryMsg(a.Priv, c.FirstEntry)
	err := sim.Dispatch(commit, reveal)
	return &e, err
}

func (b *Blockchain) Commit(a *identity.Account, extids [][]byte, content []byte) (*factom.Entry, error) {
	e := x.Entry(b.ChainID, extids, content)
	commit, _ := x.ComposeCommitEntryMsg(a.Priv, e)
	reveal, _ := x.ComposeRevealEntryMsg(a.Priv, &e)
	err := sim.Dispatch(commit, reveal)
	return &e, err
}

func Metachain() *Blockchain {

	ext := x.Ext(ptnet.Meta, ptnet.FiniteV1)

	b := Blockchain{
		ChainID: x.NewChainID(ext),
		ExtIDs:  ext,
		Tokens:  []Token{{Color: Default}},
		Contracts: map[string]contract.Contract{
			ptnet.Meta: contract.Contract{
				Schema:  ptnet.Meta,
				Machine: gen.FiniteV1.StateMachine(),
				Template: contract.RegistryTemplate(),
			},
		},
	}

	return &b
}

func DeployRegistry(a *identity.Account) (*factom.Entry, error) {
	return Metachain().Deploy(a)
}

func (b *Blockchain) Publish(a *identity.Account) (*factom.Entry, error) {
	extIDs := b.ExtIDs
	extIDs = append(extIDs, x.Encode(b.ChainID))
	extIDs = append(extIDs, b.Digest())
	m := Metachain()

	// FIXME: construct transaction & publish
	//t := finite.OfferTransaction(finite.Registry(), a.GetPrivateKey())
	//data, _ := json.Marshal(t)
	data, _ := json.Marshal(b)
	return m.Commit(a, extIDs, data)
}

func (b Blockchain) GetAccount(name string) *identity.Account {
	return identity.GetAccount(name)
}

func Offer(offer finite.Offer, a *identity.Account) finite.Transaction {
	return finite.OfferTransaction(offer, a.GetPrivateKey())
}

// REVIEW: add methods for querying memdb + factomd entries
func (b *Blockchain) Search(q map[string]string) {}
func (b *Blockchain) Execute() {}
func (b Token) Balance() {}
