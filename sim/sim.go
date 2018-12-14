package sim

import (
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/state"
	"github.com/FactomProject/ptnet-eventstore/x"
	"testing"
)

var state0 *state.State

// run a single simulator node starting and block 0
func Setup(t *testing.T, maxheight int) *state.State {
	state0 := x.SetupSim("L", map[string]string{"--debuglog": ""}, maxheight, 1, 1, t)
	return state0
}

func Dispatch(e interfaces.IMsg) {
	state0.APIQueue().Enqueue(e)
}

type BlockChain struct {

}

// TODO transplant x.* helpers to add Blockchain methods
