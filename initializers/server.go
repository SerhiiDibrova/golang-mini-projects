package initializers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Server interface {
	Start() error
	HandleRequest(w http.ResponseWriter, r *http.Request)
	getStyles(layer string) ([]string, error)
}

type NcwmsServer struct {
	styles map[string][]string
	mu     sync.RWMutex
}

func (s *NcwmsServer) Start() error {
	return nil
}

func (s *NcwmsServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/getFeatureCount" {
		w.Write([]byte("Feature count endpoint"))
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (s *NcwmsServer) getStyles(layer string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if styles, ok := s.styles[layer]; ok {
		return styles, nil
	}
	return nil, errors.New("layer not found")
}

type AlaServer struct {
	styles map[string][]string
	mu     sync.RWMutex
}

func (s *AlaServer) Start() error {
	return nil
}

func (s *AlaServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/getFeatureCount" {
		w.Write([]byte("Feature count endpoint"))
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (s *AlaServer) getStyles(layer string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if styles, ok := s.styles[layer]; ok {
		return styles, nil
	}
	return nil, errors.New("layer not found")
}

type CoreGeoserverServer struct {
	styles map[string][]string
	mu     sync.RWMutex
}

func (s *CoreGeoserverServer) Start() error {
	return nil
}

func (s *CoreGeoserverServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/getFeatureCount" {
		w.Write([]byte("Feature count endpoint"))
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (s *CoreGeoserverServer) getStyles(layer string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if styles, ok := s.styles[layer]; ok {
		return styles, nil
	}
	return nil, errors.New("layer not found")
}

type ImosGeoserverServer struct {
	styles map[string][]string
	mu     sync.RWMutex
}

func (s *ImosGeoserverServer) Start() error {
	return nil
}

func (s *ImosGeoserverServer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/getFeatureCount" {
		w.Write([]byte("Feature count endpoint"))
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (s *ImosGeoserverServer) getStyles(layer string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if styles, ok := s.styles[layer]; ok {
		return styles, nil
	}
	return nil, errors.New("layer not found")
}

type ServerImpl struct {
	servers map[string]Server
	mu      sync.RWMutex
}

func NewServerImpl() *ServerImpl {
	return &ServerImpl{
		servers: make(map[string]Server),
	}
}

func (s *ServerImpl) AddServer(name string, server Server) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.servers[name] = server
}

func (s *ServerImpl) Start() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, server := range s.servers {
		if err := server.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (s *ServerImpl) HandleRequest(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, server := range s.servers {
		server.HandleRequest(w, r)
	}
}

func (s *ServerImpl) getStyles(layer string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, server := range s.servers {
		styles, err := server.getStyles(layer)
		if err == nil {
			return styles, nil
		}
	}
	return nil, errors.New("layer not found")
}

func main() {
	ncwmsServer := &NcwmsServer{
		styles: map[string][]string{
			"layer1": {"style1", "style2"},
		},
	}
	alaServer := &AlaServer{
		styles: map[string][]string{
			"layer2": {"style3", "style4"},
		},
	}
	coreGeoserverServer := &CoreGeoserverServer{
		styles: map[string][]string{
			"layer3": {"style5", "style6"},
		},
	}
	imosGeoserverServer := &ImosGeoserverServer{
		styles: map[string][]string{
			"layer4": {"style7", "style8"},
		},
	}

	serverImpl := NewServerImpl()
	serverImpl.AddServer("ncwms", ncwmsServer)
	serverImpl.AddServer("ala", alaServer)
	serverImpl.AddServer("coreGeoserver", coreGeoserverServer)
	serverImpl.AddServer("imosGeoserver", imosGeoserverServer)

	if err := serverImpl.Start(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serverImpl.HandleRequest(w, r)
	})

	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}