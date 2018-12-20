package gen

import (
	. "github.com/stackdump/gopetri/statemachine"
)

var Spend PetriNet = PetriNet{
	Places: map[string]Place { 
		"HALTED": Place{
				Initial: 0,
				Offset: 1,
				Capacity: 1,
		},
		"PAYMENT": Place{
				Initial: 0,
				Offset: 0,
				Capacity: 0,
		},
	},
	Transitions: map[Action]Transition { 
		"EXEC": Transition{ 1,0 },
	},
}
