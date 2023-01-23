package models

import "database/sql"

type WorkerQueues struct {
	ID        int
	Driver    string
	FilePath  string
	TableName string
	Status    sql.NullString
	Error     sql.NullString
	Columns   string
}
