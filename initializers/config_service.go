package initializers

import (
	"errors"
	"fmt"
	"reflect"
)

type ConfigRepository interface {
	GetConfig() (map[string]interface{}, error)
}

type ConfigService struct {
	configRepository ConfigRepository
}

var CONFIG_KEYS_TO_IGNORE = []string{"key1", "key2"}

func NewConfigService(configRepository ConfigRepository) *ConfigService {
	return &ConfigService{configRepository: configRepository}
}

func (c *ConfigService) GetConfig() (map[string]interface{}, error) {
	config, err := c.configRepository.GetConfig()
	if err != nil {
		return nil, err
	}

	filteredConfig := make(map[string]interface{})
	for key, value := range config {
		if !c.isClosure(value) && !c.isKeyIgnored(key) {
			filteredConfig[key] = value
		}
	}

	if len(filteredConfig) == 0 {
		return nil, errors.New("no configuration found")
	}

	return filteredConfig, nil
}

func (c *ConfigService) isClosure(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.Func
}

func (c *ConfigService) isKeyIgnored(key string) bool {
	if CONFIG_KEYS_TO_IGNORE == nil {
		return false
	}
	for _, ignoredKey := range CONFIG_KEYS_TO_IGNORE {
		if key == ignoredKey {
			return true
		}
	}
	return false
}