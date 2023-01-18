package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Host       string
	Port       string
	Dbname     string
	Dbuser     string
	Dbpassword string
}

var db *gorm.DB

func NewDb(config Config) (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Dbuser, config.Dbpassword, config.Host, config.Port, config.Dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
