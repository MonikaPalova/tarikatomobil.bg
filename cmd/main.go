package main

import (
	"github.com/MonikaPalova/tarikatomobil.bg/auth"
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

	usersHandler := UsersHandler{DB: db.UsersDBHandler, PhotosDB: db.PhotosDBHandler}
	photosHandler := PhotoHandler{DB: db.PhotosDBHandler}
	loginHandler := LoginHandler{
		UserDB:    db.UsersDBHandler,
		SessionDB: db.SessionDBHandler,
	}
	sessionAuthMiddleware := auth.SessionAuthMiddleware{DB: db.SessionDBHandler}

	authRouter := router.NewRoute().Subrouter()
	authRouter.Use(sessionAuthMiddleware.Middleware)

	router.Path("/login").Methods(http.MethodPost).HandlerFunc(loginHandler.Login)

	router.Path("/users").Methods(http.MethodPost).HandlerFunc(usersHandler.Post)
	router.Path("/users/{name}").Methods(http.MethodGet).HandlerFunc(usersHandler.Get)
	authRouter.Path("/users/{name}").Methods(http.MethodPatch).HandlerFunc(usersHandler.Patch)

	authRouter.Path("/photos").Methods(http.MethodPost).HandlerFunc(photosHandler.UploadPhoto)
	router.Path("/photos/{id}").Methods(http.MethodGet).HandlerFunc(photosHandler.GetPhoto)
	authRouter.Path("/photos/{id}").Methods(http.MethodDelete).HandlerFunc(photosHandler.DeletePhoto)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui")))

	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatalf("Could not start server: %s", err.Error())
	}
}
