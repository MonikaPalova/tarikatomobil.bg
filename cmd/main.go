package main

import (
	"log"

	"github.com/MonikaPalova/tarikatomobil.bg/bootstrap"
)

func main() {
	app := bootstrap.NewApplication()
	if err := app.Start(); err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}
}
