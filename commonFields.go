package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type CommonField struct {
	FieldName        string `json:"fieldName"`
	URLSubstitutions string `json:"urlSubstitutions"`
}

type RequestParams struct {
	CommonFields []CommonField `json:"commonFields"`
}

func _loadCommonFields(r *http.Request) ([]CommonField, error) {
	if r.Body == nil {
		return nil, errors.New("request body is nil")
	}
	var requestParams RequestParams
	err := json.NewDecoder(r.Body).Decode(&requestParams)
	if err != nil {
		return nil, err
	}
	if requestParams.CommonFields == nil {
		return nil, errors.New("common fields are not properly formatted")
	}
	for _, field := range requestParams.CommonFields {
		if field.FieldName == "" || field.URLSubstitutions == "" {
			return nil, errors.New("common fields are not properly formatted")
		}
	}
	return requestParams.CommonFields, nil
}

func main() {
	r, err := http.NewRequest("POST", "", nil)
	if err != nil {
		log.Fatal(err)
	}
	r.Body = nil
	commonFields, err := _loadCommonFields(r)
	if err != nil {
		http.Error(nil, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}
	log.Println(commonFields)
}