package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// TODO: probably should be moved to messaging package

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
