package handlers

import (
	. "github.com/MonikaPalova/tarikatomobil.bg/db"
)

type TripsHandler struct {
	DB Database
}

/*func (t TripsHandler) Get(w http.ResponseWriter, r *http.Request) {
	trips, err := t.DB.GetTrips()
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not fetch trips from DB", err)
		return
	}

	bytes, err := json.Marshal(trips)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not marshal trips to JSON", err)
		return
	}

	_, _ = w.Write(bytes)
}

func (t TripsHandler) Post(w http.ResponseWriter, r *http.Request) {
	var tripToCreate model.Trip
	if err := json.NewDecoder(r.Body).Decode(&tripToCreate); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Could not parse JSON request body", err)
		return
	}

	// TODO validate reg number is of real car

}
*/