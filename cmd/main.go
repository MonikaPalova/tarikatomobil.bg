package main

import (
	. "github.com/MonikaPalova/tarikatomobil.bg/db"
	. "github.com/MonikaPalova/tarikatomobil.bg/handlers"
	"github.com/gorilla/mux"

	"log"
	"net/http"
)

func main() {
	db := Database{}
	if err := db.Connect("root", "", "tarikatomobil"); err != nil {
		log.Fatalf("Could not connect to DB: %s", err.Error())
	}

	log.Printf("Successfuly established database connection")

	router := mux.NewRouter()

	usersHandler := UsersHandler{DB: db}

	router.Path("/users").Methods(http.MethodGet).HandlerFunc(usersHandler.Get)
	router.Path("/users").Methods(http.MethodPost).HandlerFunc(usersHandler.Post)

	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}
}
