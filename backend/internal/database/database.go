package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseInterface interface {
	GetDB() *sql.DB
	WithTransaction(fn TxFunc) error
}

type Database struct {
	db *sql.DB
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}