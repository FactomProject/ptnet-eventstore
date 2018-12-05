package gen

import (
	. "github.com/stackdump/gopetri/statemachine"
)

var OctoeV1 PetriNet = PetriNet{
	Places: map[string]Place { 
		"OUT_O": Place{
				Initial: 0,
				Offset: 0,
				Capacity: 0,
		},
		"OUT_X": Place{
				Initial: 0,
				Offset: 1,
				Capacity: 0,
		},
		"REFUND": Place{
				Initial: 0,
				Offset: 2,
				Capacity: 0,
		},
		"m00": Place{
				Initial: 1,
				Offset: 3,
				Capacity: 0,
		},
		"m01": Place{
				Initial: 1,
				Offset: 4,
				Capacity: 0,
		},
		"m02": Place{
				Initial: 1,
				Offset: 5,
				Capacity: 0,
		},
		"m10": Place{
				Initial: 1,
				Offset: 6,
				Capacity: 0,
		},
		"m11": Place{
				Initial: 1,
				Offset: 7,
				Capacity: 0,
		},
		"m12": Place{
				Initial: 1,
				Offset: 8,
				Capacity: 0,
		},
		"m20": Place{
				Initial: 1,
				Offset: 9,
				Capacity: 0,
		},
		"m21": Place{
				Initial: 1,
				Offset: 10,
				Capacity: 0,
		},
		"m22": Place{
				Initial: 1,
				Offset: 11,
				Capacity: 0,
		},
		"new": Place{
				Initial: 1,
				Offset: 12,
				Capacity: 0,
		},
		"turn_o": Place{
				Initial: 0,
				Offset: 13,
				Capacity: 0,
		},
		"turn_x": Place{
				Initial: 0,
				Offset: 14,
				Capacity: 0,
		},
	},
	Transitions: map[Action]Transition { 
		"END_O": Transition{ 0,0,1,0,0,0,0,0,0,0,0,0,0,-1,0 },
		"END_X": Transition{ 0,0,1,0,0,0,0,0,0,0,0,0,0,0,-1 },
		"EXEC": Transition{ 0,0,0,0,0,0,0,0,0,0,0,0,-1,0,1 },
		"O00": Transition{ 0,0,0,-1,0,0,0,0,0,0,0,0,0,-1,1 },
		"O01": Transition{ 0,0,0,0,-1,0,0,0,0,0,0,0,0,-1,1 },
		"O02": Transition{ 0,0,0,0,0,-1,0,0,0,0,0,0,0,-1,1 },
		"O10": Transition{ 0,0,0,0,0,0,-1,0,0,0,0,0,0,-1,1 },
		"O11": Transition{ 0,0,0,0,0,0,0,-1,0,0,0,0,0,-1,1 },
		"O12": Transition{ 0,0,0,0,0,0,0,0,-1,0,0,0,0,-1,1 },
		"O20": Transition{ 0,0,0,0,0,0,0,0,0,-1,0,0,0,-1,1 },
		"O21": Transition{ 0,0,0,0,0,0,0,0,0,0,-1,0,0,-1,1 },
		"O22": Transition{ 0,0,0,0,0,0,0,0,0,0,0,-1,0,-1,1 },
		"WIN_O": Transition{ 1,0,0,0,0,0,0,0,0,0,0,0,0,0,-1 },
		"WIN_X": Transition{ 0,1,0,0,0,0,0,0,0,0,0,0,0,-1,0 },
		"X00": Transition{ 0,0,0,-1,0,0,0,0,0,0,0,0,0,1,-1 },
		"X01": Transition{ 0,0,0,0,-1,0,0,0,0,0,0,0,0,1,-1 },
		"X02": Transition{ 0,0,0,0,0,-1,0,0,0,0,0,0,0,1,-1 },
		"X10": Transition{ 0,0,0,0,0,0,-1,0,0,0,0,0,0,1,-1 },
		"X11": Transition{ 0,0,0,0,0,0,0,-1,0,0,0,0,0,1,-1 },
		"X12": Transition{ 0,0,0,0,0,0,0,0,-1,0,0,0,0,1,-1 },
		"X20": Transition{ 0,0,0,0,0,0,0,0,0,-1,0,0,0,1,-1 },
		"X21": Transition{ 0,0,0,0,0,0,0,0,0,0,-1,0,0,1,-1 },
		"X22": Transition{ 0,0,0,0,0,0,0,0,0,0,0,-1,0,1,-1 },
	},
}
