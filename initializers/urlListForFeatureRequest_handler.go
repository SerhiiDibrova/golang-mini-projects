package initializers

import (
	"encoding/json"
	"log"
	"net/http"

	"urlListForFeatureRequest_service"
)

func urlListForFeatureRequestHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("panic occurred: %v", err)
		}
	}()

	response, err := urlListForFeatureRequestService.UrlListForFeatureRequestService(r)
	if err != nil {
		if err, ok := err.(urlListForFeatureRequest_service.ServiceError); ok {
			http.Error(w, err.Message, err.StatusCode)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		log.Printf("error occurred: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("error encoding response: %v", err)
		return
	}
}