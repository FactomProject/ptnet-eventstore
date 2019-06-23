package eventstore

import (
	"fmt"
	"github.com/FactomProject/ptnet-eventstore/storage"
	"github.com/stackdump/gopflow/ptnet"
	"github.com/stackdump/gopflow/statemachine"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func NewEventStore() (es *EventStore) {
	es = &EventStore{}
	es.db = storage.Reconnect()
	es.m = StateMachines()

	return es
}

func (es *EventStore) Migrate() {
	tx, err := es.db.Begin()

	if err != nil {
		panic(err)
	}

	for schema := range es.m {
		storage.Migrate(tx, schema)
	}

	tx.Commit()
}

// load all pflow files as event schemata
func StateMachines() map[string]*statemachine.StateMachine {
	m := make(map[string]*statemachine.StateMachine)

	pflowPath, ok := os.LookupEnv("PFLOWPATH")
	if !ok {
		pflowPath = "./"
	}
	//fmt.Printf("%v\n", pflowPath)

	files, err := ioutil.ReadDir(pflowPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		match, err := regexp.MatchString(`(.*)\.pflow$`, f.Name())
		if err != nil {
			panic(err)
		}

		if match {
			name := strings.Replace(f.Name(), ".pflow", "", 1)
			path := fmt.Sprintf("%s/%s.pflow", pflowPath, name)
			m[name] = ptnet.LoadFile(path).StateMachine()
		}
	}
	return m
}
