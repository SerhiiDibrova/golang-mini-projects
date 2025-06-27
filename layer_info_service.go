package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type LayerInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LayerInfoService struct {
	db *gorm.DB
}

func NewLayerInfoService(db *gorm.DB) *LayerInfoService {
	return &LayerInfoService{db: db}
}

func (s *LayerInfoService) GetLayerInfo(id int) (*LayerInfo, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	var layerInfo LayerInfo
	if err := s.db.Where("id = ?", id).First(&layerInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("layer info not found")
		}
		return nil, err
	}
	return &layerInfo, nil
}

func (s *LayerInfoService) GetLayerInfoJSON(id int) ([]byte, error) {
	layerInfo, err := s.GetLayerInfo(id)
	if err != nil {
		return nil, err
	}
	if layerInfo == nil {
		return nil, errors.New("layer info is nil")
	}
	layerInfoJSON, err := json.Marshal(layerInfo)
	if err != nil {
		return nil, err
	}
	return layerInfoJSON, nil
}

func main() {
	db, err := gorm.Open("postgres", "host=localhost user=myuser password=mypassword dbname=mydb port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if db == nil {
		log.Fatal("database connection is nil")
	}
	defer db.Close()

	layerInfoService := NewLayerInfoService(db)
	layerInfoJSON, err := layerInfoService.GetLayerInfoJSON(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(layerInfoJSON))
}