package database

import (
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PortalBranding struct {
	gorm.Model
	Name        string
	Description string
	Logo        string
}

var db *gorm.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	if !regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).MatchString(dbHost) {
		log.Fatal("Invalid DB_HOST")
	}

	if port, err := strconv.Atoi(dbPort); err != nil || port < 1 || port > 65535 {
		log.Fatal("Invalid DB_PORT")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).MatchString(dbName) {
		log.Fatal("Invalid DB_NAME")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).MatchString(dbUser) {
		log.Fatal("Invalid DB_USER")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).MatchString(dbPassword) {
		log.Fatal("Invalid DB_PASSWORD")
	}

	dsn := dbHost + ":" + dbPort + "/" + dbName + "?" + "user=" + dbUser + "&password=" + dbPassword + "&sslmode=disable"
	var err1 error
	db, err1 = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err1 != nil {
		log.Fatalf("Error connecting to database: %v", err1)
	}
}

func GetPortalBranding() (*PortalBranding, error) {
	var portalBranding PortalBranding
	result := db.First(&portalBranding)
	if result.Error != nil {
		return nil, result.Error
	}
	return &portalBranding, nil
}