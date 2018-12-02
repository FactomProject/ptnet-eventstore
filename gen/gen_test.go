package gen_test

import (
	"github.com/FactomProject/ptnet-eventstore/gen"
	"github.com/stackdump/gopetri/statemachine"
	"io/ioutil"
	"testing"
)

func TestGenerateStateMachineCode(t *testing.T) {
	p, s := statemachine.LoadPnmlFromFile("../pnml/counter.xml")
	src := gen.Generate(p, "CounterV1")
	ioutil.WriteFile("counter.go", src.Bytes(), 0600)
	println(src.String())
	_ = s
}

func TestUsingNew(t *testing.T) {
}
