package gen

import (
	. "github.com/stackdump/gopetri/statemachine"
)

var FiniteV1 PetriNet = PetriNet{
	Places: map[string]Place { 
		"active": Place{
				Initial: 1,
				Offset: 0,
				Capacity: 0,
		},
		"edits": Place{
				Initial: 0,
				Offset: 1,
				Capacity: 0,
		},
		"inactive": Place{
				Initial: 0,
				Offset: 2,
				Capacity: 0,
		},
	},
	Transitions: map[Action]Transition { 
		"DISABLE": Transition{ -1,0,1 },
		"ENABLE": Transition{ 1,0,-1 },
		"EXEC": Transition{ 0,1,0 },
	},
}
