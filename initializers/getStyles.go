package initializers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Style struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
}

type Server struct {
	Type string `json:"type"`
	Host string `json:"host"`
}

type GetStylesRequest struct {
	Server     string `json:"server"`
	Layer      string `json:"layer"`
	ServerType string `json:"serverType"`
	Filter     string `json:"filter"`
}

type GetStylesResponse struct {
	Styles []Style `json:"styles"`
}

var allowedHosts = []string{"host1", "host2"}
var mutex = &sync.Mutex{}

func GetStyles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	server := vars["server"]
	layer := vars["layer"]
	serverType := vars["serverType"]
	filter := r.URL.Query().Get("filter")

	if !isHostAllowed(server) {
		http.Error(w, "Host not allowed", http.StatusBadGateway)
		return
	}

	if !validateInput(server, layer, serverType, filter) {
		http.Error(w, "Invalid input parameters", http.StatusBadRequest)
		return
	}

	srv, err := instantiateServer(serverType, server)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	styles, err := srv.GetStyles(layer, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := GetStylesResponse{Styles: styles}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func isHostAllowed(host string) bool {
	for _, allowedHost := range allowedHosts {
		if host == allowedHost {
			return true
		}
	}
	return false
}

func instantiateServer(serverType string, host string) (Server, error) {
	if serverType == "type1" {
		return Server{Type: serverType, Host: host}, nil
	} else if serverType == "type2" {
		return Server{Type: serverType, Host: host}, nil
	} else {
		return Server{}, errors.New("unsupported server type")
	}
}

func (s Server) GetStyles(layer string, filter string) ([]Style, error) {
	// implement logic to retrieve styles from the server
	// for demonstration purposes, return a hardcoded list of styles
	return []Style{
		{Name: "style1", URI: "uri1"},
		{Name: "style2", URI: "uri2"},
	}, nil
}

func validateInput(server string, layer string, serverType string, filter string) bool {
	if server == "" || layer == "" || serverType == "" {
		return false
	}
	if len(server) > 100 || len(layer) > 100 || len(serverType) > 100 {
		return false
	}
	if filter != "" && len(filter) > 100 {
		return false
	}
	return true
}

func authenticate(w http.ResponseWriter, r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return false
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), 12)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return false
	}
	return true
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/getStyles/{server}/{layer}/{serverType}", func(w http.ResponseWriter, r *http.Request) {
		if !authenticate(w, r) {
			return
		}
		GetStyles(w, r)
	}).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}