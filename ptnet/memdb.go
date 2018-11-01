package ptnet

import (
	"github.com/hashicorp/go-memdb"
)

const StateTable = "states"
const EventTable = "events"

var eventStore *memdb.DBSchema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		EventTable: &memdb.TableSchema{
			Name: EventTable,
			Indexes: map[string]*memdb.IndexSchema{
				"Oid": &memdb.IndexSchema{
					Name:   "Oid",
					Unique: false,
					Indexer: &memdb.StringFieldIndex{
						Field: "Oid",
					},
				},
				"id": &memdb.IndexSchema{
					Name:   "id",
					Unique: true,
					Indexer: &memdb.CompoundIndex{
						Indexes: []memdb.Indexer{
							&memdb.UintFieldIndex{
								Field: "Timestamp",
							},
							&memdb.StringFieldIndex{
								Field: "Oid",
							},
						},
					},
				},
			},
		},
		StateTable: &memdb.TableSchema{
			Name: StateTable,
			Indexes: map[string]*memdb.IndexSchema{
				"id": &memdb.IndexSchema{
					Name:   "id",
					Unique: true,
					Indexer: &memdb.StringFieldIndex{
						Field: "Oid",
					},
				},
			},
		},
	},
}

func EventStore() *memdb.MemDB {
	db, _ := memdb.NewMemDB(eventStore)
	return db
}
