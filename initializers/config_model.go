package initializers

import (
	"errors"
	"fmt"
)

type Config struct {
	Key   string
	Value string
}

func NewConfig(key string, value string) (*Config, error) {
	if key == "" || value == "" {
		return nil, errors.New("key and value cannot be empty")
	}
	return &Config{Key: key, Value: value}, nil
}

func (c *Config) GetKey() string {
	return c.Key
}

func (c *Config) GetValue() string {
	return c.Value
}

func (c *Config) SetKey(key string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	c.Key = key
	return nil
}

func (c *Config) SetValue(value string) error {
	if value == "" {
		return errors.New("value cannot be empty")
	}
	c.Value = value
	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf("Config{Key: %s, Value: %s}", c.Key, c.Value)
}