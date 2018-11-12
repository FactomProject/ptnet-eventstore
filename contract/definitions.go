package contract

import "github.com/FactomProject/ptnet-eventstore/ptnet"
import . "github.com/FactomProject/ptnet-eventstore/identity"

/*
 contract fixtures (currently for testing) are defined in this file
 eventually these structures will be stored on chain and indexed in memory when in use
*/

// define a deposit contract that makes an offer to 2 identities
// provides a way to revoke payment by original sender
// will only redeem once
// tokens should be thought of as "pay-to-script" - locked until state machine is halted
func OptionContract() Declaration {
	m := ptnet.StateMachines["option-v1"]

	return Declaration{ // array of inputs also referenced by guards and conditions
		Inputs: []AddressAmountMap{ // array of input depositors
			AddressAmountMap{DEPOSITOR, 1}, // deposit tokens
		},
		Outputs: []AddressAmountMap{
			AddressAmountMap{DEPOSITOR, 1}, // withdraw token
			AddressAmountMap{USER1, 1},     // deposit to user1
			AddressAmountMap{USER2, 1},     // deposit to user2
		},
		BlockHeight: 60221409,       // deadline for halting state
		Salt:        "|RANDOM|",     // added random salt
		ContractID:  CONTRACT_ID,    // unique ID for this contract instance
		Schema:      ptnet.OptionV1, // versioned contract schema
		State:       m.Initial,      // state machine initial state
		Actions:     m.Transitions,  // state machine defined transitions
		Guards: []Condition{ // guard clause restricts actions
			Condition{0, 0, 0, -1, 0}, // block unless contract is still open
			Condition{0, 0, 0, -1, 0}, // NOTE: don't really need this but it illustrates ability
			Condition{0, 0, 0, -1, 0}, // to restrict actions without the state machine being halted
		},
		Conditions: []Condition{ // contract conditions specify additional redeem conditions
			Condition{0, 0, 0, 0, -1}, // refund pay addr[0]
			Condition{0, -1, 0, 0, 0}, // pay addr[1]
			Condition{0, 0, -1, 0, 0}, // pay addr[2]
		},
	}
}

func TicTacToeContract() Declaration {
	m := ptnet.StateMachines["octoe-v1"]

	return Declaration{ // array of inputs/outputs also referenced by guards and conditions
		Inputs: []AddressAmountMap{ // array of input depositors
			AddressAmountMap{DEPOSITOR, 1},
		},
		Outputs: []AddressAmountMap{ // array of possible redeemers
			AddressAmountMap{DEPOSITOR, 1},
			AddressAmountMap{PLAYERX, 1},
			AddressAmountMap{PLAYERO, 1},
		},
		BlockHeight: 60221409,      // deadline for halting state
		Salt:        "|RANDOM|",    // added random salt
		ContractID:  CONTRACT_ID,   // unique ID for this contract instance
		Schema:      "octoe-v1",    // versioned contract schema
		State:       m.Initial,     // state machine initial state
		Actions:     m.Transitions, // state machine defined transitions
		Guards: []Condition{ // guard clause restricts actions
			//       00 01 02 10 11 12 20 21 22  O  X $O $X $DEP // variable labels
			Condition{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},  // Admin - 'contract owner' can take action at any time
			Condition{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0, 0, 0}, // PlayerX - players must move only when it's their turn
			Condition{0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0, 0, 0, 0}, // PlayerO
		},
		Conditions: []Condition{ // contract conditions specify additional redeem conditions
			//       00 01 02 10 11 12 20 21 22  O  X $O $X $DEP  // variable labels
			Condition{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1}, // game ended without winner tokens are unlocked for original depositor
			Condition{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0}, // game ended PlayerX wins
			Condition{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0, 0}, // game ended PlayerO wins
		},
	}
}
