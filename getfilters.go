package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Server interface {
	getFilters(layer string, filter string) ([]string, error)
}

type ServerA struct{}

func (s *ServerA) getFilters(layer string, filter string) ([]string, error) {
	return []string{"filter1", "filter2"}, nil
}

type ServerB struct{}

func (s *ServerB) getFilters(layer string, filter string) ([]string, error) {
	return []string{"filter3", "filter4"}, nil
}

func _getServerClass(serverType string) (Server, error) {
	switch serverType {
	case "A":
		return &ServerA{}, nil
	case "B":
		return &ServerB{}, nil
	default:
		return nil, errors.New("unknown server type")
	}
}

func parseParams(r *http.Request) (string, string, string, error) {
	layer := r.URL.Query().Get("layer")
	serverType := r.URL.Query().Get("serverType")
	filter := r.URL.Query().Get("filter")
	if layer == "" || serverType == "" {
		return "", "", "", errors.New("layer and serverType parameters are required")
	}
	return layer, serverType, filter, nil
}

func allowedHost(hostVerifier string, host string) bool {
	allowedHosts := []string{"host1", "host2"}
	for _, allowedHost := range allowedHosts {
		if allowedHost == host {
			return true
		}
	}
	return false
}

func getFilters(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	layer, serverType, filter, err := parseParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	host := r.Header.Get("Host")

	if !allowedHost("hostVerifier", host) {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	server, err := _getServerClass(serverType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filters, err := server.getFilters(layer, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonFilters, err := json.Marshal(filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonFilters)
}

func main() {
	http.HandleFunc("/getFilters", getFilters)
	fmt.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}