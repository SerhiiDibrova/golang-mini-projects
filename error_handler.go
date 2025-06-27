package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

type ExceptionHandler func(w http.ResponseWriter, r *http.Request, exception error)

type Request struct {
	ExceptionHandler ExceptionHandler
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
}

func exceptionHandler(w http.ResponseWriter, r *http.Request, exception error) {
	if exception != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{
			Code:    http.StatusInternalServerError,
			Message: exception.Error(),
		})
	}
}

func handleDatabaseError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorResponse{
				Code:    http.StatusNotFound,
				Message: "Resource not found",
			})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
		}
	}
}

func renderErrorView(w http.ResponseWriter, r *http.Request, exception error) {
	if exception != nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "<h1>Internal Server Error</h1>")
		fmt.Fprint(w, "<p>An error occurred while processing your request.</p>")
	}
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	req := Request{}
	db, err := sql.Open("postgres", "user=myuser password=mypass dbname=mydb sslmode=disable")
	if err != nil {
		errorHandler(w, r, err)
		return
	}
	defer db.Close()

	_, err = db.Exec("SELECT * FROM mytable")
	if err != nil {
		handleDatabaseError(w, r, err)
		return
	}

	if req.ExceptionHandler != nil {
		req.ExceptionHandler(w, r, nil)
	} else {
		err := http.Error(w, "Exception handler not set", http.StatusInternalServerError)
		if err != nil {
			log.Println(err)
		}
	}
}

func HandleRequestWithExceptionHandler(w http.ResponseWriter, r *http.Request) {
	req := Request{
		ExceptionHandler: func(w http.ResponseWriter, r *http.Request, exception error) {
			if exception != nil {
				exceptionHandler(w, r, exception)
			}
		},
	}
	db, err := sql.Open("postgres", "user=myuser password=mypass dbname=mydb sslmode=disable")
	if err != nil {
		req.ExceptionHandler(w, r, err)
		return
	}
	defer db.Close()

	_, err = db.Exec("SELECT * FROM mytable")
	if err != nil {
		req.ExceptionHandler(w, r, err)
		return
	}

	if req.ExceptionHandler != nil {
		req.ExceptionHandler(w, r, nil)
	} else {
		err := http.Error(w, "Exception handler not set", http.StatusInternalServerError)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	http.HandleFunc("/", HandleRequestWithExceptionHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}