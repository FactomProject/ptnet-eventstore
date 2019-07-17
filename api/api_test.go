package api_test

import (
	"context"
	"encoding/json"
	"github.com/FactomProject/ptnet-eventstore/event"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/FactomProject/ptnet-eventstore/api"
)

var ctx, _ = context.WithTimeout(context.Background(), time.Second)

func Serve() {
	// set path to example schemata
	pflowPath, _ := filepath.Abs("../examples")
	os.Setenv("PFLOWPATH", pflowPath)

	go api.Serve()
	time.Sleep(time.Microsecond * 1)
}

func TestCounterEvents(t *testing.T) {
	Serve()
	defer api.Stop()

	c := api.TestClient()
	oid := event.NewUuid().String()
	schema := "counter"

	commit := func(action string, multiple int, testFail bool) { // dispatch command
		s, err := c.Dispatch(ctx, schema, oid, map[string]uint64{action: uint64(multiple)}, map[string]string{"foo": "bar"}, []uint64{})

		if s == nil || err != nil && testFail != true {
			t.Fatalf("API call failed %v", err)
			t.Logf("FAIL: %v, %v, %v, %v, %v ,%v", s.Id, s.Schema, s.Head, s.State, s.Updated, s.Created)
		} else {
			if err != nil {
				t.Logf("ROLLBACK: %s", err)
			} else {
				t.Logf("COMMIT: %v, %v, %v, %v, %v ,%v", s.Id, s.Schema, s.Head, s.State, s.Updated, s.Created)
			}
		}
	}

	xPass := func(action string, multiple int) { // dispatch command
		commit(action, multiple, false)
	}

	xFail := func(action string, multiple int) { // dispatch command
		commit(action, multiple, true)
	}

	ping := func() {
		ok := c.Ping(ctx)
		if !ok {
			t.Fatal("Failed to ping server")
		}
	}

	state := func() {
		s, err := c.GetState(ctx, schema, oid, "")

		if s == nil || err != nil {
			t.Fatalf("API call failed %v", err)
		}

		for _, state := range s {
			t.Logf("state: %v, %v, %v, %v", state.Schema, state.Id, state.Head, state.State)
		}
	}

	event := func() {
		el, err := c.GetEvent(ctx, schema, oid, "")

		if el == nil || err != nil {
			t.Fatalf("API call failed %v", err)
		}

		for _, e := range el {
			t.Logf("evt: %v, %v, %v, %v, %v ", e.Id, e.Schema, e.Uuid, e.State, e.Parent)
		}
	}

	machine := func(schema string) {
		s, err := c.GetMachine(ctx, schema, oid, "")

		if s == nil || err != nil {
			t.Fatalf("API call failed %v", err)
		}

		j, _ := json.Marshal(s)
		t.Logf("machine: %s\n", j)
	}

	machines := func() {
		s, err := c.ListMachines(ctx)

		if s == nil || err != nil {
			t.Fatalf("API call failed %v", err)
		}

		j, _ := json.Marshal(s)
		t.Logf("machines: %s\n", j)
	}

	ping()

	xPass("INC0", 1)
	xPass("INC0", 1)
	xPass("INC1", 2)

	xPass("INC1.INC0.INC0", 1) // compound action

	xFail("DEC1", 3)
	xFail("FAKE1", 1)

	event()
	state()
	machine("counter")
	machine("octoe")
	machine("octoe2step")
	machines()
}
