package initializers

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorHandler(w http.ResponseWriter, err error, code int) {
	if w == nil {
		log.Println("http.ResponseWriter is nil")
		return
	}
	if err == nil {
		log.Println("error is nil")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	errorResponse := &ErrorResponse{
		Code:    code,
		Message: err.Error(),
	}
	jsonError, err := json.Marshal(errorResponse)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
	_, err = w.Write(jsonError)
	if err != nil {
		log.Println(err)
	}
}

func HandleMissingParameter(w http.ResponseWriter, param string) {
	if w == nil {
		log.Println("http.ResponseWriter is nil")
		return
	}
	if param == "" {
		log.Println("param is empty")
		return
	}
	err := &ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "Missing parameter: " + param,
	}
	jsonError, err := json.Marshal(err)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(jsonError)
	if err != nil {
		log.Println(err)
	}
}

func HandleInvalidParameter(w http.ResponseWriter, param string) {
	if w == nil {
		log.Println("http.ResponseWriter is nil")
		return
	}
	if param == "" {
		log.Println("param is empty")
		return
	}
	err := &ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "Invalid parameter: " + param,
	}
	jsonError, err := json.Marshal(err)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(jsonError)
	if err != nil {
		log.Println(err)
	}
}

func HandleServiceUnavailable(w http.ResponseWriter) {
	if w == nil {
		log.Println("http.ResponseWriter is nil")
		return
	}
	err := &ErrorResponse{
		Code:    http.StatusServiceUnavailable,
		Message: "Service is currently unavailable",
	}
	jsonError, err := json.Marshal(err)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusServiceUnavailable)
	_, err = w.Write(jsonError)
	if err != nil {
		log.Println(err)
	}
}

func HandleDatabaseConnectionError(w http.ResponseWriter, err error) {
	if w == nil {
		log.Println("http.ResponseWriter is nil")
		return
	}
	if err == nil {
		log.Println("error is nil")
		return
	}
	errorResponse := &ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: "Database connection error: " + err.Error(),
	}
	jsonError, err := json.Marshal(errorResponse)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write(jsonError)
	if err != nil {
		log.Println(err)
	}
}