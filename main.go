package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type CssHandler struct{}

func (h *CssHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("CssHandler"))
}

type Filters struct {
	Server string `json:"server"`
	Layer  string `json:"layer"`
}

func getFilters(w http.ResponseWriter, r *http.Request) {
	filters := Filters{
		Server: "server1",
		Layer:  "layer1",
	}
	json.NewEncoder(w).Encode(filters)
}

func getLayerInfoJson(w http.ResponseWriter, r *http.Request) {
	layerInfo := map[string]string{
		"layer": "layer1",
	}
	json.NewEncoder(w).Encode(layerInfo)
}

func BeforeInterceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("BeforeInterceptor")
		next.ServeHTTP(w, r)
	})
}

func connectToDatabase() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbString := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"

	return sql.Open("postgres", dbString)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cssHandler := &CssHandler{}

	http.Handle("/css", cssHandler)
	http.HandleFunc("/getFilters", getFilters)
	http.HandleFunc("/getLayerInfoJson", getLayerInfoJson)

	db, err := connectToDatabase()
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	defer db.Close()

	log.Println("Database connection established")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), BeforeInterceptor(http.DefaultServeMux)))
}