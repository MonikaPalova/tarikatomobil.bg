package handlers

import (
	"encoding/json"
	"github.com/MonikaPalova/tarikatomobil.bg/auth"
	"github.com/MonikaPalova/tarikatomobil.bg/db"
	"github.com/MonikaPalova/tarikatomobil.bg/httputils"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type TripsHandler struct {
	DB *db.TripsDBHandler
}

func (th *TripsHandler) Post(w http.ResponseWriter, r *http.Request) {
	caller := auth.GetUserFromRequest(r)
	if caller == "" {
		httputils.RespondWithError(w, http.StatusUnauthorized, "Only logged-in users can create trips", nil, false)
		return
	}

	var tripToCreate model.Trip
	if err := json.NewDecoder(r.Body).Decode(&tripToCreate); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Could not parse request body as JSON", err, false)
		return
	}

	if err := tripToCreate.Validate(); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Cannot create trip", err, false)
		return
	}

	tripToCreate.DriverName = caller
	tripID, err := uuid.NewUUID()
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not generate a trip ID", err, true)
		return
	}
	tripToCreate.ID = tripID.String()

	if dbErr := th.DB.CreateTrip(tripToCreate); dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, "Could not create trip")
		return
	}

	if err := json.NewEncoder(w).Encode(tripToCreate); err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not serialize created trip", err, true)
		return
	}

	log.Printf("Trip from %s to %s created successfully by %s", tripToCreate.From, tripToCreate.To, caller)
}

func (th *TripsHandler) Delete(w http.ResponseWriter, r *http.Request) {

}

func (th *TripsHandler) GetAll(w http.ResponseWriter, r *http.Request) {

}

func (th *TripsHandler) GetSingle(w http.ResponseWriter, r *http.Request) {
	tripID := mux.Vars(r)["id"]
	trip, dbErr := th.DB.GetTrip(tripID)
	if dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, "Could not get trip")
		return
	}

	if err := json.NewEncoder(w).Encode(trip); err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not serialize trip", err, true)
		return
	}
}
