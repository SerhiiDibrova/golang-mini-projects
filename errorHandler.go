package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ErrorHandler struct {
	exceptionHandlers map[string]func(error) (bool, string)
	logger            *log.Logger
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		exceptionHandlers: make(map[string]func(error) (bool, string)),
		logger:            log.Default(),
	}
}

func (e *ErrorHandler) AddExceptionHandler(exceptionType string, handler func(error) (bool, string)) {
	e.exceptionHandlers[exceptionType] = handler
}

func (e *ErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	e.LogError(err)
	for exceptionType, handler := range e.exceptionHandlers {
		if errors.Is(err, &exceptionType{}) {
			handled, message := handler(err)
			if handled {
				e.RenderErrorMessage(w, r, message)
				return
			}
		}
	}
	e.RenderErrorMessage(w, r, err.Error())
}

func (e *ErrorHandler) RenderErrorMessage(w http.ResponseWriter, r *http.Request, message string) {
	errorMessage := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(errorMessage)
}

func (e *ErrorHandler) LogError(err error) {
	e.logger.Printf("error: %v", err)
}

func main() {
	errorHandler := NewErrorHandler()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := func() error {
			// endpoint logic
			return nil
		}()
		if err != nil {
			errorHandler.HandleError(w, r, err)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}