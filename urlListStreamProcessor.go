package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

type UrlListStream struct {
	UrlList []string `json:"url_list"`
}

type FieldNameSubstitution struct {
	FieldName string `json:"field_name"`
	Substitution string `json:"substitution"`
}

type UrlSubstitution struct {
	Url string `json:"url"`
	Substitution string `json:"substitution"`
}

type Config struct {
	FieldNameSubstitutions []FieldNameSubstitution `json:"field_name_substitutions"`
	UrlSubstitutions []UrlSubstitution `json:"url_substitutions"`
}

func validateConfig(config Config) error {
	if config.FieldNameSubstitutions == nil || config.UrlSubstitutions == nil {
		return errors.New("config is invalid")
	}
	return nil
}

func validateUrlListStream(urlListStream UrlListStream) error {
	if urlListStream.UrlList == nil {
		return errors.New("url list stream is invalid")
	}
	return nil
}

func urlListStreamProcessor(jsonStream []byte, config Config) (func() []string, error) {
	if len(jsonStream) == 0 {
		return nil, errors.New("json stream is empty")
	}

	err := validateConfig(config)
	if err != nil {
		return nil, err
	}

	var urlListStream UrlListStream
	err = json.Unmarshal(jsonStream, &urlListStream)
	if err != nil {
		return nil, err
	}

	err = validateUrlListStream(urlListStream)
	if err != nil {
		return nil, err
	}

	processedUrlList := make([]string, len(urlListStream.UrlList))
	for i, url := range urlListStream.UrlList {
		for _, fieldNameSubstitution := range config.FieldNameSubstitutions {
			url = strings.ReplaceAll(url, fieldNameSubstitution.FieldName, fieldNameSubstitution.Substitution)
		}
		for _, urlSubstitution := range config.UrlSubstitutions {
			url = strings.ReplaceAll(url, urlSubstitution.Url, urlSubstitution.Substitution)
		}
		processedUrlList[i] = url
	}

	return func() []string {
		return processedUrlList
	}, nil
}

func main() {
	jsonStream := []byte(`{"url_list": ["http://example.com/field1", "http://example.com/field2"]}`)
	config := Config{
		FieldNameSubstitutions: []FieldNameSubstitution{
			{FieldName: "field1", Substitution: "value1"},
			{FieldName: "field2", Substitution: "value2"},
		},
		UrlSubstitutions: []UrlSubstitution{
			{Url: "http://example.com", Substitution: "https://new-example.com"},
		},
	}

	processor, err := urlListStreamProcessor(jsonStream, config)
	if err != nil {
		log.Fatal(err)
	}

	processedUrlList := processor()
	fmt.Println(processedUrlList)
}