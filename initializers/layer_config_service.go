package initializers

import (
	"errors"
	"log"
)

type LayerConfigService interface {
	GetLayerConfig() ([]LayerConfig, error)
	ProcessBaseLayers([]LayerConfig) []LayerConfig
	ProcessDataLayers([]LayerConfig) []LayerConfig
	CombineLayerConfigs([]LayerConfig, []LayerConfig) []LayerConfig
	FilterOutNullLayers([]LayerConfig) []LayerConfig
}

type layerConfigService struct{}

func (l *layerConfigService) GetLayerConfig() ([]LayerConfig, error) {
	layerConfigs, err := retrieveFromDatabase()
	if err != nil {
		return nil, err
	}
	if layerConfigs == nil {
		return nil, errors.New("layer configurations are nil")
	}
	return layerConfigs, nil
}

func (l *layerConfigService) ProcessBaseLayers(layerConfigs []LayerConfig) []LayerConfig {
	baseLayers := make([]LayerConfig, len(layerConfigs))
	copy(baseLayers, layerConfigs)
	for i := range baseLayers {
		baseLayers[i].IsBaseLayer = true
		baseLayers[i].Queryable = false
	}
	return baseLayers
}

func (l *layerConfigService) ProcessDataLayers(layerConfigs []LayerConfig) []LayerConfig {
	dataLayers := make([]LayerConfig, len(layerConfigs))
	copy(dataLayers, layerConfigs)
	for i := range dataLayers {
		dataLayers[i].IsDataLayer = true
		dataLayers[i].DisplayInLayerSwitcher = true
		dataLayers[i].Queryable = false
		dataLayers[i].Visibility = false
	}
	return dataLayers
}

func (l *layerConfigService) CombineLayerConfigs(baseLayers []LayerConfig, dataLayers []LayerConfig) []LayerConfig {
	if baseLayers == nil {
		baseLayers = []LayerConfig{}
	}
	if dataLayers == nil {
		dataLayers = []LayerConfig{}
	}
	return append(baseLayers, dataLayers...)
}

func (l *layerConfigService) FilterOutNullLayers(layerConfigs []LayerConfig) []LayerConfig {
	if layerConfigs == nil {
		return []LayerConfig{}
	}
	var filteredLayerConfigs []LayerConfig
	for _, layerConfig := range layerConfigs {
		if layerConfig != (LayerConfig{}) {
			filteredLayerConfigs = append(filteredLayerConfigs, layerConfig)
		}
	}
	return filteredLayerConfigs
}

type LayerConfig struct {
	IsBaseLayer          bool
	IsDataLayer          bool
	DisplayInLayerSwitcher bool
	Queryable            bool
	Visibility           bool
	Name                 string
}

func retrieveFromDatabase() ([]LayerConfig, error) {
	layerConfigs := []LayerConfig{
		{IsBaseLayer: true, Name: "Base Layer"},
		{IsDataLayer: true, Name: "Data Layer"},
	}
	return layerConfigs, nil
}

func NewLayerConfigService() LayerConfigService {
	return &layerConfigService{}
}

func main() {
	layerConfigService := NewLayerConfigService()
	layerConfigs, err := layerConfigService.GetLayerConfig()
	if err != nil {
		log.Fatal(err)
	}
	baseLayers := layerConfigService.ProcessBaseLayers(layerConfigs)
	dataLayers := layerConfigService.ProcessDataLayers(layerConfigs)
	combinedLayers := layerConfigService.CombineLayerConfigs(baseLayers, dataLayers)
	filteredLayers := layerConfigService.FilterOutNullLayers(combinedLayers)
	log.Println(filteredLayers)
}