package ptnet

/*
declare state machines
REVIEW: consider loading definitions from file on demand

These state machines are re-usable components that can be
referenced and integrated w/ smart contract definitions

This may mean that part of this protocol is to publish versioned
state machines to a distinct revision-control chain on Factom
*/

const OptionV1 string = "option-v1" // version contract definitions by using schema name
const OctoeV1 string = "octoe-v1"   // this allows for future mechanism to 'upconvert' v1 -> v2

var StateMachines map[string]Machine = map[string]Machine{
	"counter": counterMachine,
	OctoeV1:   octoeMachine,
	OptionV1:  optionMachine,
}

// choice between 2 options with one approver
var optionMachine Machine = Machine{
	Initial: StateVector{1, 1, 1, 1, 1}, // use 'all ones' for initial vector
	Transitions: map[string]Transition{
		BEGIN:   Transition{0, -1, -1, 0, -1}, // unsure if we should introduce standard action names
		"OPT_1": Transition{-1, 1, 0, 0, 0},   // pay addr[1]
		"OPT_2": Transition{-1, 0, 1, 0, 0},   // pay addr[2]
		CANCEL:  Transition{-1, 0, 0, -1, 1},  // refund pay addr[0] - use like exit code
		END:     Transition{0, 0, 0, -1, 0},   // confirm transaction
	},
	db: EventStore(),
}

// this state machine never halts
// simple counter structure used mostly for
// testing & benchmarking state machine and eventstore
var counterMachine Machine = Machine{
	Initial: StateVector{1, 1},
	Transitions: map[string]Transition{
		BEGIN:   Transition{-1, -1},
		"INC_0": Transition{1, 0},
		"DEC_0": Transition{-1, 1},
		"INC_1": Transition{0, 1},
		"DEC_1": Transition{0, -1},
	},
	db: EventStore(),
}

var octoeMachine Machine = Machine{
	//                   00 01 02 10 11 12 20 21 22 O  X $O $X $DEP
	Initial: StateVector{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},

	Transitions: map[string]Transition{

		//                00 01 02 10 11 12 20 21 22  O  X $O $X $DEP
		BEGIN: Transition{0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0, -1, -1, -1},

		"X00": Transition{-1, 0, 0, 0, 0, 0, 0, 0, 0, 1, -1, 0, 0, 0},
		"X01": Transition{0, -1, 0, 0, 0, 0, 0, 0, 0, 1, -1, 0, 0, 0},
		"X02": Transition{0, 0, -1, 0, 0, 0, 0, 0, 0, 1, -1, 0, 0, 0},
		"X10": Transition{0, 0, 0, -1, 0, 0, 0, 0, 0, 1, -1, 0, 0, 0},
		"X11": Transition{0, 0, 0, 0, -1, 0, 0, 0, 0, 1, -1, 0, 0, 0},
		"X12": Transition{0, 0, 0, 0, 0, -1, 0, 0, 0, 1, -1, 0, 0, 0},
		"X20": Transition{0, 0, 0, 0, 0, 0, -1, 0, 0, 1, -1, 0, 0, 0},
		"X21": Transition{0, 0, 0, 0, 0, 0, 0, -1, 0, 1, -1, 0, 0, 0},
		"X22": Transition{0, 0, 0, 0, 0, 0, 0, 0, -1, 1, -1, 0, 0, 0},

		"O00": Transition{-1, 0, 0, 0, 0, 0, 0, 0, 0, -1, 1, 0, 0, 0},
		"O01": Transition{0, -1, 0, 0, 0, 0, 0, 0, 0, -1, 1, 0, 0, 0},
		"O02": Transition{0, 0, -1, 0, 0, 0, 0, 0, 0, -1, 1, 0, 0, 0},
		"O10": Transition{0, 0, 0, -1, 0, 0, 0, 0, 0, -1, 1, 0, 0, 0},
		"O11": Transition{0, 0, 0, 0, -1, 0, 0, 0, 0, -1, 1, 0, 0, 0},
		"O12": Transition{0, 0, 0, 0, 0, -1, 0, 0, 0, -1, 1, 0, 0, 0},
		"O20": Transition{0, 0, 0, 0, 0, 0, -1, 0, 0, -1, 1, 0, 0, 0},
		"O21": Transition{0, 0, 0, 0, 0, 0, 0, -1, 0, -1, 1, 0, 0, 0},
		"O22": Transition{0, 0, 0, 0, 0, 0, 0, 0, -1, -1, 1, 0, 0, 0},

		"WINX": Transition{0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0, 0, 1, 0}, // Depositor acts as judge
		"WINO": Transition{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 1, 0, 0},
		"ENDX": Transition{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0, 0, 1},
		"ENDO": Transition{0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0, 0, 0, 1},
	},
	db: EventStore(),
}
