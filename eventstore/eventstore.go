package eventstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/FactomProject/ptnet-eventstore/event"
	"github.com/FactomProject/ptnet-eventstore/storage"
	"github.com/stackdump/gopflow/statemachine"
	"regexp"
	"strconv"
	"strings"
)

type EventStore struct {
	db *sql.DB
	m  map[string]*statemachine.StateMachine
}

func unknownSchema(schema string) error {
	return errors.New(fmt.Sprintf("Unknown Schema %s", schema))
}

func unknownAction(action string, schema string) error {
	return errors.New(fmt.Sprintf("Unknown Action %s.%s", schema, action))
}

func (es *EventStore) GetEvent(schema, oid string) *event.Event {

	txn, errTx := es.db.Begin()
	defer txn.Rollback()
	if errTx != nil {
		panic(errTx)
	}

	rows, err := txn.Query(storage.GetEvent(schema, oid))
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		evt := &event.Event{}
		err = rows.Scan(&evt.Id, &evt.Schema, &evt.Action, &evt.Multiple, &evt.Payload, &evt.State, &evt.TS, &evt.Uuid, &evt.Parent)
		if err != nil {
			panic(err)
		}
		return evt
	}

	return nil
}

// evaluate valid events and persist to eventstore db
func (es *EventStore) Commit(ctx context.Context, evt *event.Event) (*event.State, error) {
	roles := ctx.Value("roles").(map[statemachine.Role]bool)

	var err error
	var outState []int64
	var s *event.State

	txn, err := es.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})
	if err != nil {
		panic(err)
	}

	s, err = es.getState(txn, evt.Schema, evt.Id.String())
	if err != nil {
		panic(err)
	}

	m, ok := es.m[evt.Schema]
	if !ok {
		return s, unknownSchema(evt.Schema)
	}

	outState, err = execute(m, s, evt, roles)

	if err == nil {
		eventId := event.NewUuid()
		evt.Uuid = eventId
		evt.Parent = s.Head
		evt.State = outState

		s.State = outState
		s.Head = eventId

		// persist to db
		err := es.appendEvent(txn, evt)
		if err != nil {
			panic(err)
		}

		err = es.setState(txn, s)
		if err != nil {
			panic(err)
		}

		// REVIEW: occasional error returned - doesn't seem to interfere w/ transaction
		// seems to be a warning from the driver 'unexpected tag INSERT...'
		txn.Commit()
	} else {
		txn.Rollback()
	}

	return s, err
}

func (es *EventStore) getState(txn *sql.Tx, schema string, oid string) (s *event.State, err error) {
	s = &event.State{}
	s.Id = event.ParseUuid(oid)
	s.Schema = schema

	{ // query existing state
		var rows *sql.Rows

		rows, err = txn.Query(storage.GetState(schema, oid))
		if err != nil {
			panic(err)
		}

		for rows.Next() {
			err = rows.Scan(&s.Id, &s.Schema, &s.State, &s.Head, &s.Created, &s.Updated)
			if err != nil {
				panic(err)
			}

			return s, nil
		}

		err = rows.Close()
		if err != nil {
			panic(err)
		}
	}

	sm := es.m[schema]
	if sm == nil {
		return s, unknownSchema(schema)
	}

	for _, p := range es.m[schema].Initial {
		s.State = append(s.State, int64(p))
	}
	return s, nil
}

// assure that our writes actually write
func assertRowsAffected(res sql.Result, err error) error {

	if err != nil {
		panic(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}

	if rows == 0 {
		return errors.New("failed to update")
	}

	return nil
}

// add event to DB
func (es *EventStore) appendEvent(txn *sql.Tx, e *event.Event) error {
	res, err := txn.Exec(storage.AppendEvent(e.Schema, e.Id.String(), e.Uuid.String(), e.Action, e.Multiple, e.Payload, e.State, e.Parent.String()))
	return assertRowsAffected(res, err)
}

// set state in DB
func (es *EventStore) setState(txn *sql.Tx, s *event.State) error {
	res, err := txn.Exec(storage.SetState(s.Schema, s.Id.String(), s.State, s.Head.String()))
	return assertRowsAffected(res, err)
}

// compare vectors a & b
// passes if b is empty
func matchVectorPrecondition(a []int64, b []int64) bool {
	for k, v := range b {
		if a[k] != v {
			return false
		}
	}
	return true
}

var actionRegex = regexp.MustCompile(`^(\w+)\(?(\d*)\)?`)

func execute(m *statemachine.StateMachine, s *event.State, evt *event.Event, roles map[statemachine.Role]bool) (outState []int64, err error) {

	var role statemachine.Role
	inState := event.PqArrayToUint(s.State)

	var command string
	var mult uint64

	for _, action := range strings.Split(evt.Action, ".") {

		match := actionRegex.FindStringSubmatch(action)
		if len(match) == 3 {
			if match[2] == "" {
				mult = 1
			} else {
				x, err := strconv.ParseInt(match[2], 10, 64)
				if err != nil {
					panic(err)
				}
				mult = uint64(x)
			}
			command = match[1]
		} else {
			panic("no match")
		}

		_, ok := m.Transitions[statemachine.Action(command)]
		if !ok {
			return outState, unknownAction(evt.Schema, action)
		}
		outState, role, err = m.Transform(inState, command, mult)

		switch {
		case err != nil:
			{
				return outState, err
			}
		case !roles[event.SuperUser] && !roles[role]:
			{
				return outState, errors.New(fmt.Sprintf("insufficent privledge %v", role))
			}
		default: // use output as next input
			{
				for i, v := range outState {
					inState[i] = uint64(v)
				}
			}
		}
	}

	if !matchVectorPrecondition(outState, evt.State) {
		return outState, errors.New("precondition mismatch")
	} else {
		return outState, nil
	}

}

func (es *EventStore) GetMachine(schema string, uuid string) (*statemachine.StateMachine, bool) {
	m, ok := es.m[schema]
	return m, ok
}

func (es *EventStore) ListMachines() []string {
	ml := make([]string, 0)

	for m := range es.m {
		ml = append(ml, m)
	}
	return ml
}

// REVIEW: should also support getting state by schema & head/uuid
func (es *EventStore) GetState(schema string, oid string) *event.State {
	txn, errTx := es.db.Begin()
	defer txn.Rollback()

	if errTx != nil {
		panic(errTx)
	}

	s, err := es.getState(txn, schema, oid)
	if err != nil {
		panic(err)
	}

	return s
}

func (es *EventStore) GetEvents(schema string, oid string) []*event.Event {
	res := make([]*event.Event, 0)

	txn, errTx := es.db.Begin()
	defer txn.Rollback()
	if errTx != nil {
		panic(errTx)
	}

	rows, err := txn.Query(storage.GetEvents(schema, oid))
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		evt := &event.Event{}
		err = rows.Scan(&evt.Id, &evt.Schema, &evt.Action, &evt.Multiple, &evt.Payload, &evt.State, &evt.TS, &evt.Uuid, &evt.Parent)
		if err != nil {
			panic(err)
		}
		res = append(res, evt)
	}

	return res
}

// drop all tables
func Drop() {
	db := storage.Reconnect()
	for schema := range StateMachines() {
		storage.Drop(db, schema)
	}
	db.Close()
}
