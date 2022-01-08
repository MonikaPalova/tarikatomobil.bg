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
)

type ReviewsHandler struct {
	DB *db.ReviewsDBHandler
}

func (rh *ReviewsHandler) Post(w http.ResponseWriter, r *http.Request) {
	caller := auth.GetUserFromRequest(r)
	usernamePathParam := mux.Vars(r)["name"]
	if caller != usernamePathParam {
		httputils.RespondWithError(w, http.StatusUnauthorized,
			"Reviews can only be created from the account you are logged in", nil, false)
		return
	}

	var reviewToCreate model.Review
	if err := json.NewDecoder(r.Body).Decode(&reviewToCreate); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Could not parse request body as JSON", err, false)
		return
	}
	reviewToCreate.FromUser = caller
	if err := reviewToCreate.Validate(); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "Could not create review", err, false)
		return
	}

	reviewID, err := uuid.NewUUID()
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not generate a review ID", err, true)
		return
	}
	reviewToCreate.ID = reviewID.String()

	if dbErr := rh.DB.CreateReview(reviewToCreate); dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, "Could not create review")
		return
	}

	log.Printf("Review from user %s to user %s with rating %d successfully created", reviewToCreate.FromUser,
		reviewToCreate.ForUser, reviewToCreate.Rating)
}

func (rh *ReviewsHandler) Get(w http.ResponseWriter, r *http.Request) {
	relation := r.URL.Query().Get("relation")
	username := mux.Vars(r)["name"]
	var reviews []model.Review
	var dbErr *db.DBError

	if relation == "for" {
		reviews, dbErr = rh.DB.GetReviewsForUser(username)
	} else {
		reviews, dbErr = rh.DB.GetReviewsFromUser(username)
	}

	if dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, "Could not get reviews")
		return
	}

	if err := json.NewEncoder(w).Encode(reviews); err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Could not serialize reviews into JSON",
			err, true)
		return
	}
}

func (rh *ReviewsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	caller := auth.GetUserFromRequest(r)
	usernamePathParam := mux.Vars(r)["name"]
	if caller != usernamePathParam {
		httputils.RespondWithError(w, http.StatusUnauthorized, "Deleting reviews for other users is forbidden", nil, false)
		return
	}

	reviewID := mux.Vars(r)["review_id"]
	if dbErr := rh.DB.DeleteReview(reviewID, caller); dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, fmt.Sprintf("Could not delete review %s", reviewID))
		return
	}

	log.Printf("Review %s was successfully deleted", reviewID)
	w.WriteHeader(http.StatusNoContent)
}
