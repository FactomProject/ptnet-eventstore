package finite

import (
	"bytes"
	"github.com/FactomProject/ptnet-eventstore/contract"
	"github.com/FactomProject/ptnet-eventstore/identity"
	"github.com/FactomProject/ptnet-eventstore/ptnet"
	"text/template"
)

type Offer struct {
	contract.Declaration
	ChainID string
}

type Transaction struct {
	*ptnet.Event
}

func OfferTransaction(o Offer, privKey identity.PrivateKey) Transaction {
	event, err := contract.Create(o.Declaration, o.ChainID, privKey)

	if err != nil {
		panic("failed to create contract")
	}

	return Transaction{Event: event}
}

func ExecuteTransaction(x contract.Command, privKey identity.PrivateKey) (Transaction, error) {
	event, err := contract.Transform(x, func(evt *ptnet.Event) error {
		contract.SignEvent(evt, privKey)
		return nil
	})

	return Transaction{Event: event}, err
}

var offerFormat string = `
ChainID: {{.ChainID}}
Inputs: {{ range $_, $input := .Inputs}}
Address: {{ printf "%x" $input.Address }} Amount: {{ $input.Amount }} Color: {{ $input.Token }} {{ end }}
Outputs: {{ range $_, $output := .Outputs}}
	Address: {{ printf "%x" $output.Address }} Amount: {{ $output.Amount }} Color {{ $output.Token }} {{ end }}
BlockHeight: {{ .BlockHeight }}
ContractID: {{ .ContractID }}
Schema: {{ .Schema }}
State: {{ .GetState }}
Actions: {{ range $key, $action := .Actions }}
	{{$key}}: {{ $action }}{{ end }}
Guards: {{ range $i, $guard := .Guards }}
	{{$i}}: {{ $guard }}{{ end }}
Conditions: {{ range $i, $condition := .Conditions }}
	{{$i}}: {{ $condition }}{{ end }}
`
var offerTemplate *template.Template = template.Must(
	template.New("").Parse(offerFormat),
)

func (offer Offer) String() string {
	b := &bytes.Buffer{}
	offerTemplate.Execute(b, offerSource{offer})
	return b.String()
}

func (s offerSource) GetState() []uint64 {
	return ptnet.ToVector(s.State)
}

type offerSource struct {
	Offer
}
