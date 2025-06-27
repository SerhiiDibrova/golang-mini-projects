package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var cache = sync.Map{}
var logger = log.New(log.Writer(), "ifAllowed: ", log.Ldate|log.Ltime|log.Lshortfile)
var allowedURLs []string
var cacheExpiration time.Duration = 1 * time.Hour
var updateInterval time.Duration = 1 * time.Minute

func init() {
	var err error
	dsn := "host=localhost user=myuser password=mypassword dbname=mydb port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal(err)
	}
	go updateAllowedURLs()
}

func updateAllowedURLs() {
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := updateAllowedURLsFromDB()
			if err != nil {
				logger.Println(err)
			}
		}
	}
}

func updateAllowedURLsFromDB() error {
	var newAllowedURLs []string
	result := db.Find(&newAllowedURLs)
	if result.Error != nil {
		return result.Error
	}
	allowedURLs = newAllowedURLs
	return nil
}

func ifAllowed(url string) bool {
	if val, ok := cache.Load(url); ok {
		return val.(bool)
	}

	for _, allowedURL := range allowedURLs {
		if url == allowedURL {
			cache.Store(url, true)
			go expireCacheEntry(url, cacheExpiration)
			return true
		}
	}

	cache.Store(url, false)
	go expireCacheEntry(url, cacheExpiration)
	return false
}

func expireCacheEntry(key string, expiration time.Duration) {
	time.Sleep(expiration)
	cache.Delete(key)
}

func main() {
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if ifAllowed(url) {
			w.Write([]byte("Allowed"))
		} else {
			w.Write([]byte("Not Allowed"))
		}
	})
	http.ListenAndServe(":8080", nil)
}