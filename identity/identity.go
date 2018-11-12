package identity

// FIXME public addresses FAxxxx
// FIXME private addresses FSxxxx

var DEPOSITOR string = "|DEPOSITOR|"
var DEPOSITOR_SECRET string = "|DEPOSITOR_SECRET|"

var USER1 string = "|USER1|"
var USER1_SECRET string = "|USER1_SECRET|"

var USER2 string = "|USER2|"
var USER2_SECRET string = "|USER2_SECRET|"

var PLAYERX string = "|PLAYERX|"
var PLAYERX_SECRET string = "|PLAYERX_SECRET|"

var PLAYERO string = "|PLAYERO|"
var PLAYERO_SECRET string = "|PLAYERO_SECRET|"

// KLUDGE: public/private keypairs for testing
var Identity map[string]string = map[string]string{
	DEPOSITOR: DEPOSITOR_SECRET,
	PLAYERX:   PLAYERX_SECRET,
	PLAYERO:   PLAYERO_SECRET,
	USER1:     USER1_SECRET,
	USER2:     USER2_SECRET,
}
