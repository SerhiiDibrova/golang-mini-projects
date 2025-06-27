package initializers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

type Coordinates struct {
	Lat float64 `xml:"lat"`
	Lon float64 `xml:"lon"`
}

func _generateServiceUrl(lat, lon float64) (string, error) {
	if lat < -90 || lat > 90 {
		return "", errors.New("invalid latitude")
	}
	if lon < -180 || lon > 180 {
		return "", errors.New("invalid longitude")
	}
	return fmt.Sprintf("https://example.com/service?lat=%f&lon=%f", lat, lon), nil
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")

	if lat == "" || lon == "" {
		http.Error(w, "Missing lat or lon parameter", http.StatusBadRequest)
		return
	}

	if !regexp.MustCompile(`^-?\d{1,3}\.?\d{0,10}$`).MatchString(lat) {
		http.Error(w, "Invalid lat parameter", http.StatusBadRequest)
		return
	}

	if !regexp.MustCompile(`^-?\d{1,3}\.?\d{0,10}$`).MatchString(lon) {
		http.Error(w, "Invalid lon parameter", http.StatusBadRequest)
		return
	}

	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		http.Error(w, "Invalid lat parameter", http.StatusBadRequest)
		return
	}

	lonFloat, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		http.Error(w, "Invalid lon parameter", http.StatusBadRequest)
		return
	}

	serviceUrl, err := _generateServiceUrl(latFloat, lonFloat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := http.Get(serviceUrl)
	if err != nil {
		http.Error(w, "Failed to call service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Service returned error", http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}