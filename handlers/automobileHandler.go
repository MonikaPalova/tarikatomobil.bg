package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/MonikaPalova/tarikatomobil.bg/auth"
	"github.com/MonikaPalova/tarikatomobil.bg/db"
	"github.com/MonikaPalova/tarikatomobil.bg/httputils"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
)

type AutomobileHandler struct {
	DB *db.AutomobileDBHandler
}

func (ah *AutomobileHandler) Post(w http.ResponseWriter, r *http.Request) {
	caller := auth.GetUserFromRequest(r)
	usernamePathParam := mux.Vars(r)["name"]
	if caller != usernamePathParam {
		httputils.RespondWithError(w, http.StatusUnauthorized, "Creating automobiles for other users is forbidden", nil, false)
		return
	}

	var automobileToCreate model.Automobile
	if err := json.NewDecoder(r.Body).Decode(&automobileToCreate); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Could not parse request body as JSON", err, false)
		return
	}

	if automobileToCreate.PhotoID == "" {
		automobileToCreate.PhotoID = model.DefaultPhotoID
	}
	automobileToCreate.OwnerName = caller

	if dbErr := ah.DB.CreateAutomobile(automobileToCreate); dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, fmt.Sprintf("Could not create automobile for user %s", caller))
		return
	}

	log.Printf("Successfully created an automobile with reg. number %s for user %s", automobileToCreate.RegNumber, caller)
}

func (ah *AutomobileHandler) Patch(w http.ResponseWriter, r *http.Request) {
	caller := auth.GetUserFromRequest(r)
	usernamePathParam := mux.Vars(r)["name"]
	if caller != usernamePathParam {
		httputils.RespondWithError(w, http.StatusUnauthorized, "Updating automobiles for other users is forbidden", nil, false)
		return
	}

	var automobilePatch model.AutomobilePatch
	if err := json.NewDecoder(r.Body).Decode(&automobilePatch); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Could not parse request body as JSON", err, false)
		return
	}

	if dbErr := ah.DB.UpdateAutomobile(automobilePatch, caller); dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, "Could not update automobile")
		return
	}
	log.Printf("User %s's automobile was successfuly updated", caller)
}

func (ah *AutomobileHandler) Get(w http.ResponseWriter, r *http.Request) {
	usernamePathParam := mux.Vars(r)["name"]
	automobile, dbErr := ah.DB.GetUserAutomobile(usernamePathParam)
	if dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, fmt.Sprintf("Could not get user %s's automobile", usernamePathParam))
		return
	}

	if err := json.NewEncoder(w).Encode(&automobile); err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "", err, true)
		return
	}
}

func (ah *AutomobileHandler) Delete(w http.ResponseWriter, r *http.Request) {
	caller := auth.GetUserFromRequest(r)
	usernamePathParam := mux.Vars(r)["name"]
	if caller != usernamePathParam {
		httputils.RespondWithError(w, http.StatusUnauthorized, "Deleting automobiles for other users is forbidden", nil, false)
		return
	}

	if dbErr := ah.DB.DeleteAutomobile(caller); dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, "Could not delete automobile")
		return
	}

	log.Printf("User %s's automobile was successfully deleted", caller)
	w.WriteHeader(http.StatusNoContent)
}
