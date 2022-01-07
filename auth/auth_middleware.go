package auth

import (
	"context"
	"github.com/MonikaPalova/tarikatomobil.bg/db"
	"github.com/MonikaPalova/tarikatomobil.bg/httputils"
	"github.com/MonikaPalova/tarikatomobil.bg/model"
	"log"
	"net/http"
)

type sessionUserCtxKey string

const usernameCtxKey sessionUserCtxKey = "username"

type SessionAuthMiddleware struct {
	DB *db.SessionDBHandler
}

func (sam SessionAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(model.SessionCookieName)
		if err != nil {
			httputils.RespondWithError(w, http.StatusUnauthorized, "This action requires being logged in", err, false)
			return // Terminate the request
		}
		session, dbErr := sam.DB.GetSessionByID(sessionCookie.Value)
		if dbErr != nil {
			httputils.RespondWithError(w, http.StatusUnauthorized, "Session is invalid or expired", dbErr.Err, true)
			return // Terminate the request
		}

		ctx := context.WithValue(r.Context(), usernameCtxKey, session.Owner)
		log.Printf("User %s successfully authenticated", session.Owner)
		next.ServeHTTP(w, r.WithContext(ctx)) // Set the username in the request's context
	})
}

func GetUserFromRequest(r *http.Request) string {
	return r.Context().Value(usernameCtxKey).(string)
}
