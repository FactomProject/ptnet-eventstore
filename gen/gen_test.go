package gen_test

import (
	"github.com/FactomProject/ptnet-eventstore/gen"
	"github.com/stackdump/gopetri/statemachine"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestGenerateStateMachines(t *testing.T) {

	generate := func(path string, varName string, filename string) statemachine.PetriNet {
		p, _ := statemachine.LoadPnmlFromFile(path)
		src := statemachine.Generate(p, varName)
		ioutil.WriteFile(filename, src.Bytes(), 0600)
		return p
	}

	// NOTE: overwrites source files
	generate("../pnml/counter.xml", "CounterV1", "counter.go")
	generate("../pnml/option.xml", "OptionV1", "option.go")
	generate("../pnml/octoe.xml", "OctoeV1", "octoe.go")
	generate("../pnml/registry.xml", "FiniteV1", "registry.go")
}

func TestUsingGeneratedSource(t *testing.T) {
	m := gen.CounterV1.StateMachine()
	m.Init()
	_, err := m.Transform("INC_0", 1)
	assert.Nil(t, err)
	assert.NotNil(t, gen.OptionV1)
	assert.NotNil(t, gen.OctoeV1)
	assert.NotNil(t, gen.FiniteV1)
}
