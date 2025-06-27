package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var router *mux.Router
var portalBranding string
var once sync.Once

type PortalBranding struct {
	gorm.Model
	Name        string
	Description string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB_DSN")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	once.Do(func() {
		var pb PortalBranding
		db.First(&pb)
		portalBranding = pb.Name
	})

	router = mux.NewRouter()
	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/footer-content", footerContentHandler).Methods("GET")
	router.HandleFunc("/layer-config", layerConfigHandler).Methods("GET")

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Index page: " + portalBranding))
}

func footerContentHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Footer content: " + portalBranding))
}

func layerConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Layer config: " + portalBranding))
}