package main

import (
	. "github.com/MonikaPalova/tarikatomobil.bg/db"
	. "github.com/MonikaPalova/tarikatomobil.bg/handlers"
	"github.com/gorilla/mux"

	"log"
	"net/http"
)

func main() {
	var db *Database
	var err error

	if db, err = InitDB("root", "", "tarikatomobil"); err != nil {
		log.Fatalf("Could not connect to DB: %s", err.Error())
	}
	log.Printf("Successfuly established database connection")

	router := mux.NewRouter()

	usersHandler := UsersHandler{DB: db.UsersDBHandler}
	photosHandler := PhotoHandler{DB: db.PhotosDBHandler}

	router.Path("/users").Methods(http.MethodPost).HandlerFunc(usersHandler.Post)
	router.Path("/users/{name}").Methods(http.MethodGet).HandlerFunc(usersHandler.Get)
	router.Path("/users/{name}").Methods(http.MethodPatch).HandlerFunc(usersHandler.Patch)

	router.Path("/photos").Methods(http.MethodPost).HandlerFunc(photosHandler.UploadPhoto)
	router.Path("/photos/{id}").Methods(http.MethodGet).HandlerFunc(photosHandler.GetPhoto)
	router.Path("/photos/{id}").Methods(http.MethodDelete).HandlerFunc(photosHandler.DeletePhoto)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui")))

	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}
}
