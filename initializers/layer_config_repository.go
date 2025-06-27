package initializers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type LayerConfig struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Settings    string `json:"settings"`
}

type LayerConfigRepository interface {
	GetLayerConfig(id string) (*LayerConfig, error)
	StoreLayerConfig(config *LayerConfig) error
}

type layerConfigRepository struct {
	db *sql.DB
}

func NewLayerConfigRepository(db *sql.DB) LayerConfigRepository {
	return &layerConfigRepository{db: db}
}

func (r *layerConfigRepository) GetLayerConfig(id string) (*LayerConfig, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if !strings.HasPrefix(id, "lc_") {
		id = "lc_" + id
	}
	var config LayerConfig
	row := r.db.QueryRow("SELECT id, name, description, settings FROM layer_configs WHERE id = $1", id)
	err := row.Scan(&config.ID, &config.Name, &config.Description, &config.Settings)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("layer config not found")
		}
		return nil, err
	}
	return &config, nil
}

func (r *layerConfigRepository) StoreLayerConfig(config *LayerConfig) error {
	if config == nil {
		return errors.New("config is required")
	}
	if config.ID == "" {
		return errors.New("id is required")
	}
	if !strings.HasPrefix(config.ID, "lc_") {
		config.ID = "lc_" + config.ID
	}
	_, err := r.db.Exec("INSERT INTO layer_configs (id, name, description, settings) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET name = $2, description = $3, settings = $4", config.ID, config.Name, config.Description, config.Settings)
	if err != nil {
		return err
	}
	return nil
}

func GetLayerConfigFromFile(id string) (*LayerConfig, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if !strings.HasPrefix(id, "lc_") {
		id = "lc_" + id
	}
	filePath := fmt.Sprintf("layer_configs/%s.json", id)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("layer config file not found")
		}
		return nil, err
	}
	var config LayerConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func StoreLayerConfigToFile(config *LayerConfig) error {
	if config == nil {
		return errors.New("config is required")
	}
	if config.ID == "" {
		return errors.New("id is required")
	}
	if !strings.HasPrefix(config.ID, "lc_") {
		config.ID = "lc_" + config.ID
	}
	filePath := fmt.Sprintf("layer_configs/%s.json", config.ID)
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}