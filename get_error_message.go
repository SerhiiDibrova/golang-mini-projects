package main

import (
	"encoding/json"
	"errors"
	"log"
)

type ExecutionStatusResponse struct {
	ErrorMessage string `json:"error_message"`
}

func _getErrorMessage(executionStatusResponse []byte) (string, error) {
	var response ExecutionStatusResponse
	err := json.Unmarshal(executionStatusResponse, &response)
	if err != nil {
		log.Printf("Error unmarshaling execution status response: %v", err)
		return "", errors.New("invalid execution status response")
	}
	if response.ErrorMessage == "" {
		log.Println("Error message not found in execution status response")
		return "", errors.New("error message not found")
	}
	return response.ErrorMessage, nil
}

func main() {
	executionStatusResponse := []byte(`{"error_message": "Test error message"}`)
	errorMessage, err := _getErrorMessage(executionStatusResponse)
	if err != nil {
		log.Printf("Error getting error message: %v", err)
	} else {
		log.Println(errorMessage)
	}
}