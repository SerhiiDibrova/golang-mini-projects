package initializers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"archive/zip"
)

type commonFields struct {
	FieldName      string `json:"fieldName"`
	DownloadFilename string `json:"downloadFilename"`
}

func downloadFilesForLayer(w http.ResponseWriter, r *http.Request) {
	commonFields := _loadCommonFields(r)
	if commonFields.FieldName == "" {
		http.Error(w, "Field name is required", http.StatusBadRequest)
		return
	}
	if commonFields.DownloadFilename == "" {
		http.Error(w, "Download filename is required", http.StatusBadRequest)
		return
	}

	output, err := _executeExternalRequest(commonFields.FieldName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	archive := _generateArchive(output)
	if archive == nil {
		http.Error(w, "Failed to generate archive", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", commonFields.DownloadFilename))
	w.Header().Set("Content-Type", "application/zip")

	_, err = io.Copy(w, archive)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func _loadCommonFields(r *http.Request) commonFields {
	var commonFields commonFields
	err := json.NewDecoder(r.Body).Decode(&commonFields)
	if err != nil {
		log.Println(err)
		return commonFields
	}
	return commonFields
}

func _executeExternalRequest(fieldName string) (io.Reader, error) {
	url := "https://example.com/" + fieldName
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to execute external request: %s", resp.Status)
	}

	return resp.Body, nil
}

func _generateArchive(output io.Reader) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		zipWriter := zip.NewWriter(pw)
		defer zipWriter.Close()
		defer pw.Close()

		_, err := io.Copy(zipWriter, output)
		if err != nil {
			log.Println(err)
			pw.CloseWithError(err)
		}
	}()

	return pr
}