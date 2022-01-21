package handlers

import (
	"encoding/json"
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
	"strings"
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
	trips, dbErr := th.DB.GetTrips()
	if dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, "Could not retrieve trips")
		return
	}

	if err := filterTrips(&trips, r.URL.Query()); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Query string is not properly formatted", err, false)
		return
	}

	if err := json.NewEncoder(w).Encode(trips); err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not marshal retrieved trips",
			err, true)
		return
	}
	log.Println("Trips successfully retrieved")
}

func filterTrips(trips *[]model.Trip, query url.Values) error {
	if fromParam, ok := query["from"]; ok {
		removeUnless(trips, func(trip model.Trip) bool {
			return strings.EqualFold(trip.From, fromParam[0])
		})
	}

	if toParam, ok := query["to"]; ok {
		removeUnless(trips, func(trip model.Trip) bool {
			return strings.EqualFold(trip.To, toParam[0])
		})
	}

	var err error
	const timeFormat = "2006-Jan-02"
	if beforeParam, ok := query["before"]; ok {
		var before time.Time
		if before, err = time.Parse(timeFormat, beforeParam[0]); err != nil {
			return fmt.Errorf("could not parse the before query param: %s", err.Error())
		}
		removeUnless(trips, func(trip model.Trip) bool {
			return trip.When.Before(before)
		})
	}

	if afterParam, ok := query["after"]; ok {
		var after time.Time
		if after, err = time.Parse(timeFormat, afterParam[0]); err != nil {
			return fmt.Errorf("could not parse the after query param: %s", err.Error())
		}
		removeUnless(trips, func(trip model.Trip) bool {
			return trip.When.After(after)
		})
	}
	if maxPriceParam, ok := query["maxPrice"]; ok {
		var maxPrice float64
		if maxPrice, err = strconv.ParseFloat(maxPriceParam[0], 32); err != nil {
			return fmt.Errorf("could not parse the maxPrice query param as a float: %s", err.Error())
		}
		removeUnless(trips, func(trip model.Trip) bool {
			return trip.Price <= maxPrice
		})
	}

	if airConditioningParam, ok := query["airConditioning"]; ok {
		var airConditioning bool
		if airConditioning, err = strconv.ParseBool(airConditioningParam[0]); err != nil {
			return fmt.Errorf("could not parse the airConditioning query param as a boolean: %s", err.Error())
		}
		removeUnless(trips, func(trip model.Trip) bool {
			return trip.AirConditioning == airConditioning
		})
	}
	if smokingParam, ok := query["smoking"]; ok {
		var smoking bool
		if smoking, err = strconv.ParseBool(smokingParam[0]); err != nil {
			return fmt.Errorf("could not parse the smoking query param as a boolean: %s", err.Error())
		}
		removeUnless(trips, func(trip model.Trip) bool {
			return trip.Smoking == smoking
		})
	}
	if petsParam, ok := query["pets"]; ok {
		var pets bool
		if pets, err = strconv.ParseBool(petsParam[0]); err != nil {
			return fmt.Errorf("could not parse the pets query param as a boolean: %s", err.Error())
		}
		removeUnless(trips, func(trip model.Trip) bool {
			return trip.Pets == pets
		})
	}
	return nil
}

func removeUnless(trips *[]model.Trip, condition func(model.Trip) bool) {
	n := 0
	for _, trip := range *trips {
		if condition(trip) {
			(*trips)[n] = trip
			n++
		}
	}
	*trips = (*trips)[:n]
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

	isDriver := r.URL.Query().Get("isDriver") == "true"

	trips, dbErr := th.DB.GetTripsForUser(caller, isDriver)
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
