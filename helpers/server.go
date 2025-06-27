package helpers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Server interface {
	getFilterValues(layer string, filter string) ([]string, error)
}

type NcwmsServer struct{}

func (s *NcwmsServer) getFilterValues(layer string, filter string) ([]string, error) {
	resp, err := http.Get("https://example.com/ncwms/" + layer + "/" + filter)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to retrieve filter values")
	}

	var filterValues []string
	// implementation to parse filter values from response
	return filterValues, nil
}

type AlaServer struct{}

func (s *AlaServer) getFilterValues(layer string, filter string) ([]string, error) {
	resp, err := http.Get("https://example.com/ala/" + layer + "/" + filter)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to retrieve filter values")
	}

	var filterValues []string
	// implementation to parse filter values from response
	return filterValues, nil
}

type CoreGeoserverServer struct{}

func (s *CoreGeoserverServer) getFilterValues(layer string, filter string) ([]string, error) {
	resp, err := http.Get("https://example.com/geoserver/" + layer + "/" + filter)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to retrieve filter values")
	}

	var filterValues []string
	// implementation to parse filter values from response
	return filterValues, nil
}

type WpsUniqueValuesFilterService struct{}

func (s *WpsUniqueValuesFilterService) getFilterValues(layer string, filter string) ([]string, error) {
	resp, err := http.Get("https://example.com/wps/" + layer + "/" + filter)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to retrieve filter values")
	}

	var filterValues []string
	// implementation to parse filter values from response
	return filterValues, nil
}

func _getServerClass(serverType string) (Server, error) {
	switch serverType {
	case "NcwmsServer":
		return &NcwmsServer{}, nil
	case "AlaServer":
		return &AlaServer{}, nil
	case "CoreGeoserverServer":
		return &CoreGeoserverServer{}, nil
	case "WpsUniqueValuesFilterService":
		return &WpsUniqueValuesFilterService{}, nil
	default:
		return nil, errors.New("unsupported server type")
	}
}

func GetFilterValues(serverType string, layer string, filter string) ([]string, error) {
	server, err := _getServerClass(serverType)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return server.getFilterValues(layer, filter)
}

func main() {
	serverType := "NcwmsServer"
	layer := "layer1"
	filter := "filter1"
	values, err := GetFilterValues(serverType, layer, filter)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(values)
	}
}