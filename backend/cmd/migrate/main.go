package main

import (
	"errors"
	"flag"
	"gomess/internal/config"
	"gomess/internal/database"
	"log"

	"github.com/golang-migrate/migrate/v4"
	mysqlmigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	direction := flag.String("direction", "up", `"up" or "down"`)
	flag.Parse()

	cfg := config.NewConfig()

	db, err := database.NewMySql(cfg)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := mysqlmigrate.WithInstance(db.GetDB(), &mysqlmigrate.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
	if err != nil {
		log.Fatal(err)
	}

	switch *direction {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	default:
		log.Fatalf("unknown direction: %s", *direction)
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	log.Println("migration completed:", *direction)
}