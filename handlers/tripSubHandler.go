package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/MonikaPalova/tarikatomobil.bg/auth"
	"github.com/MonikaPalova/tarikatomobil.bg/db"
	"github.com/MonikaPalova/tarikatomobil.bg/httputils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type TripSubscriptionHandler struct {
	DB *db.TripSubscriptionDBHandler
}

func (tsh *TripSubscriptionHandler) Post(w http.ResponseWriter, r *http.Request) {
	caller := auth.GetUserFromRequest(r)
	if caller == "" {
		httputils.RespondWithError(w, http.StatusUnauthorized, "Only logged-in users can subscribe to a trip",
			nil, false)
		return
	}
	tripID := mux.Vars(r)["id"]

	freeSpots, dbErr := tsh.DB.GetFreeSpotsCount(tripID)
	if dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, "Could not check free spots for trip")
		return
	}
	if freeSpots <= 0 {
		httputils.RespondWithError(w, http.StatusBadRequest, "The trip is full", nil, false)
		return
	}

	if dbErr = tsh.DB.SubscribePassenger(caller, tripID); dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, fmt.Sprintf("Could not subscribe user %s for trip %s", caller, tripID))
		return
	}

	log.Printf("User %s successfully subscribed for trip %s", caller, tripID)
}

func (tsh *TripSubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	caller := auth.GetUserFromRequest(r)
	if caller == "" {
		httputils.RespondWithError(w, http.StatusUnauthorized, "Only logged-in users can unsubscribe from a trip",
			nil, false)
		return
	}
	tripID := mux.Vars(r)["id"]

	passengerToUnsubscribe := caller

	userToKick := r.URL.Query().Get("kickUser")
	if userToKick != "" {
		// Then check that the caller is actually the owner of the trip
		tripOwner, dbErr := tsh.DB.GetTripOwner(tripID)
		if dbErr != nil {
			httputils.RespondWithDBError(w, dbErr, "Kicking user failed")
			return
		}
		if tripOwner != caller {
			httputils.RespondWithError(w, http.StatusBadRequest, "Only trip owners can kick users", nil, false)
			return
		}
		passengerToUnsubscribe = userToKick
	}

	if dbErr := tsh.DB.UnsubscribePassenger(passengerToUnsubscribe, tripID); dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, "Unsubscription failed")
		return
	}
	log.Printf("User %s successfully unubscribed from trip %s", caller, tripID)
}

func (tsh *TripSubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	tripID := mux.Vars(r)["id"]
	passengers, dbErr := tsh.DB.GetPassengersForTrip(tripID)
	if dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, fmt.Sprintf("Could not get passengers for trip %s", tripID))
		return
	}

	if err := json.NewEncoder(w).Encode(passengers); err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not serialize passengers", err, true)
		return
	}

	log.Printf("Successfully retrieved passengers for trip %s", tripID)
}
