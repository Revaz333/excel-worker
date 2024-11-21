package models

import "database/sql"

type Types struct {
	Field   string
	Type    string
	Null    sql.NullString
	Key     sql.NullString
	Default sql.NullString
	Extra   sql.NullString
}
