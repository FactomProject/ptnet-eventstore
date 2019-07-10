package event

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strings"
	"time"
)

func NewStateVector() pq.Int64Array {
	return pq.Int64Array{}
}

const SuperUser = "SuperUser"

type Event struct {
	Id       uuid.UUID
	Schema   string
	Action   string
	Multiple uint64
	Payload  json.RawMessage
	State    pq.Int64Array
	TS       time.Time
	Uuid     uuid.UUID
	Parent   uuid.UUID
}

type State struct {
	Id      uuid.UUID
	Schema  string
	State   pq.Int64Array
	Head    uuid.UUID
	Created time.Time
	Updated time.Time
}

func NewUuid() uuid.UUID {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return id
}

func ParseUuid(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		panic(err)
	}
	return id
}

// fails with a panic
func NewEvent(id string, schema string, action []string, multiple uint64, payload interface{}) *Event {
	e, err := newEvent(id, schema, action, multiple, payload)
	if err != nil {
		panic(err)
	}
	return e
}

// return error if conversion fails
func PrepareEvent(id string, schema string, action []string, multiple uint64, payload interface{}) (*Event, error) {
	return newEvent(id, schema, action, multiple, payload)
}

// for empty or truncated inputs
const emptyJsonErrorMessage = "json: error calling MarshalJSON for type json.RawMessage: unexpected end of JSON input"

func newEvent(id string, schema string, action []string, multiple uint64, payload interface{}) (*Event, error) {
	if len(action) == 0 {
		panic("empty action")
	}

	if payload == nil {
		panic("empty payload")
	}

	j, err := json.Marshal(payload)
	if err != nil && err.Error() == emptyJsonErrorMessage {
		// set payload as json encoded empty map
		j, err = json.Marshal(make(map[string]string))
	}

	var oid uuid.UUID
	oid, err = uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return &Event{
		Id:       oid,
		Schema:   schema,
		Action:   strings.Join(action, "."),
		Multiple: multiple,
		Payload:  j,
		State:    nil,
		TS:       time.Now(),
	}, nil
}

// convert db-typed to state-machine-typed state vectors
func PqArrayToUint(ar pq.Int64Array) []uint64 {
	a := make([]uint64, len(ar))
	for i, r := range ar {
		a[i] = uint64(r)
	}
	return a
}

func (evt Event) String() string {
	return fmt.Sprintf("(%s, %s) %v(%v) => %v, %s", evt.Id, evt.Uuid, evt.Action, evt.Multiple, evt.State, evt.Payload)
}

func (s State) String() string {
	return fmt.Sprintf("(%s, %s) %v => %v", s.Id, s.Head, s.Schema, s.State)
}
