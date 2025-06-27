package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type SystemController struct{}

func (sc *SystemController) error(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

type urlListForLayerRequest struct {
	FieldName string `json:"fieldName"`
}

type urlListForLayerResponse struct {
	URLs []string `json:"urls"`
}

func _loadCommonFields(r *http.Request) (map[string]string, error) {
	commonFields := make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&commonFields)
	if err != nil {
		return nil, err
	}
	return commonFields, nil
}

func _performProxyingIfAllowed(r *http.Request) error {
	proxyAllowed := r.Header.Get("Proxy-Allow")
	if proxyAllowed == "true" {
		// implement proxying logic here
	}
	return nil
}

func urlListStreamProcessor(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")
	encoder := json.NewEncoder(w)
	for i := 0; i < 10; i++ {
		url := fmt.Sprintf("https://example.com/url-%d", i)
		err := encoder.Encode(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		flusher.Flush()
	}
}

func _executeExternalRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func urlListForLayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fieldName := vars["fieldName"]
	if fieldName == "" {
		sc := &SystemController{}
		sc.error(w, errors.New("field name is required"))
		return
	}
	var requestBody urlListForLayerRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		sc := &SystemController{}
		sc.error(w, err)
		return
	}
	if requestBody.FieldName == "" {
		sc := &SystemController{}
		sc.error(w, errors.New("field name is required in request body"))
		return
	}
	commonFields, err := _loadCommonFields(r)
	if err != nil {
		sc := &SystemController{}
		sc.error(w, err)
		return
	}
	err = _performProxyingIfAllowed(r)
	if err != nil {
		sc := &SystemController{}
		sc.error(w, err)
		return
	}
	urlListStreamProcessor(w, r)
	url := "https://example.com/external-request"
	body, err := _executeExternalRequest(url)
	if err != nil {
		sc := &SystemController{}
		sc.error(w, err)
		return
	}
	var response urlListForLayerResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		sc := &SystemController{}
		sc.error(w, err)
		return
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/url-list-for-layer/{fieldName}", urlListForLayerHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}