package main

import (
	"encoding/json"
	"errors"
	"testing"
)

type Filter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  map[string]string `json:"parameters"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

func (f *Filter) UnmarshalJSON(data []byte) error {
	var aux struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Parameters  map[string]string `json:"parameters"`
		CreatedAt   string            `json:"created_at"`
		UpdatedAt   string            `json:"updated_at"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	f.Name = aux.Name
	f.Description = aux.Description
	f.Parameters = aux.Parameters
	f.CreatedAt = aux.CreatedAt
	f.UpdatedAt = aux.UpdatedAt
	return nil
}

func (f Filter) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Parameters  map[string]string `json:"parameters"`
		CreatedAt   string            `json:"created_at"`
		UpdatedAt   string            `json:"updated_at"`
	}{
		Name:        f.Name,
		Description: f.Description,
		Parameters:  f.Parameters,
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	})
}

func NewFilter(name string, description string, parameters map[string]string) (*Filter, error) {
	if name == "" {
		return nil, errors.New("filter name cannot be empty")
	}
	if description == "" {
		return nil, errors.New("filter description cannot be empty")
	}
	if parameters == nil {
		return nil, errors.New("filter parameters cannot be nil")
	}
	return &Filter{
		Name:        name,
		Description: description,
		Parameters:  parameters,
	}, nil
}

func (f *Filter) Validate() error {
	if f.Name == "" {
		return errors.New("filter name cannot be empty")
	}
	if f.Description == "" {
		return errors.New("filter description cannot be empty")
	}
	if f.Parameters == nil || len(f.Parameters) == 0 {
		return errors.New("filter parameters cannot be empty")
	}
	return nil
}

func TestNewFilter(t *testing.T) {
	_, err := NewFilter("", "", nil)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	_, err = NewFilter("test", "", nil)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	_, err = NewFilter("test", "test", nil)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	_, err = NewFilter("test", "test", map[string]string{})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	_, err = NewFilter("test", "test", map[string]string{"key": "value"})
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestFilterValidate(t *testing.T) {
	f := &Filter{}
	err := f.Validate()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	f.Name = "test"
	err = f.Validate()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	f.Description = "test"
	err = f.Validate()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	f.Parameters = map[string]string{}
	err = f.Validate()
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	f.Parameters = map[string]string{"key": "value"}
	err = f.Validate()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestFilterUnmarshalJSON(t *testing.T) {
	jsonData := `{"name": "test", "description": "test", "parameters": {"key": "value"}}`
	var f Filter
	err := json.Unmarshal([]byte(jsonData), &f)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if f.Name != "test" {
		t.Errorf("expected test, got %s", f.Name)
	}
	if f.Description != "test" {
		t.Errorf("expected test, got %s", f.Description)
	}
	if f.Parameters["key"] != "value" {
		t.Errorf("expected value, got %s", f.Parameters["key"])
	}
}

func TestFilterMarshalJSON(t *testing.T) {
	f := Filter{
		Name:        "test",
		Description: "test",
		Parameters:  map[string]string{"key": "value"},
	}
	jsonData, err := json.Marshal(f)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	var aux struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Parameters  map[string]string `json:"parameters"`
		CreatedAt   string            `json:"created_at"`
		UpdatedAt   string            `json:"updated_at"`
	}
	err = json.Unmarshal(jsonData, &aux)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if aux.Name != "test" {
		t.Errorf("expected test, got %s", aux.Name)
	}
	if aux.Description != "test" {
		t.Errorf("expected test, got %s", aux.Description)
	}
	if aux.Parameters["key"] != "value" {
		t.Errorf("expected value, got %s", aux.Parameters["key"])
	}
}

func main() {
	filter, err := NewFilter("test", "test filter", map[string]string{"key": "value"})
	if err != nil {
		panic(err)
	}
	jsonData, err := json.Marshal(filter)
	if err != nil {
		panic(err)
	}
	println(string(jsonData))
}