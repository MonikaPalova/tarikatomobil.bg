package model

import (
	"github.com/MonikaPalova/tarikatomobil.bg/utils"
	"time"
)

const (
	SessionIDLength    = 64
	SessionDuration    = time.Hour
	SessionCookieName  = "TARIKATOMOBIL-SESSION-ID"
	UsernameCookieName = "TARIKATOMOBIL-USERNAME"
)

type Session struct {
	ID         string    // the value of the cookie (secure and identifies the session in the database)
	Owner      string    // the name of the user who owns the session
	Expiration time.Time // the session is not longer valid once Expiration is in the past
}

func NewSession(owner string) (Session, error) {
	s := Session{
		ID:         "",
		Owner:      owner,
		Expiration: time.Now().Add(SessionDuration),
	}
	sessionID, err := utils.GenerateRandomString(SessionIDLength)
	if err != nil {
		return s, err
	}
	s.ID = sessionID
	return s, nil
}
