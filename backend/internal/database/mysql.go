package database

import (
	"database/sql"
	"gomess/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseInterface interface {
	GetDB() *sql.DB
}

type Database struct {
	db *sql.DB
}

func NewMySql(cfg config.ConfigInterface) (*Database, error) {
	db, err := sql.Open("mysql", cfg.GetDBUrl())

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}