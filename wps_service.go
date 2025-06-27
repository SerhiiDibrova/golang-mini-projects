package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type WpsService struct {
	db  *gorm.DB
	log *log.Logger
	cache map[string]string
	mu    sync.RWMutex
}

func NewWpsService(dsn string) (*WpsService, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	err = db.Exec("SELECT 1").Error
	if err != nil {
		return nil, err
	}

	return &WpsService{
		db:  db,
		log: log.Default(),
		cache: make(map[string]string),
	}, nil
}

func (s *WpsService) GetUrl() (string, error) {
	s.mu.RLock()
	url, ok := s.cache["url"]
	s.mu.RUnlock()
	if ok {
		return url, nil
	}

	var wpsExecutionStatus struct {
		Url string `gorm:"column:url"`
	}

	err := s.db.Model(&wpsExecutionStatus).First(&wpsExecutionStatus).Error
	if err != nil {
		s.log.Printf("error retrieving url from database: %v", err)
		return "", err
	}

	s.mu.Lock()
	s.cache["url"] = wpsExecutionStatus.Url
	s.mu.Unlock()

	return wpsExecutionStatus.Url, nil
}

func (s *WpsService) StartServer() {
	http.HandleFunc("/wps", func(w http.ResponseWriter, r *http.Request) {
		url, err := s.GetUrl()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			s.log.Printf("error handling request: %v", err)
			return
		}

		w.Write([]byte(url))
	})

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			s.mu.Lock()
			s.cache = make(map[string]string)
			s.mu.Unlock()
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	wpsService, err := NewWpsService("host=localhost user=myuser password=mypassword dbname=mydb port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	wpsService.StartServer()
}