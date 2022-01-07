package handlers

import (
	"encoding/json"
	. "github.com/MonikaPalova/tarikatomobil.bg/db"
	"github.com/MonikaPalova/tarikatomobil.bg/httputils"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type PhotoHandler struct {
	DB *PhotoDBHandler
}

// GETs a photo by ID
func (ph *PhotoHandler) GetPhoto(w http.ResponseWriter, r *http.Request) {
	photoID, ok := mux.Vars(r)["id"]
	if !ok { // Should not happen
		httputils.RespondWithError(w, http.StatusInternalServerError, "Mux did not forward the request correctly", nil)
		return
	}

	photo, dbErr := ph.DB.GetPhotoByID(photoID)
	if dbErr != nil {
		httputils.RespondWithDBError(w, dbErr)
		return
	}

	bytes, err := json.Marshal(photo)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not marshal photo to JSON", err)
		return
	}
	_, _ = w.Write(bytes)
}

// Uploads a new photo
func (ph *PhotoHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	var photoToUpload model.Photo
	if err := json.NewDecoder(r.Body).Decode(&photoToUpload); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Could not parse request body as JSON", err)
		return
	}

	photoID, err := uuid.NewUUID()
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not generate a user ID", err)
		return
	}
	photoToUpload.ID = photoID.String()

	if dbErr := ph.DB.UploadPhoto(&photoToUpload); dbErr != nil {
		httputils.RespondWithDBError(w, dbErr)
		return
	}

	if err := json.NewEncoder(w).Encode(&photoToUpload); err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not marshal the created photo", err)
		return
	}

	log.Printf("Successfully created a photo with ID %s", photoID)
}

// Deletes a photo by ID.
func (ph *PhotoHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	photoID, ok := mux.Vars(r)["id"]
	if !ok { // Should not happen
		httputils.RespondWithError(w, http.StatusInternalServerError, "Mux did not forward the request correctly", nil)
	}
	if err := ph.DB.DeletePhoto(photoID); err != nil {
		httputils.RespondWithDBError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
