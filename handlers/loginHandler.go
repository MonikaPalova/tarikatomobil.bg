package handlers

import (
	"github.com/MonikaPalova/tarikatomobil.bg/db"
	"github.com/MonikaPalova/tarikatomobil.bg/httputils"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
	"log"
	"net/http"
)

type LoginHandler struct {
	UserDB    *db.UsersDBHandler
	SessionDB *db.SessionDBHandler
}

func (h LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	username, pass, ok := r.BasicAuth()
	if !ok {
		httputils.RespondWithError(w, http.StatusBadRequest, "Basic Authentication header was not provided", nil, false)
		return
	}

	if dbErr := h.UserDB.ValidateUserCredentials(username, pass); dbErr != nil {
		if dbErr.ErrorType == db.ErrNotFound {
			httputils.RespondWithError(w, http.StatusUnauthorized, "Invalid login credentials", nil, false)
		} else {
			httputils.RespondWithDBError(w, dbErr, "")
		}
		return
	}

	session, err := model.NewSession(username)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, "Login failed", err, true)
		return
	}

	if dbErr := h.SessionDB.InsertSession(session); dbErr != nil {
		httputils.RespondWithDBError(w, dbErr, "Login failed")
		return
	}

	sessionCookie := http.Cookie{
		Name:    model.SessionCookieName,
		Value:   session.ID,
		Expires: session.Expiration,
	}

	http.SetCookie(w, &sessionCookie)
	log.Printf("User %s successfully logged in", username)
	_, _ = w.Write([]byte("Login successful!"))
}
