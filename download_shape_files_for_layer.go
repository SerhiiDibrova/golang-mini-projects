package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

type CommonFields struct {
	FieldName string `json:"field_name"`
}

func _loadCommonFields(r *http.Request) (CommonFields, error) {
	commonFields := CommonFields{}
	fieldName := r.URL.Query().Get("field_name")
	if fieldName == "" {
		return commonFields, errors.New("field name is not provided")
	}
	commonFields.FieldName = fieldName
	return commonFields, nil
}

func _addDownloadTokenCookie(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:  "download_token",
		Value: "token_value",
	}
	http.SetCookie(w, cookie)
}

func downloadShapeFilesForLayer(w http.ResponseWriter, r *http.Request) {
	commonFields, err := _loadCommonFields(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if commonFields.FieldName == "" {
		http.Error(w, "Field name is not provided", http.StatusBadRequest)
		return
	}
	shapeFileURL := fmt.Sprintf("https://example.com/%s/shapefile.zip", commonFields.FieldName)
	resp, err := http.Get(shapeFileURL)
	if err != nil {
		http.Error(w, "Error downloading shape file", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Shape file not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(shapeFileURL)))
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Error writing shape file to response", http.StatusInternalServerError)
		return
	}
	_addDownloadTokenCookie(w, r)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/downloadShapeFilesForLayer", downloadShapeFilesForLayer).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}