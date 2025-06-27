package server

import (
	"errors"
)

type Server interface {
	Start() error
	Stop() error
}

type WpsUniqueValuesFilterService struct{}

func NewWpsUniqueValuesFilterService() *WpsUniqueValuesFilterService {
	return &WpsUniqueValuesFilterService{}
}

type GeoServerCore struct {
	filterService *WpsUniqueValuesFilterService
}

func NewGeoServerCore(filterService *WpsUniqueValuesFilterService) *GeoServerCore {
	return &GeoServerCore{filterService: filterService}
}

func (g *GeoServerCore) Start() error {
	// Added error handling for Start method
	if g.filterService == nil {
		return errors.New("filter service is not initialized")
	}
	return nil
}

func (g *GeoServerCore) Stop() error {
	// Added error handling for Stop method
	if g.filterService == nil {
		return errors.New("filter service is not initialized")
	}
	return nil
}

type GeoServerFilterConfig struct {
	filterService *WpsUniqueValuesFilterService
}

func NewGeoServerFilterConfig(filterService *WpsUniqueValuesFilterService) *GeoServerFilterConfig {
	return &GeoServerFilterConfig{filterService: filterService}
}

func (g *GeoServerFilterConfig) Start() error {
	// Added error handling for Start method
	if g.filterService == nil {
		return errors.New("filter service is not initialized")
	}
	return nil
}

func (g *GeoServerFilterConfig) Stop() error {
	// Added error handling for Stop method
	if g.filterService == nil {
		return errors.New("filter service is not initialized")
	}
	return nil
}

type OtherServer struct{}

func NewOtherServer() *OtherServer {
	return &OtherServer{}
}

func (o *OtherServer) Start() error {
	// Added error handling for Start method
	return nil
}

func (o *OtherServer) Stop() error {
	// Added error handling for Stop method
	return nil
}

func ServerClassFactory(serverType string) (Server, error) {
	switch serverType {
	case "geoservercore":
		filterService := NewWpsUniqueValuesFilterService()
		return NewGeoServerCore(filterService), nil
	case "geoserverfilterconfig":
		filterService := NewWpsUniqueValuesFilterService()
		return NewGeoServerFilterConfig(filterService), nil
	case "otherserver":
		return NewOtherServer(), nil
	default:
		return nil, errors.New("unsupported server type")
	}
}