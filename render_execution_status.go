package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"text/template"
	"time"
)

type ExecutionStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type RenderExecutionStatus struct {
	templateEngine *template.Template
	cache          map[string]string
	cacheMutex     sync.RWMutex
	logger         *log.Logger
}

func NewRenderExecutionStatus(templateEngine *template.Template, logger *log.Logger) *RenderExecutionStatus {
	return &RenderExecutionStatus{
		templateEngine: templateEngine,
		cache:          make(map[string]string),
		logger:         logger,
	}
}

func (r *RenderExecutionStatus) _renderExecutionStatus(executionStatus ExecutionStatus) (string, error) {
	if executionStatus.Status == "" || executionStatus.Message == "" {
		return "", errors.New("execution status response is empty")
	}

	if r.templateEngine == nil {
		r.logger.Printf("error: template engine is unavailable")
		return "", errors.New("template engine is unavailable")
	}

	var renderedStatus string
	var err error

	r.cacheMutex.RLock()
	if cachedStatus, ok := r.cache[executionStatus.Status]; ok {
		renderedStatus = cachedStatus
		r.cacheMutex.RUnlock()
		return renderedStatus, nil
	}
	r.cacheMutex.RUnlock()

	data := map[string]string{
		"Status":  executionStatus.Status,
		"Message": executionStatus.Message,
	}

	var buf []byte
	buf, err = r.templateEngine.ExecuteTemplate([]byte{}, "executionStatus", data)
	if err != nil {
		r.logger.Printf("error rendering execution status: %v", err)
		return "", err
	}

	renderedStatus = string(buf)

	r.cacheMutex.Lock()
	r.cache[executionStatus.Status] = renderedStatus
	if len(r.cache) > 100 {
		keys := make([]string, 0, len(r.cache))
		for k := range r.cache {
			keys = append(keys, k)
		}
		r.cache[keys[0]] = ""
		delete(r.cache, keys[0])
	}
	r.cacheMutex.Unlock()

	return renderedStatus, nil
}

func main() {
	logger := log.New(log.Writer(), "render-execution-status: ", log.Ldate|log.Ltime|log.Lshortfile)
	templateEngine := template.Must(template.New("executionStatus").Parse("Status: {{ .Status }}, Message: {{ .Message }}"))
	renderExecutionStatus := NewRenderExecutionStatus(templateEngine, logger)

	executionStatus := ExecutionStatus{
		Status:  "success",
		Message: "execution completed successfully",
	}

	renderedStatus, err := renderExecutionStatus._renderExecutionStatus(executionStatus)
	if err != nil {
		logger.Printf("error rendering execution status: %v", err)
		return
	}

	fmt.Println(renderedStatus)

	time.Sleep(1 * time.Second)

	renderedStatus, err = renderExecutionStatus._renderExecutionStatus(executionStatus)
	if err != nil {
		logger.Printf("error rendering execution status: %v", err)
		return
	}

	fmt.Println(renderedStatus)
}