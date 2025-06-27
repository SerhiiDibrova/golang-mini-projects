package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

type RequestParams struct {
	PropertyName string `json:"propertyName"`
}

var cacheInstance *cache.Cache
var mutex sync.RWMutex

func init() {
	cacheInstance = cache.New(5*60*time.Minute, 10*time.Minute)
}

func urlListForFeatureRequestService(propertyName string) ([]string, error) {
	// This function should be implemented according to the requirements
	// For now, it returns a mock response
	return []string{"url1", "url2", "url3"}, nil
}

func getCachedResponse(propertyName string) ([]string, bool) {
	mutex.RLock()
	value, found := cacheInstance.Get(propertyName)
	mutex.RUnlock()
	if found {
		return value.([]string), true
	}
	return nil, false
}

func setCacheResponse(propertyName string, response []string) {
	mutex.Lock()
	cacheInstance.Set(propertyName, response, cache.NoExpiration)
	mutex.Unlock()
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var params RequestParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if params.PropertyName == "" {
		http.Error(w, "Property name is required", http.StatusBadRequest)
		return
	}

	cachedResponse, found := getCachedResponse(params.PropertyName)
	if found {
		json.NewEncoder(w).Encode(cachedResponse)
		return
	}

	urls, err := urlListForFeatureRequestService(params.PropertyName)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	setCacheResponse(params.PropertyName, urls)
	json.NewEncoder(w).Encode(urls)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/urlListForFeatureRequest", handleRequest).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}