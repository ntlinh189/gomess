package main

import (
	"gomess/internal/app"
	"gomess/internal/config"
	"log"
)

func main() {
	cfg := config.NewConfig()

	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
