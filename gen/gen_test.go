package gen_test

import (
	"github.com/FactomProject/ptnet-eventstore/gen"
	"github.com/stackdump/gopetri/statemachine"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestGenerateStateMachineCode(t *testing.T) {
	p, s := statemachine.LoadPnmlFromFile("../pnml/counter.xml")
	src := statemachine.Generate(p, "CounterV1")
	ioutil.WriteFile("counter.go", src.Bytes(), 0600)
	println(src.String())
	assert.Nil(t, nil)
	_ = s
}

func TestUsingGeneratedSource(t *testing.T) {
	m := gen.CounterV1.StateMachine()
	m.Init()
	_, err := m.Transform("INC_0", 1)
	assert.Nil(t, err)

}
