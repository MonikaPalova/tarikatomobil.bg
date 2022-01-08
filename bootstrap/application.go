package bootstrap

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/MonikaPalova/tarikatomobil.bg/auth"
	. "github.com/MonikaPalova/tarikatomobil.bg/db"
	. "github.com/MonikaPalova/tarikatomobil.bg/handlers"
)

type Application struct {
	db     *Database
	router *mux.Router
}

func NewApplication() Application {
	var a Application
	a.db = initDBConnection()
	a.router = createHTTPRouter(a.db)
	return a
}

func (a Application) Start() error {
	return http.ListenAndServe(":80", a.router)
}

func initDBConnection() *Database {
	db, err := InitDB("root", "", "tarikatomobil")
	if err != nil {
		log.Fatalf("Could not connect to DB: %s", err.Error())
	}
	log.Printf("Successfuly established database connection")
	return db
}

func createHTTPRouter(db *Database) *mux.Router {
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

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui")))

	return router
}
