package httputils

import (
	"fmt"
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
