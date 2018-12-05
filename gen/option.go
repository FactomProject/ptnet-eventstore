package gen

import (
	. "github.com/stackdump/gopetri/statemachine"
)

var OptionV1 PetriNet = PetriNet{
	Places: map[string]Place { 
		"NEW": Place{
				Initial: 1,
				Offset: 0,
				Capacity: 1,
		},
		"OPEN": Place{
				Initial: 0,
				Offset: 1,
				Capacity: 1,
		},
		"OUT1": Place{
				Initial: 0,
				Offset: 2,
				Capacity: 1,
		},
		"OUT2": Place{
				Initial: 0,
				Offset: 3,
				Capacity: 1,
		},
		"REFUND": Place{
				Initial: 0,
				Offset: 4,
				Capacity: 1,
		},
	},
	Transitions: map[Action]Transition { 
		"CANCEL": Transition{ 0,-1,0,0,1 },
		"EXEC": Transition{ -1,1,0,0,0 },
		"HALT": Transition{ 0,-1,0,0,0 },
		"OPT_1": Transition{ 0,-1,1,0,0 },
		"OPT_2": Transition{ 0,-1,0,1,0 },
	},
}
