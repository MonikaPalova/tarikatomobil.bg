package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MonikaPalova/tarikatomobil.bg/auth"
	"github.com/MonikaPalova/tarikatomobil.bg/db"
	"github.com/MonikaPalova/tarikatomobil.bg/httputils"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type TripsHandler struct {
	DB            *db.TripsDBHandler
	AutomobilesDB *db.AutomobileDBHandler
}

func (th *TripsHandler) Post(w http.ResponseWriter, r *http.Request) {
	caller := auth.GetUserFromRequest(r)
	if caller == "" {
		httputils.RespondWithError(w, http.StatusUnauthorized, "Only logged-in users can create trips", nil, false)
		return
	}

	if _, dbErr := th.AutomobilesDB.GetUserAutomobile(caller); dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, "Users without automobiles cannot create trips")
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
	caller := auth.GetUserFromRequest(r)
	if caller == "" {
		httputils.RespondWithError(w, http.StatusUnauthorized, "Only logged-in users can delete trips", nil, false)
		return
	}
	tripID := mux.Vars(r)["id"]
	if dbErr := th.DB.DeleteTrip(tripID, caller); dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, "Could not delete trip")
		return
	}
	w.WriteHeader(http.StatusNoContent)
	log.Printf("Trip %s successfuly deleted by user %s", tripID, caller)
}

func (th *TripsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	filter, err := parseTripFilterQuery(r.URL.Query())
	if err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Query string is not properly formatted", err, false)
		return
	}
	trips, dbErr := th.DB.GetTrips(filter)
	if dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, "Could not retrieve trips")
		return
	}

	if err := json.NewEncoder(w).Encode(trips); err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not marshal retrieved trips",
			err, true)
		return
	}
	log.Println("Trips successfully retrieved")
}

func parseTripFilterQuery(query url.Values) (model.TripFilter, error) {
	filter := model.DefaultTripFilter()
	filter.From = query.Get("from")
	if filter.From == "" {
		return filter, errors.New("when filtering, the from query param is mandatory")
	}

	filter.To = query.Get("to")
	if filter.To == "" {
		return filter, errors.New("when filtering, the to query param is mandatory")
	}

	var err error
	if before, ok := query["before"]; ok {
		if filter.Before, err = time.Parse("2022-01-08T11:56:00Z", before[0]); err != nil {
			return filter, fmt.Errorf("could not parse the before query param: %s", err.Error())
		}
	}

	if after, ok := query["after"]; ok {
		if filter.After, err = time.Parse("2022-01-08T11:56:00Z", after[0]); err != nil {
			return filter, fmt.Errorf("could not parse the after query param: %s", err.Error())
		}
	}
	if maxPrice, ok := query["maxPrice"]; ok {
		if filter.MaxPrice, err = strconv.ParseFloat(maxPrice[0], 32); err != nil {
			return filter, fmt.Errorf("could not parse the maxPrice query param as a float: %s", err.Error())
		}
	}
	if query.Get("airConditioning") == "true" {
		filter.AirConditioning = true
	}
	if query.Get("smoking") == "true" {
		filter.Smoking = true
	}
	if query.Get("pets") == "true" {
		filter.Pets = true
	}
	return filter, nil
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

func (th *TripsHandler) GetTripsForUser(w http.ResponseWriter, r *http.Request) {
	caller := auth.GetUserFromRequest(r)
	if caller == "" {
		httputils.RespondWithError(w, http.StatusUnauthorized, "Only logged-in users can get their trips", nil, false)
		return
	}

	trips, dbErr := th.DB.GetTripsForUser(caller)
	if dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, fmt.Sprintf("Could not get trips for user %s", caller))
		return
	}

	if err := json.NewEncoder(w).Encode(trips); err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not serialize trips", err, true)
		return
	}

	log.Printf("Successfully retrieved trips for user %s", caller)
}
