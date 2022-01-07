package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/MonikaPalova/tarikatomobil.bg/auth"
	. "github.com/MonikaPalova/tarikatomobil.bg/db"
	"github.com/MonikaPalova/tarikatomobil.bg/httputils"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type UsersHandler struct {
	DB       *UsersDBHandler
	PhotosDB *PhotoDBHandler
}

func (u UsersHandler) Get(w http.ResponseWriter, r *http.Request) {
	user, dbErr := u.DB.GetUserByName(mux.Vars(r)["name"])
	if dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, "Could not get user")
		return
	}

	bytes, err := json.Marshal(user)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not marshal users: %s", err, true)
		return
	}
	_, _ = w.Write(bytes)
}

func (u UsersHandler) Post(w http.ResponseWriter, r *http.Request) {
	var userToCreate model.User
	if err := json.NewDecoder(r.Body).Decode(&userToCreate); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Could not parse request body as JSON", err, false)
		return
	}

	if err := userToCreate.ValidateUserData(); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Could not create user with these fields", err, false)
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
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not marshal created user", err, true)
		return
	}

	log.Printf("Successfully created a user with name %s", userToCreate.Name)
}

func (u UsersHandler) Patch(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["name"]
	sessionUsername := auth.GetUserFromRequest(r)
	if username != sessionUsername {
		httputils.RespondWithError(w, http.StatusUnauthorized, "Modifying other users is forbidden", nil, false)
		return
	}

	var userPatch model.UserPatch
	if err := json.NewDecoder(r.Body).Decode(&userPatch); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Could not parse request body as JSON", err, false)
		return
	}

	if err := userPatch.ValidateNonEmptyUserData(); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Could not update the user with these fields", err, false)
		return
	}

	if dbErr := u.DB.UpdateUser(username, userPatch); dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, "Could not update the user")
		return
	}
	_, _ = w.Write([]byte(fmt.Sprintf("Successfully updated user %s", username)))
}
