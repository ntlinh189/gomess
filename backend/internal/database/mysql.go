package database

import (
	"database/sql"
	"gomess/internal/config"
)

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