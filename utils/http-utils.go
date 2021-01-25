package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// RespondWithJSON returns result in JSON format
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		RespondWithError(w, 500, "Unknown error")
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		err = errors.Wrap(err, "Response cannot be converted to JSON format")
		log.Panic(err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
}

// RespondWithError returns error in JSON format
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}