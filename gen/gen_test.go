package gen_test

import (
	"github.com/FactomProject/ptnet-eventstore/gen"
	"github.com/stackdump/gopetri/statemachine"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func genStateMachine(path string, varName string, filename string) statemachine.PetriNet {
	p, _ := statemachine.LoadPnmlFromFile(path)
	src := statemachine.Generate(p, varName)
	ioutil.WriteFile(filename, src.Bytes(), 0600)
	//println(src.String())
	return p
}

func TestGenerateStateMachines(t *testing.T) {
	// NOTE: overwrites source files
	//genStateMachine("../pnml/counter.xml", "CounterV1", "counter.go")
	//genStateMachine("../pnml/option.xml", "OptionV1", "option.go")
	genStateMachine("../pnml/octoe.xml", "OctoeV1", "octoe.go")
}

func TestUsingGeneratedSource(t *testing.T) {
	m := gen.CounterV1.StateMachine()
	m.Init()
	_, err := m.Transform("INC_0", 1)
	assert.Nil(t, err)
	assert.NotNil(t, gen.OptionV1)
	assert.NotNil(t, gen.OctoeV1)
}
