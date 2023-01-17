package utils

import (
	"Golangcode/config"
	"encoding/json"
	"net/http"
)

// RespondWithJSON ...
func RespondWithJSON(msg string, w http.ResponseWriter, r *http.Request) {

	body := config.Body{ResponseCode: 200, Message: msg}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)

}
