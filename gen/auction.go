package gen

import (
	. "github.com/stackdump/gopetri/statemachine"
)

var AuctionV1 PetriNet = PetriNet{
	Places: map[string]Place { 
		"ACCEPTED": Place{
				Initial: 0,
				Offset: 0,
				Capacity: 1,
		},
		"NEW": Place{
				Initial: 1,
				Offset: 1,
				Capacity: 1,
		},
		"OPEN": Place{
				Initial: 0,
				Offset: 2,
				Capacity: 1,
		},
		"PRICE": Place{
				Initial: 0,
				Offset: 3,
				Capacity: 0,
		},
		"REJECTED": Place{
				Initial: 0,
				Offset: 4,
				Capacity: 1,
		},
	},
	Transitions: map[Action]Transition { 
		"BID": Transition{ 0,0,0,1,0 },
		"EXEC": Transition{ 0,-1,1,0,0 },
		"HALT": Transition{ 0,0,-1,0,1 },
		"SOLD": Transition{ 1,0,-1,0,0 },
	},
}
