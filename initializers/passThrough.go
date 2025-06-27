package initializers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"

	"your-project/externalRequestService"
	"your-project/downloadTokenService"
	"your-project/proxyService"
	"your-project/requestSingleFieldParamProcessor"
	"your-project/streamProcessor"
	"your-project/urlListStreamProcessor"
	"your-project/calculateSumStreamProcessor"
)

type PassThroughHandler struct {
	proxyService          proxyService.ProxyService
	externalRequestService externalRequestService.ExternalRequestService
	downloadTokenService  downloadTokenService.DownloadTokenService
	mu                    sync.RWMutex
}

func NewPassThroughHandler(proxyService proxyService.ProxyService, externalRequestService externalRequestService.ExternalRequestService, downloadTokenService downloadTokenService.DownloadTokenService) *PassThroughHandler {
	return &PassThroughHandler{
		proxyService:          proxyService,
		externalRequestService: externalRequestService,
		downloadTokenService:  downloadTokenService,
	}
}

func (p *PassThroughHandler) HandlePassThrough(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
		}
	}()

	if p.proxyService == nil {
		http.Error(w, "Proxy service is not initialized", http.StatusInternalServerError)
		return
	}

	if !p.proxyService.IsProxyingAllowed() {
		http.Error(w, "Proxying is not allowed", http.StatusForbidden)
		return
	}

	url := p.proxyService.PerformProxying()
	if url == "" {
		http.Error(w, "Failed to perform proxying", http.StatusInternalServerError)
		return
	}

	streamProcessor := streamProcessor.NewStreamProcessor()
	if streamProcessor == nil {
		http.Error(w, "Failed to create stream processor", http.StatusInternalServerError)
		return
	}

	resultStream := &bytes.Buffer{}

	err := p.externalRequestService.ExecuteExternalRequest(url, streamProcessor, resultStream)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fieldName := "singleField"
	err = requestSingleFieldParamProcessor.Process(fieldName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fieldName = "urlList"
	urlSubstitutions := map[string]string{}
	err = urlListStreamProcessor.Process(fieldName, urlSubstitutions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filenameFieldName := "filename"
	sizeFieldName := "size"
	err = calculateSumStreamProcessor.Process(filenameFieldName, sizeFieldName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = p.downloadTokenService.AddDownloadTokenCookie(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(w, resultStream)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	proxyService := proxyService.NewProxyService()
	externalRequestService := externalRequestService.NewExternalRequestService()
	downloadTokenService := downloadTokenService.NewDownloadTokenService()

	passThroughHandler := NewPassThroughHandler(proxyService, externalRequestService, downloadTokenService)

	http.HandleFunc("/passThrough", passThroughHandler.HandlePassThrough)
	log.Fatal(http.ListenAndServe(":8080", nil))
}