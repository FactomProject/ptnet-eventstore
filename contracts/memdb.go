package contracts

import (
	"github.com/hashicorp/go-memdb"
)

const ContractTable = "contracts"

var contractStore *memdb.DBSchema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		ContractTable: &memdb.TableSchema{
			Name: ContractTable,
			Indexes: map[string]*memdb.IndexSchema{
				"id": &memdb.IndexSchema{
					Name:   "id",
					Unique: true,
					Indexer: &memdb.StringFieldIndex{
						Field: "ContractID",
					},
				},
			},
		},
	},
}

func ContractStore() *memdb.MemDB {
	db, _ := memdb.NewMemDB(contractStore)
	return db
}
