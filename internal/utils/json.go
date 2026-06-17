package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// RespondWithError sends a JSON response with an error message and the specified HTTP status code.
func RespondWithError(w http.ResponseWriter, code int, msg string) {
	resp := map[string]string{"error": msg}
	dat, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(dat))
}

// RespondWithJSON sends a JSON response with the specified payload and HTTP status code.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(dat))
}
