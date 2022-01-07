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
	DB       *UsersDBHandler
	PhotosDB *PhotoDBHandler
}

func (u UsersHandler) Get(w http.ResponseWriter, r *http.Request) {
	username, ok := mux.Vars(r)["name"]
	if !ok { // Should not happen
		httputils.RespondWithError(w, http.StatusInternalServerError, "Mux did not forward the request correctly", nil)
		return
	}

	user, dbErr := u.DB.GetUserByName(username)
	if dbErr != nil {
		httputils.RespondWithDBError(w, dbErr)
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
		httputils.RespondWithDBError(w, dbErr)
		return
	}

	if err := json.NewEncoder(w).Encode(&userToCreate); err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not marshal created user", err)
		return
	}

	log.Printf("Successfully created a user with name %s", userToCreate.Name)
}

func (u UsersHandler) Patch(w http.ResponseWriter, r *http.Request) {
	username, ok := mux.Vars(r)["name"]
	if !ok { // Should not happen
		httputils.RespondWithError(w, http.StatusInternalServerError, "Mux did not forward the request correctly", nil)
		return
	}

	var userPatch model.UserPatch
	if err := json.NewDecoder(r.Body).Decode(&userPatch); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Could not parse request body as JSON", err)
		return
	}

	if err := userPatch.ValidateNonEmptyUserData(); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Could not update the user with these fields", err)
		return
	}

	// Get user so that we can see if his photo id was updated
	userBeforeUpdate, dbErr := u.DB.GetUserByName(username)
	if dbErr != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError,
			fmt.Sprintf("Could not get user %s", username), dbErr.Err)
		return
	}

	if dbErr = u.DB.UpdateUser(username, userPatch); dbErr != nil {
		httputils.RespondWithDBError(w, dbErr)
		return
	}

	if userPatch.PhotoID != "" && // a photo ID was provided
		userBeforeUpdate.PhotoID != userPatch.PhotoID && // the photo ID is different than the previous photo ID
		userBeforeUpdate.PhotoID != model.DefaultPhotoID { // the old photo is not the default photo
		// Then delete the previous photo
		dbErr = u.PhotosDB.DeletePhoto(userBeforeUpdate.PhotoID)
		if dbErr != nil {
			// Just log the error, the patch was successful
			log.Printf("User %s's photo was updated from %s to %s. Could not delete the old photo: %s",
				username, userBeforeUpdate.PhotoID, userPatch.PhotoID, dbErr.Err.Error())
		} else {
			log.Printf("User %s's photo was updated from %s to %s. Successfully deleted photo %s",
				username, userBeforeUpdate.PhotoID, userPatch.PhotoID, userBeforeUpdate.PhotoID)
		}
	}

	_, _ = w.Write([]byte(fmt.Sprintf("Successfully updated user %s", username)))
}
