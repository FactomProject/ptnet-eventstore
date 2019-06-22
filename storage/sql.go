package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"os"
	"strings"
)

const STATE = "_state"
const EVENT = "_event"

func EventTable(schema string) string {
	return pq.QuoteIdentifier(schema + EVENT)
}

func StateTable(schema string) string {
	return pq.QuoteIdentifier(schema + STATE)
}

func Reconnect() *sql.DB {
	pgpass, ok := os.LookupEnv("PGPASS")
	if !ok {
		pgpass = "pflow"
	}

	pguser, ok := os.LookupEnv("PGUSER")
	if !ok {
		pguser = "pflow"
	}

	pgdatabase, ok := os.LookupEnv("PGDATABASE")
	if !ok {
		pgdatabase = "pflow"
	}

	pghost, ok := os.LookupEnv("PGHOST")
	if !ok {
		pghost = "127.0.0.1"
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", pguser, pgpass, pgdatabase, pghost)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return db
}

const createEvents = `
CREATE TABLE IF NOT EXISTS %s (
    id UUID,
    schema VARCHAR DEFAULT '',
    action VARCHAR,
    multiple INT,
    payload JSONB NOT NULL,
    state INT[],
    ts TIMESTAMP DEFAULT now(),
    uuid UUID,
    parent UUID,
    PRIMARY KEY (id, schema, uuid)
)
`

const createStates = `
CREATE TABLE IF NOT EXISTS %s (
    id UUID,
    schema VARCHAR DEFAULT '',
    state INT[],
    head UUID,
    created TIMESTAMP default now(),
    updated TIMESTAMP ,
    PRIMARY KEY (id, schema)
)
`

// create event and state table for given schema
func Migrate(tx *sql.Tx, schema string) {
	//fmt.Printf("migrate: %v\n", schema)
	var err error

	_, err = tx.Exec(fmt.Sprintf(createEvents, EventTable(schema)))
	if err != nil { panic(err) }

	_, err = tx.Exec(fmt.Sprintf(createStates, StateTable(schema)))
	if err != nil { panic(err) }
}

const dropTable = `
DROP TABLE IF EXISTS %s
`

// drop event and state table for given schema
func Drop(db *sql.DB, schema string) {
	var err error

	_, err = db.Exec(fmt.Sprintf(dropTable, pq.QuoteIdentifier(schema+STATE)))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(fmt.Sprintf(dropTable, pq.QuoteIdentifier(schema+EVENT)))
	if err != nil {
		panic(err)
	}
}

const appendEvent = `
INSERT INTO %s
    (uuid, id, schema, action, multiple, payload, state, parent, ts)
VALUES
    ('%s', '%s', '%s', '%s', %d, '%s', '%s', '%s', now())
`

func AppendEvent(schema string, oid string, eventUuid string, action string, multiple uint64, payload json.RawMessage, state pq.Int64Array, parent string) string {
	return fmt.Sprintf(appendEvent,
		EventTable(schema),
		eventUuid,
		oid,
		schema,
		action,
		multiple,
		payload,
		vectorLiteral(state),
		parent,
	)
}

const listEvents = `
SELECT * FROM %s
WHERE
    id = '%s' 
AND
    schema = '%s' 
`

func GetEvents(schema string, oid string) string {
	return fmt.Sprintf(listEvents, EventTable(schema), oid, schema)
}

const getEvent = `
SELECT * FROM %s
WHERE
    uuid = '%s'
AND
    schema = '%s' 
`

// TODO: refactor rename oid/uuid
func GetEvent(schema string, oid string) string {
	return fmt.Sprintf(getEvent, EventTable(schema), oid, schema)
}

const getState = `
SELECT * FROM %s
WHERE
    id = '%s'
AND
    schema = '%s' 
`

func GetState(schema string, oid string) string {
	return fmt.Sprintf(getState, StateTable(schema), oid, schema)
}

const setState = `
INSERT INTO %s
    (id, schema, state, head, updated)
VALUES ('%s', '%s', '%s', '%s', now())
ON CONFLICT(id, schema) DO
    UPDATE SET state = '%s', head = '%s'
WHERE
    excluded.schema = '%s' 
AND
    excluded.id = '%s' 
`

// convert to postgres literal array
func vectorLiteral(ar pq.Int64Array) string {
	raw, err := json.Marshal(ar)
	if err != nil {
		panic(err)
	}

	s := string(raw)
	s = strings.Replace(s, "[", "{", 1)
	s = strings.Replace(s, "]", "}", 1)
	return s
}

func SetState(schema string, oid string, state pq.Int64Array, head string) string {
	stateLiteral := vectorLiteral(state)
	return fmt.Sprintf(setState, StateTable(schema), oid, schema, stateLiteral, head, stateLiteral, head, schema, oid)
}
