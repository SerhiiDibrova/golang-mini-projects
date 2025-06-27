package initializers

import (
	"encoding/json"
	"fmt"
)

type LayerConfig struct {
	IsBaseLayer       bool `json:"isBaseLayer"`
	Queryable         bool `json:"queryable"`
	IsDataLayer       bool `json:"isDataLayer"`
	DisplayInLayerSwitcher bool `json:"displayInLayerSwitcher"`
	Visibility        bool `json:"visibility"`
}

type BaseLayerConfig struct {
	IsBaseLayer bool `json:"isBaseLayer"`
	Queryable   bool `json:"queryable"`
}

type DataLayerConfig struct {
	IsDataLayer       bool `json:"isDataLayer"`
	DisplayInLayerSwitcher bool `json:"displayInLayerSwitcher"`
	Queryable         bool `json:"queryable"`
	Visibility        bool `json:"visibility"`
}

type Config struct {
	Baselayers []BaseLayerConfig `json:"baselayers"`
	Datalayers []DataLayerConfig `json:"datalayers"`
}

func GetLayerConfig(grailsApplicationConfig string) (*Config, error) {
	var config Config
	err := json.Unmarshal([]byte(grailsApplicationConfig), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func GetBaseLayerConfigs(config *Config) []BaseLayerConfig {
	return config.Baselayers
}

func GetDataLayerConfigs(config *Config) []DataLayerConfig {
	return config.Datalayers
}

func GetLayerConfigFromGrails(grailsApplicationConfig map[string]interface{}) (*Config, error) {
	baselayers, ok := grailsApplicationConfig["baselayers"]
	if !ok {
		return nil, fmt.Errorf("baselayers not found in grailsApplicationConfig")
	}
	datalayers, ok := grailsApplicationConfig["datalayers"]
	if !ok {
		return nil, fmt.Errorf("datalayers not found in grailsApplicationConfig")
	}
	var config Config
	config.Baselayers = make([]BaseLayerConfig, 0)
	config.Datalayers = make([]DataLayerConfig, 0)
	for _, baseLayer := range baselayers.([]interface{}) {
		baseLayerConfig := BaseLayerConfig{
			IsBaseLayer: baseLayer.(map[string]interface{})["isBaseLayer"].(bool),
			Queryable:   baseLayer.(map[string]interface{})["queryable"].(bool),
		}
		config.Baselayers = append(config.Baselayers, baseLayerConfig)
	}
	for _, dataLayer := range datalayers.([]interface{}) {
		dataLayerConfig := DataLayerConfig{
			IsDataLayer:       dataLayer.(map[string]interface{})["isDataLayer"].(bool),
			DisplayInLayerSwitcher: dataLayer.(map[string]interface{})["displayInLayerSwitcher"].(bool),
			Queryable:         dataLayer.(map[string]interface{})["queryable"].(bool),
			Visibility:        dataLayer.(map[string]interface{})["visibility"].(bool),
		}
		config.Datalayers = append(config.Datalayers, dataLayerConfig)
	}
	return &config, nil
}

func main() {
	grailsApplicationConfig := map[string]interface{}{
		"baselayers": []interface{}{
			map[string]interface{}{
				"isBaseLayer": true,
				"queryable":   true,
			},
		},
		"datalayers": []interface{}{
			map[string]interface{}{
				"isDataLayer":       true,
				"displayInLayerSwitcher": true,
				"queryable":         true,
				"visibility":        true,
			},
		},
	}
	config, err := GetLayerConfigFromGrails(grailsApplicationConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	baseLayerConfigs := GetBaseLayerConfigs(config)
	dataLayerConfigs := GetDataLayerConfigs(config)
	fmt.Println(baseLayerConfigs)
	fmt.Println(dataLayerConfigs)
}