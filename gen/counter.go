package gen

import (
	. "github.com/stackdump/gopetri/statemachine"
)

var CounterV1 PetriNet = PetriNet{
	Places: map[string]Place { 
		"p0": Place{
				Initial: 0,
				Offset: 0,
				Capacity: 0,
		},
		"p1": Place{
				Initial: 0,
				Offset: 1,
				Capacity: 0,
		},
	},
	Transitions: map[Action]Transition { 
		"DEC_0": Transition{ -1,0 },
		"DEC_1": Transition{ 0,-1 },
		"INC_0": Transition{ 1,0 },
		"INC_1": Transition{ 0,1 },
	},
}
