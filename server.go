package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Filter struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Server interface {
	getFilters(layer string) ([]Filter, error)
}

type NcwmsServer struct{}

func (s *NcwmsServer) getFilters(layer string) ([]Filter, error) {
	if layer == "" {
		return nil, errors.New("layer is required")
	}
	filters := []Filter{
		{Name: "filter1", Type: "type1", Value: "value1"},
		{Name: "filter2", Type: "type2", Value: "value2"},
	}
	return filters, nil
}

type AlaServer struct{}

func (s *AlaServer) getFilters(layer string) ([]Filter, error) {
	if layer == "" {
		return nil, errors.New("layer is required")
	}
	filters := []Filter{
		{Name: "filter3", Type: "type3", Value: "value3"},
		{Name: "filter4", Type: "type4", Value: "value4"},
	}
	return filters, nil
}

type GeoserverCoreServer struct{}

func (s *GeoserverCoreServer) getFilters(layer string) ([]Filter, error) {
	if layer == "" {
		return nil, errors.New("layer is required")
	}
	filters := []Filter{
		{Name: "filter5", Type: "type5", Value: "value5"},
		{Name: "filter6", Type: "type6", Value: "value6"},
	}
	return filters, nil
}

type GeoserverFilterConfigServer struct{}

func (s *GeoserverFilterConfigServer) getFilters(layer string) ([]Filter, error) {
	if layer == "" {
		return nil, errors.New("layer is required")
	}
	filters := []Filter{
		{Name: "filter7", Type: "type7", Value: "value7"},
		{Name: "filter8", Type: "type8", Value: "value8"},
	}
	return filters, nil
}

type ImosServer struct{}

func (s *ImosServer) getFilters(layer string) ([]Filter, error) {
	if layer == "" {
		return nil, errors.New("layer is required")
	}
	filters := []Filter{
		{Name: "filter9", Type: "type9", Value: "value9"},
		{Name: "filter10", Type: "type10", Value: "value10"},
	}
	return filters, nil
}

func GetServer(serverType string) (Server, error) {
	if serverType == "" {
		return nil, errors.New("server type is required")
	}
	switch strings.ToLower(serverType) {
	case "ncwms":
		return &NcwmsServer{}, nil
	case "ala":
		return &AlaServer{}, nil
	case "geoservercore":
		return &GeoserverCoreServer{}, nil
	case "geoserverfilterconfig":
		return &GeoserverFilterConfigServer{}, nil
	default:
		return &ImosServer{}, nil
	}
}

func GetFilters(w http.ResponseWriter, r *http.Request) {
	serverType := r.URL.Query().Get("serverType")
	layer := r.URL.Query().Get("layer")

	if serverType == "" {
		http.Error(w, "server type is required", http.StatusBadRequest)
		return
	}

	if layer == "" {
		http.Error(w, "layer is required", http.StatusBadRequest)
		return
	}

	server, err := GetServer(serverType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filters, err := server.getFilters(layer)
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
	http.HandleFunc("/filters", GetFilters)
	log.Fatal(http.ListenAndServe(":8080", nil))
}