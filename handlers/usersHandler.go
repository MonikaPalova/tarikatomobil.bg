package handlers

import (
	"encoding/json"
	"fmt"
	. "github.com/MonikaPalova/tarikatomobil.bg/db"
	"github.com/MonikaPalova/tarikatomobil.bg/httputils"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type UsersHandler struct {
	DB *UsersDBHandler
}

func (u UsersHandler) Get(w http.ResponseWriter, r *http.Request) {
	username, ok := mux.Vars(r)["name"]
	if !ok { // Should not happen
		httputils.RespondWithError(w, http.StatusInternalServerError, "Mux did not forward the request correctly", nil)
		return
	}

	user, dbErr := u.DB.GetUserByName(username)
	if dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, fmt.Sprintf("Could not get user with name %s", username))
		return
	}

	bytes, err := json.Marshal(user)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not marshal users: %s", err)
		return
	}
	_, _ = w.Write(bytes)
}

func (u UsersHandler) Post(w http.ResponseWriter, r *http.Request) {
	var userToCreate model.User
	if err := json.NewDecoder(r.Body).Decode(&userToCreate); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Could not parse request body as JSON", err)
		return
	}

	if err := userToCreate.ValidateUserData(); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Could not create user with these fields", err)
		return
	}

	if userToCreate.PhotoID == "" {
		userToCreate.PhotoID = model.DefaultPhotoID
	}
	userToCreate.TimesDriver = 0
	userToCreate.TimesPassenger = 0

	if dbErr := u.DB.CreateUser(userToCreate); dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, "Could not create user")
		return
	}

	if err := json.NewEncoder(w).Encode(&userToCreate); err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not marshal created user", err)
		return
	}

	log.Printf("Successfully created a user with name %s", userToCreate.Name)
}
