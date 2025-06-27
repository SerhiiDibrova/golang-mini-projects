package helpers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ErrorHandler struct {
	logger *log.Logger
}

type ErrorLogger struct {
	logger *log.Logger
}

type ErrorResponse struct {
	StatusCode    int    `json:"statusCode"`
	ErrorMessage  string `json:"errorMessage"`
	LayerName     string `json:"layerName"`
}

type ErrorHandlingLogic struct {
	errorHandler *ErrorHandler
}

func NewErrorLogger(logger *log.Logger) *ErrorLogger {
	return &ErrorLogger{logger: logger}
}

func NewErrorHandler(logger *log.Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

func NewErrorHandlingLogic(errorHandler *ErrorHandler) *ErrorHandlingLogic {
	return &ErrorHandlingLogic{errorHandler: errorHandler}
}

func (e *ErrorLogger) LogError(err error) {
	e.logger.Println(err)
}

func (e *ErrorHandler) HandleError(w http.ResponseWriter, err error, layerName string) {
	if layerName == "" {
		layerName = "unknown"
	}

	e.logger.Println(err)
	var statusCode int
	var errorMessage string

	switch err.Error() {
	case "server configuration error":
		statusCode = http.StatusInternalServerError
		errorMessage = "Internal Server Error"
	case "knownServers configuration error":
		statusCode = http.StatusBadGateway
		errorMessage = "Bad Gateway"
	default:
		statusCode = http.StatusInternalServerError
		errorMessage = "Internal Server Error"
	}

	errorResponse := ErrorResponse{
		StatusCode:    statusCode,
		ErrorMessage:  errorMessage,
		LayerName:     layerName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

func (e *ErrorHandlingLogic) HandleClientLogError(w http.ResponseWriter, err error, layerName string) {
	if err != nil {
		e.errorHandler.HandleError(w, err, layerName)
	}
}

func (e *ErrorHandlingLogic) HandleNotificationError(w http.ResponseWriter, err error, layerName string) {
	if err != nil {
		e.errorHandler.HandleError(w, err, layerName)
	}
}

func (e *ErrorHandlingLogic) HandleLoggingError(w http.ResponseWriter, err error, layerName string) {
	if err != nil {
		e.errorHandler.HandleError(w, err, layerName)
	}
}

func (e *ErrorHandlingLogic) HandleServerError(w http.ResponseWriter, err error, layerName string) {
	if err != nil {
		e.errorHandler.HandleError(w, errors.New("server configuration error"), layerName)
	}
}

func (e *ErrorHandlingLogic) HandleKnownServersError(w http.ResponseWriter, err error, layerName string) {
	if err != nil {
		e.errorHandler.HandleError(w, errors.New("knownServers configuration error"), layerName)
	}
}