package helpers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/your-project/logger"
)

type logLayerErrorRequest struct {
	Layer string `json:"layer"`
}

type jsonResponse struct {
	Message string `json:"message"`
}

func logLayerError(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	if r.Header.Get("Authorization") == "" {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	var request logLayerErrorRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Layer == "" {
		jsonResponse := jsonResponse{
			Message: "No layer specified",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(jsonResponse)
		return
	}

	err = logger.LogError(request.Layer)
	if err != nil {
		jsonResponse := jsonResponse{
			Message: "Error logging layer error",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(jsonResponse)
		return
	}

	jsonResponse := jsonResponse{
		Message: "Layer error logged successfully",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonResponse)
}

func main() {
	http.HandleFunc("/log-layer-error", logLayerError)
	log.Fatal(http.ListenAndServe(":8080", nil))
}