package main

import (
	"log"
	"os"

	"github.com/oren/standard-package-layout/http"
	"github.com/oren/standard-package-layout/postgres"
)

func main() {
	// Connect to database.
	db, err := postgres.Open(os.Getenv("DB"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create services.
	us := &postgres.UserService{DB: db}

	// Attach to HTTP handler.
	var h http.Handler
	h.UserService = us

	// start http server...
}
