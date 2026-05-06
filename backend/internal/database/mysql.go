package database

import (
	"database/sql"
	"gomess/internal/config"
	"log"
)

func NewMySQL(cfg *config.Config) *sql.DB {
	db, err := sql.Open("mysql", cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}