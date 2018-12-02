package gen

import (
	"bytes"
	"fmt"
	"github.com/stackdump/gopetri/statemachine"
	"strings"
	"text/template"
)

var sourceFileFormat string = `package gen

import (
	. "github.com/stackdump/gopetri/statemachine"
)

var {{ .Name }} PetriNet = PetriNet{
	Places: map[string]Place { {{ range $key, $place := .Places}}
		"{{ $key }}": Place{
				Initial: {{ $place.Initial }},
				Offset: {{ $place.Initial }},
				Capacity: {{ $place.Capacity }},
		},{{ end }}
	},
	Transitions: map[Action]Transition { {{ range $action, $transition := .TransitionLiterals}}
		"{{ $action }}": Transition{ {{ $transition }} },{{ end }}
	},
}
`

func (s SourceFile) TransitionLiterals() map[string]string {
	out := map[string]string{}
	for action, t := range s.Transitions {
		convert := []string{}
		for _, v := range t {
			convert = append(convert, fmt.Sprintf("%v", v))
		}
		out[string(action)] = strings.Join(convert, ",")
	}
	return out
}

type SourceFile struct {
	statemachine.PetriNet
	Name string
}

var sourceTemplate *template.Template = template.Must(
	template.New("").Parse(sourceFileFormat),
)

func Generate(net statemachine.PetriNet, filename string) *bytes.Buffer {
	f := SourceFile{net, filename}
	b := &bytes.Buffer{}
	err := sourceTemplate.Execute(b, f)
	if nil != err {
		panic(err)
	}
	return b
}
