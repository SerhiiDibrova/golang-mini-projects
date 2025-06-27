package initializers

import (
	"encoding/json"
	"log"
	"net/http"

	"your-project/layer-config-service"
)

type LayerConfigController struct {
	layerConfigService layer_config_service.LayerConfigService
}

func NewLayerConfigController(layerConfigService layer_config_service.LayerConfigService) *LayerConfigController {
	return &LayerConfigController{layerConfigService: layerConfigService}
}

func (c *LayerConfigController) GetLayerConfig(w http.ResponseWriter, r *http.Request) {
	if c.layerConfigService == nil {
		http.Error(w, "Layer config service is not initialized", http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	layerConfig, err := c.layerConfigService.GetLayerConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if layerConfig == nil {
		http.Error(w, "Layer config is not found", http.StatusNotFound)
		return
	}

	baseLayers, err := c.layerConfigService.ProcessBaseLayers(layerConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dataLayers, err := c.layerConfigService.ProcessDataLayers(layerConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	combinedLayers, err := c.layerConfigService.CombineLayerConfigs(baseLayers, dataLayers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filteredLayers, err := c.layerConfigService.FilterOutNullLayers(combinedLayers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if filteredLayers == nil {
		http.Error(w, "Filtered layers are not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(filteredLayers); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}