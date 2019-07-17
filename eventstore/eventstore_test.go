package eventstore_test

import (
	"context"
	"github.com/FactomProject/ptnet-eventstore/event"
	"github.com/FactomProject/ptnet-eventstore/eventstore"
	"github.com/stackdump/gopflow/statemachine"
	"os"
	"path/filepath"
	"testing"
)

func TestCounter(t *testing.T) {
	// set path to example schemata
	pflowPath, _ := filepath.Abs("../examples")
	os.Setenv("PFLOWPATH", pflowPath)

	eventstore.Drop()
	es := eventstore.NewEventStore()
	es.Migrate()

	var schema = "counter"

	oid := event.NewUuid()
	ctx := context.Background()
	history := make([]*event.Event, 0)

	// REVIEW: you must bring your own roles
	// it is assumed role affiliation will be validated w/ DID/Verifiable Credentials
	roles := map[statemachine.Role]bool{
		"SuperUser": true,
	}

	commit := func(action string, multiple int, payload interface{}, assertError bool) {
		evt := event.NewEvent(oid.String(), schema, map[string]uint64{action: uint64(multiple)}, payload)
		st, err := es.Commit(context.WithValue(ctx, "roles", roles), evt)
		_ = st

		//t.Logf("%v\n", st.String())
		switch {
		case err != nil && !assertError:
			{
				t.Fatal(err)
			}
		case err != nil:
			{
				t.Logf("Found Error %v %v", err, evt)
			}
		case err == nil:
			{
				history = append(history, evt)
				if assertError {
					t.Fatal("Expected Error")
				}
			}
		}

	}

	// expect action to fail
	xFail := func(action string, multiple int, payload interface{}) {
		commit(action, multiple, payload, true)

	}
	// expect action to pass
	xPass := func(action string, multiple int, payload interface{}) {
		commit(action, multiple, payload, false)
	}

	xPass("INC0", 1, map[string]string{"hello": "world"})
	if history[0] != nil {
		if len(history[0].State) == 0 {
			t.Fatalf("Failed to persist state")
		}
	}

	xPass("INC1", 2, map[string]string{"hello": "again"})
	if history[1].State[0] != 1 || history[1].State[1] != 2 {
		t.Fatal("Failed to persist state")
	}

	// test parent relation
	if history[0].Uuid != history[1].Parent {
		t.Logf("Parent: %v", history[1].Parent)
		t.Fatalf("Failed to set parent event uuid")
	}

	// trigger invalid output -1
	xFail("DEC0", 3, map[string]string{"hello": "failure"})

	// compound action
	xPass("INC0.INC1.INC1.INC1", 1, map[string]string{"hello": "compound"})

	// compound action mixed values
	xPass("INC0(3).DEC0(2).INC1.INC1.INC1", 1, map[string]string{"hello": "muli-value"})

	t.Log("GetEvents")
	for _, evt := range es.GetEvents("counter", oid.String()) {
		t.Logf(evt.String())
	}

	t.Log("GetEvent")
	t.Log(es.GetEvent("counter", history[0].Uuid.String()).String())

	t.Log("GetState")
	t.Log(es.GetState("counter", history[0].Id.String()).String())

}
