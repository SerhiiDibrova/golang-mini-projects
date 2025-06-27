package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"
)

var (
	templateEngine *template.Template
	cache          = make(map[string]string)
	cacheMutex     sync.RWMutex
	logger         *log.Logger
)

func init() {
	var err error
	templateEngine, err = template.ParseFiles("error_template.html")
	if err != nil {
		log.Fatal(err)
	}
	logger = log.New(log.Writer(), "render_execution_failed: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func _renderExecutionFailed(w http.ResponseWriter, errorMessage string) {
	if errorMessage == "" {
		http.Error(w, "error message is empty", http.StatusBadRequest)
		logger.Println("error message is empty")
		return
	}

	cacheMutex.RLock()
	if cachedView, ok := cache[errorMessage]; ok {
		cacheMutex.RUnlock()
		fmt.Fprint(w, cachedView)
		return
	}
	cacheMutex.RUnlock()

	var view string
	err := templateEngine.ExecuteTemplate(&viewBuffer{w}, "error_template.html", struct {
		ErrorMessage string
	}{
		ErrorMessage: errorMessage,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Println(err)
		return
	}

	cacheMutex.Lock()
	cache[errorMessage] = view
	if len(cache) > 100 {
		for key := range cache {
			delete(cache, key)
			break
		}
	}
	cacheMutex.Unlock()
}

type viewBuffer struct {
	w http.ResponseWriter
}

func (v *viewBuffer) Write(p []byte) (int, error) {
	v.w.Write(p)
	return len(p), nil
}

func main() {
	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		_renderExecutionFailed(w, "Something went wrong")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}