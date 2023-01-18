package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host       string
	Port       string
	Dbname     string
	Dbuser     string
	Dbpassword string
}

var db *sql.DB

func NewDb(config Config) (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		config.Dbuser, config.Dbpassword, config.Host, config.Port, config.Dbname)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}
