package httputils

import (
	"fmt"
	"github.com/MonikaPalova/tarikatomobil.bg/db"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, statusCode int, errorMsg string, err error) {
	finalErrorMsg := errorMsg
	if err != nil {
		finalErrorMsg = fmt.Sprintf("%s :%s", finalErrorMsg, err.Error())
	}

	w.WriteHeader(statusCode)
	log.Println(finalErrorMsg)
	if statusCode != http.StatusInternalServerError {
		_, _ = w.Write([]byte(finalErrorMsg)) // Only write the error if it's not 500
	}
}

func RespondWithDBError(w http.ResponseWriter, dbError *db.DBError, additionalMsg ...string) {
	var statusCode int
	switch dbError.ErrorType {
	case db.ErrNotFound:
		statusCode = http.StatusNotFound
	case db.ErrInternal:
		statusCode = http.StatusInternalServerError
	case db.ErrConflict:
		statusCode = http.StatusConflict
	}

	errorMsg := ""
	if len(additionalMsg) > 0 {
		errorMsg = additionalMsg[0]
	}

	RespondWithError(w, statusCode, errorMsg, dbError.Err)
}
