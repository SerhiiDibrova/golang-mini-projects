package initializers

import (
	"encoding/json"
	"net/http"

	"github.com/config-service/config"
)

type ConfigHandler struct {
	configService *config.ConfigService
}

func NewConfigHandler(configService *config.ConfigService) *ConfigHandler {
	return &ConfigHandler{configService: configService}
}

func (c *ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	filteredConfig, err := c.configService.GetConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonConfig, err := json.Marshal(filteredConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonConfig)
}

func (c *ConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		c.GetConfig(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}