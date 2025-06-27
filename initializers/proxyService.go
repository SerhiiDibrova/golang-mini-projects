package initializers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/externalRequestService"
)

type ProxyService struct {
	externalRequestService externalRequestService.ExternalRequestService
	systemConfig            SystemConfig
}

type SystemConfig struct {
	ProxyEnabled bool
	ProxyUrl     string
}

type RequestParams struct {
	Url string
}

func NewProxyService(externalRequestService externalRequestService.ExternalRequestService, systemConfig SystemConfig) *ProxyService {
	return &ProxyService{
		externalRequestService: externalRequestService,
		systemConfig:            systemConfig,
	}
}

func (p *ProxyService) Proxy(w http.ResponseWriter, r *http.Request) {
	if r == nil {
		http.Error(w, "Request is nil", http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		http.Error(w, "Request body is nil", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	requestParams := RequestParams{}
	err := json.NewDecoder(r.Body).Decode(&requestParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !p.systemConfig.ProxyEnabled {
		http.Error(w, "Proxying is not enabled", http.StatusForbidden)
		return
	}

	url := requestParams.Url
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	if p.systemConfig.ProxyUrl != "" && url != p.systemConfig.ProxyUrl {
		http.Error(w, "Proxy URL does not match", http.StatusForbidden)
		return
	}

	streamProcessor := NewStreamProcessor()
	resultStream := NewResultStream()

	err = p.externalRequestService.ExecuteExternalRequest(url, streamProcessor, resultStream)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, resultStream)
	if err != nil {
		log.Println(err)
	}
}

type StreamProcessor struct{}

func NewStreamProcessor() *StreamProcessor {
	return &StreamProcessor{}
}

func (s *StreamProcessor) Process(stream io.Reader) (io.Reader, error) {
	if stream == nil {
		return nil, errors.New("stream is nil")
	}
	return stream, nil
}

type ResultStream struct {
	stream io.Reader
}

func NewResultStream() *ResultStream {
	return &ResultStream{}
}

func (r *ResultStream) Read(p []byte) (n int, err error) {
	if r.stream == nil {
		return 0, errors.New("stream is not set")
	}
	return r.stream.Read(p)
}

func (r *ResultStream) SetStream(stream io.Reader) {
	r.stream = stream
}