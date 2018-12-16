// Package sim runs the Factomd blockchain simulator to test prototype Factom Asset Token implementations
package sim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/ptnet-eventstore/contract"
	"github.com/FactomProject/ptnet-eventstore/gen"
	"github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"github.com/FactomProject/ptnet-eventstore/x"
	"text/template"
)

const Meta string = "Meta"

type Color uint8

const (
	Default = Color(iota)
)

type Token struct {
	Color Color
}

type Identity struct {
}

type Account struct {
	Identity
}

type Blockchain struct {
	ChainID string
	ExtIDs  [][]byte
	Tokens  []Token
	Account
	Contracts map[string]contract.Contract
}

func NewBlockchain(extids ...string) (*Blockchain, error) {
	ext := [][]byte{}
	for _, id := range extids {
		ext = append(ext, x.Encode(id))
	}

	b := Blockchain{
		ChainID:   x.NewChainID(ext),
		ExtIDs:    ext,
		Tokens:    []Token{ { Color: Default } },
		Contracts: contract.Contracts,
	}

	return &b, nil
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
	err := Dispatch(commit, reveal)
	return &e, err
}

func (b *Blockchain) Commit(a *identity.Account, extids [][]byte, content []byte) (*factom.Entry, error) {
	e := x.Entry(b.ChainID, extids, content)
	commit, _ := x.ComposeCommitEntryMsg(a.Priv, e)
	reveal, _ := x.ComposeRevealEntryMsg(a.Priv, &e)
	err := Dispatch(commit, reveal)
	return &e, err
}

func Ext(extIDs ...string) [][]byte {
	ext := [][]byte{}

	for _, id := range extIDs {
		ext = append(ext, x.Encode(id))
	}

	return ext
}


// definition for deploying the registry chain
func Metachain() *Blockchain{

	ext := Ext(Meta, ptnet.FiniteV1)

	b := Blockchain{
		ChainID:   x.NewChainID(ext),
		ExtIDs:    ext,
		Tokens:    []Token{ { Color: Default } },
		Contracts: map[string]contract.Contract{
				ptnet.FiniteV1: contract.Contract{
					Schema:   ptnet.FiniteV1,
					Machine:  gen.FiniteV1.StateMachine(),
					Template: contract.Declaration{
						BlockHeight: 0,
						ContractID:	x.NewChainID(ext),
						Schema: ptnet.FiniteV1,
						Capacity: gen.FiniteV1.GetCapacityVector(),
						State: gen.FiniteV1.GetCapacityVector(),
						Actions: gen.FiniteV1.Transitions,
					},
				},
		},
	}

	return &b
}

func DeployRegistry (a *identity.Account) (*factom.Entry, error) {
	return Metachain().Deploy(a)
}

func (b *Blockchain) Publish(a *identity.Account) (*factom.Entry, error) {
	data , _ := json.Marshal(b)
	extIDs := b.ExtIDs
	extIDs = append(extIDs, x.Encode(b.ChainID))
	extIDs = append(extIDs, b.Digest())
	return Metachain().Commit(a, extIDs, data)
}

func (b *Blockchain) Search(q map[string]string) {}
func (b *Blockchain) Offer() {}
func (b *Blockchain) Execute() {}

func (b Token) Balance() {}
func (b Account) List()  {}

func (b Account) GetAccount(name string) *identity.Account {
	return identity.GetAccount(name)
}

