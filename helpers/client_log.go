package helpers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type LogRequest struct {
	Level  string `json:"level"`
	Message string `json:"message"`
}

type LogResponse struct {
	Status string `json:"status"`
}

func ClientLog(w http.ResponseWriter, r *http.Request) {
	var logRequest LogRequest
	err := json.NewDecoder(r.Body).Decode(&logRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	level := strings.ToLower(logRequest.Level)
	message := logRequest.Message

	log.Printf("%s: %s", level, message)

	logResponse := LogResponse{Status: "success"}
	jsonResponse, err := json.Marshal(logResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}