package helpers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type JobCompleteRequest struct {
	UUID        string `json:"uuid"`
	Successful  bool   `json:"successful"`
	ErrorMessage string `json:"errorMessage"`
}

type WpsService struct{}

func (s *WpsService) HandleNotification(uuid string, successful bool, errorMessage string) error {
	// handle the notification logic
	log.Println("Handling notification for job with UUID:", uuid)
	if !successful {
		log.Println("Job with UUID:", uuid, "failed with error message:", errorMessage)
	}
	return nil
}

func GetWpsService() *WpsService {
	return &WpsService{}
}

func JobComplete(w http.ResponseWriter, r *http.Request) {
	var jobCompleteRequest JobCompleteRequest
	err := json.NewDecoder(r.Body).Decode(&jobCompleteRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if jobCompleteRequest.UUID == "" {
		http.Error(w, "UUID is required", http.StatusBadRequest)
		return
	}

	if jobCompleteRequest.Successful != true && jobCompleteRequest.Successful != false {
		http.Error(w, "Successful must be a boolean value", http.StatusBadRequest)
		return
	}

	if jobCompleteRequest.ErrorMessage != "" && !jobCompleteRequest.Successful {
		// validation check for error message
	} else if jobCompleteRequest.ErrorMessage != "" && jobCompleteRequest.Successful {
		http.Error(w, "Error message is not allowed for successful jobs", http.StatusBadRequest)
		return
	}

	wpsService := GetWpsService()
	err = wpsService.HandleNotification(jobCompleteRequest.UUID, jobCompleteRequest.Successful, jobCompleteRequest.ErrorMessage)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type response struct {
		Message string `json:"message"`
	}
	resp := response{Message: "Job completed successfully"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}