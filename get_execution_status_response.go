package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ExecutionStatusResponse struct {
	Status string `json:"status"`
}

func _getExecutionStatusResponse(url string) (*ExecutionStatusResponse, error) {
	maxRetries := 3
	retryDelay := 500 * time.Millisecond
	var response *http.Response
	var err error
	var executionStatusResponse *ExecutionStatusResponse

	for attempt := 0; attempt <= maxRetries; attempt++ {
		response, err = http.Get(url)
		if err != nil {
			log.Printf("Error sending request to %s: %v", url, err)
			if attempt < maxRetries {
				time.Sleep(retryDelay)
				continue
			}
			return nil, err
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			log.Printf("Invalid response status code from %s: %d", url, response.StatusCode)
			if attempt < maxRetries {
				time.Sleep(retryDelay)
				continue
			}
			return nil, fmt.Errorf("invalid response status code: %d", response.StatusCode)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("Error reading response body from %s: %v", url, err)
			if attempt < maxRetries {
				time.Sleep(retryDelay)
				continue
			}
			return nil, err
		}

		err = json.Unmarshal(body, &executionStatusResponse)
		if err != nil {
			log.Printf("Error unmarshaling response from %s: %v", url, err)
			if attempt < maxRetries {
				time.Sleep(retryDelay)
				continue
			}
			return nil, err
		}

		return executionStatusResponse, nil
	}

	return nil, fmt.Errorf("max retries exceeded")
}

func main() {
	url := "http://example.com/execution-status"
	executionStatusResponse, err := _getExecutionStatusResponse(url)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Execution status response: %+v", executionStatusResponse)
}