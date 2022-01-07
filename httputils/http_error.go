package httputils

import (
	"fmt"
	"github.com/MonikaPalova/tarikatomobil.bg/db"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, statusCode int, userResponse string, err error, hideErrorCause bool) {
	finalErrorMsg := userResponse
	if err != nil {
		finalErrorMsg = fmt.Sprintf("%s :%s", finalErrorMsg, err.Error())
	}

	w.WriteHeader(statusCode)
	log.Println(finalErrorMsg) // always log the full error
	if statusCode == http.StatusInternalServerError {
		// Then do not even respond with the userResponse
		_, _ = w.Write([]byte("Internal server error"))
	}
	if !hideErrorCause {
		_, _ = w.Write([]byte(finalErrorMsg)) // If the error cause is not to be hidden, respond with the full error
	} else {
		_, _ = w.Write([]byte(userResponse)) // If the error cause is to be hidden, only send the userResponse back
	}
}

func RespondWithDBError(w http.ResponseWriter, dbError *db.DBError, errorMsg string) {
	var statusCode int
	switch dbError.ErrorType {
	case db.ErrNotFound:
		statusCode = http.StatusNotFound
	case db.ErrInternal:
		statusCode = http.StatusInternalServerError
	case db.ErrConflict:
		statusCode = http.StatusConflict
	}

	userResponse := errorMsg
	if dbError.UserMessage != "" { // Add the user-friendly DB error if it is non-empty
		userResponse = fmt.Sprintf("%s: %s", userResponse, dbError.UserMessage)
	}

	RespondWithError(w, statusCode, userResponse, dbError.Err, true /* always hide the underlying DB error*/)
}
