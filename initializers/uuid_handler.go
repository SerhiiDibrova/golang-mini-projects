package initializers

import (
	"net/http"
	"strings"
)

func uuidHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("id")
	redirectURL := search_service(uuid)
	w.Header().Set("Location", redirectURL)
	w.WriteHeader(http.StatusFound)
}

func search_service(uuid string) string {
	baseURL := "https://example.com/"
	path := "/redirect/"
	query := "?param1=value1&param2=value2"
	return baseURL + path + uuid + query
}

func init() {
	http.HandleFunc("/uuid", uuidHandler)
}