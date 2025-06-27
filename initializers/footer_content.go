package initializers

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type FooterContent struct {
	Content string `json:"content"`
}

func getFooterContent(db *gorm.DB) (string, error) {
	var portalBranding struct {
		FooterContent string `json:"footer_content"`
	}
	result := db.First(&portalBranding)
	if result.Error != nil {
		log.Println("Error retrieving portal branding: ", result.Error)
		return "", result.Error
	}
	return portalBranding.FooterContent, nil
}

func footerContentHandler(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value("db").(*gorm.DB)
	if !ok {
		http.Error(w, "Database connection not found", http.StatusInternalServerError)
		log.Println("Database connection not found")
		return
	}
	footerContent, err := getFooterContent(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error retrieving footer content: ", err)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(footerContent))
}

func InitializeFooterContentEndpoint(router *http.ServeMux, db *gorm.DB) {
	httpHandler := http.HandlerFunc(footerContentHandler)
	router.Handle("/footer-content", httpHandler)
}