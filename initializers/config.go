package initializers

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ConfigHandler struct {
	db *gorm.DB
}

func (c *ConfigHandler) StartServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusBadRequest)
			return
		}
		w.Write([]byte("REST API server started"))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadEnvironmentVariables() (map[string]string, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	envVars := map[string]string{
		"DB_USERNAME": os.Getenv("DB_USERNAME"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
	}
	for _, value := range envVars {
		if value == "" {
			return nil, errors.New("environment variable is not set")
		}
	}
	return envVars, nil
}

func initializeDatabase(envVars map[string]string) (*gorm.DB, error) {
	dsn := envVars["DB_USERNAME"] + ":" + envVars["DB_PASSWORD"] + "@tcp(" + envVars["DB_HOST"] + ":" + envVars["DB_PORT"] + ")/" + envVars["DB_NAME"] + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func hashPassword(password string) (string, error) {
	hash := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

func verifyPassword(storedPassword string, providedPassword string) (bool, error) {
	hash, err := base64.StdEncoding.DecodeString(storedPassword)
	if err != nil {
		return false, err
	}
	providedHash := sha256.Sum256([]byte(providedPassword))
	return subtle.ConstantTimeCompare(hash, providedHash[:]) == 1, nil
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func main() {
	envVars, err := loadEnvironmentVariables()
	if err != nil {
		log.Fatal(err)
	}
	db, err := initializeDatabase(envVars)
	if err != nil {
		log.Fatal(err)
	}
	configHandler := &ConfigHandler{db: db}
	configHandler.StartServer()
}