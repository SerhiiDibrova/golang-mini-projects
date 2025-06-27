package helpers

import (
	"encoding/json"
	"errors"
	"net/http"

	"host_verifier"
	"server"
)

type filterValuesResponse struct {
	FilterValues interface{} `json:"filterValues"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func parseParams(params map[string]string) (string, string, string, map[string]string, error) {
	server, ok := params["server"]
	if !ok {
		return "", "", "", nil, errors.New("server parameter is required")
	}

	layer, ok := params["layer"]
	if !ok {
		return "", "", "", nil, errors.New("layer parameter is required")
	}

	serverType, ok := params["serverType"]
	if !ok {
		return "", "", "", nil, errors.New("serverType parameter is required")
	}

	filter, ok := params["filter"]
	if !ok {
		return "", "", "", nil, errors.New("filter parameter is required")
	}

	if filter == "" {
		return "", "", "", nil, errors.New("filter parameter cannot be empty")
	}

	filterParams := make(map[string]string)
	for key, value := range params {
		if key != "server" && key != "layer" && key != "serverType" && key != "filter" {
			filterParams[key] = value
		}
	}

	return server, layer, serverType, filterParams, nil
}

func getFilterValues(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	server, layer, serverType, filterParams, err := parseParams(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !host_verifier.AllowedHost(server) {
		http.Error(w, "Host not allowed", http.StatusBadGateway)
		return
	}

	serverClass, err := server.GetServerClass(serverType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	filterValues, err := serverClass.GetFilterValues(layer, filterParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	if filterValues == nil {
		http.Error(w, "Filter values not found", http.StatusBadGateway)
		return
	}

	response := filterValuesResponse{FilterValues: filterValues}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}