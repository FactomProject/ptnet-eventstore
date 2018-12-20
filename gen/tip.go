package gen

import (
	. "github.com/stackdump/gopetri/statemachine"
)

var Tip PetriNet = PetriNet{
	Places: map[string]Place { 
		"ANTIKARMA": Place{
				Initial: 0,
				Offset: 0,
				Capacity: 0,
		},
		"CHARITY": Place{
				Initial: 0,
				Offset: 1,
				Capacity: 0,
		},
		"HALTED": Place{
				Initial: 0,
				Offset: 2,
				Capacity: 0,
		},
		"KARMA": Place{
				Initial: 0,
				Offset: 3,
				Capacity: 0,
		},
		"PAYMENT": Place{
				Initial: 0,
				Offset: 4,
				Capacity: 0,
		},
	},
	Transitions: map[Action]Transition { 
		"TIP": Transition{ 0,0,0,1,1 },
		"WARN": Transition{ 1,1,0,0,0 },
	},
}
