// @title						GoMess API
// @version					1.0
// @description				Backend API for GoMess
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
// @description				Example: Bearer eyJhbGciOiJIUzI1NiIs...
package main

import (
	_ "gomess/docs"
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
