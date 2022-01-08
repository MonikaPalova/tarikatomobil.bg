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
	db         *Database
	router     *mux.Router
	authRouter *mux.Router
}

func NewApplication() Application {
	var a Application
	a.initDBConnection()
	a.setupHTTPRoutes()
	return a
}

func (a *Application) Start() error {
	return http.ListenAndServe(":80", a.router)
}

func (a *Application) initDBConnection() {
	var err error
	a.db, err = InitDB("root", "", "tarikatomobil")
	if err != nil {
		log.Fatalf("Could not connect to DB: %s", err.Error())
	}
	log.Printf("Successfuly established database connection")
}

func (a *Application) setupHTTPRoutes() {
	a.router = mux.NewRouter()
	sessionAuthMiddleware := auth.SessionAuthMiddleware{DB: a.db.SessionDBHandler}
	a.authRouter = a.router.NewRoute().Subrouter()
	a.authRouter.Use(sessionAuthMiddleware.Middleware)

	a.setupPhotoHandler()
	a.setupLoginHandler()
	a.setupUserHandler()
	a.setupAutomobileHandler()

	// Serve UI files
	a.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui")))
}

func (a *Application) setupPhotoHandler() {
	photosHandler := PhotoHandler{DB: a.db.PhotosDBHandler}
	a.authRouter.Path("/photos").Methods(http.MethodPost).HandlerFunc(photosHandler.UploadPhoto)
	a.router.Path("/photos/{id}").Methods(http.MethodGet).HandlerFunc(photosHandler.GetPhoto)
}

func (a *Application) setupLoginHandler() {
	loginHandler := LoginHandler{
		UserDB:    a.db.UsersDBHandler,
		SessionDB: a.db.SessionDBHandler,
	}
	a.router.Path("/login").Methods(http.MethodPost).HandlerFunc(loginHandler.Login)
}

func (a *Application) setupUserHandler() {
	usersHandler := UsersHandler{DB: a.db.UsersDBHandler, PhotosDB: a.db.PhotosDBHandler}
	a.router.Path("/users").Methods(http.MethodPost).HandlerFunc(usersHandler.Post)
	a.router.Path("/users/{name}").Methods(http.MethodGet).HandlerFunc(usersHandler.Get)
	a.authRouter.Path("/users/{name}").Methods(http.MethodPatch).HandlerFunc(usersHandler.Patch)
}

func (a *Application) setupAutomobileHandler() {
	automobilesHandler := AutomobileHandler{DB: a.db.AutomobileDBHandler}
	a.authRouter.Path("/users/{name}/automobile").Methods(http.MethodPost).HandlerFunc(automobilesHandler.Post)
	a.authRouter.Path("/users/{name}/automobile").Methods(http.MethodPatch).HandlerFunc(automobilesHandler.Patch)
	//a.router.Path("/users/{name}/automobile").Methods().HandlerFunc(automobilesHandler.)
	a.router.Path("/users/{name}/automobile").Methods(http.MethodGet).HandlerFunc(automobilesHandler.Get)
}