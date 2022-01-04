package handlers

import (
	"encoding/json"
	. "github.com/MonikaPalova/tarikatomobil.bg/db"
	"github.com/MonikaPalova/tarikatomobil.bg/httputils"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type UsersHandler struct {
	DB *UsersDBHandler
}

func (u UsersHandler) Get(w http.ResponseWriter, r *http.Request) {
	users, err := u.DB.GetUsers()
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not fetch users from DB", err)
		return
	}

	bytes, err := json.Marshal(users)
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

	id, err := uuid.NewUUID()
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not generate a user ID", err)
		return
	}
	userToCreate.ID = id.String()
	userToCreate.TimesDriver = 0
	userToCreate.TimesPassenger = 0

	if err := u.DB.CreateUser(userToCreate); err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not create user", err)
		return
	}

	if err := json.NewEncoder(w).Encode(&userToCreate); err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not marshal created user", err)
		return
	}

	log.Printf("Successfully created a user with ID %s", userToCreate.ID)
}
