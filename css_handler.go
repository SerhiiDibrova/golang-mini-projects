package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type PortalBranding struct {
	CssContent string `json:"css_content"`
}

type CssHandler struct{}

func (h *CssHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./portal_branding.db")
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Printf("Failed to ping database: %v", err)
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}

	var portalBranding PortalBranding
	err = db.QueryRow("SELECT css_content FROM portal_branding").Scan(&portalBranding.CssContent)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found in portal branding table")
			http.Error(w, "Failed to retrieve portal branding data", http.StatusNotFound)
		} else {
			log.Printf("Failed to retrieve portal branding data: %v", err)
			http.Error(w, "Failed to retrieve portal branding data", http.StatusInternalServerError)
		}
		return
	}

	cssContent, err := h.GetCssContent(portalBranding)
	if err != nil {
		log.Printf("Failed to generate CSS content: %v", err)
		http.Error(w, "Failed to generate CSS content", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/css")
	w.Write([]byte(cssContent))
}

func (h *CssHandler) GetCssContent(portalBranding PortalBranding) (string, error) {
	return portalBranding.CssContent, nil
}

func main() {
	http.HandleFunc("/css", func(w http.ResponseWriter, r *http.Request) {
		cssHandler := &CssHandler{}
		cssHandler.HandleGet(w, r)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}