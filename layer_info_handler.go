package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
)

type Server struct {
	Host string
}

func (s *Server) GetLayerInfo(layer string) (map[string]interface{}, error) {
	// implement layer info retrieval logic here
	layerInfo := make(map[string]interface{})
	layerInfo["name"] = layer
	layerInfo["description"] = "This is a sample layer"
	return layerInfo, nil
}

func verifyHost(host string) bool {
	_, err := url.Parse(host)
	return err == nil
}

func getLayerInfoHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
		}
	}()

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	serverParam := r.URL.Query().Get("server")
	layerParam := r.URL.Query().Get("layer")

	if serverParam == "" || layerParam == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	if !verifyHost(serverParam) {
		http.Error(w, "Invalid host", http.StatusBadRequest)
		return
	}

	s := &Server{Host: serverParam}

	layerInfo, err := s.GetLayerInfo(layerParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(layerInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func main() {
	http.HandleFunc("/getLayerInfoJson", getLayerInfoHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}