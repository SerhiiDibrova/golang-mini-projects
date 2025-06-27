package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type RequestParams struct {
	Url      string `json:"url"`
	ProxyUrl string `json:"proxy_url"`
}

func _performProxyingIfAllowed(params RequestParams) bool {
	if params.ProxyUrl == "" || params.Url == "" {
		return false
	}

	req, err := http.NewRequest("GET", params.Url, nil)
	if err != nil {
		log.Println(err)
		return false
	}

	req.Header.Set("User-Agent", "ProxyClient")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to retrieve data from URL")
		return false
	}

	proxyReq, err := http.NewRequest("POST", params.ProxyUrl, resp.Body)
	if err != nil {
		log.Println(err)
		return false
	}

	proxyReq.Header.Set("User-Agent", "ProxyClient")
	proxyReq.Header.Set("Content-Type", resp.Header.Get("Content-Type"))

	proxyClient := &http.Client{}
	proxyResp, err := proxyClient.Do(proxyReq)
	if err != nil {
		log.Println(err)
		return false
	}
	defer proxyResp.Body.Close()

	if proxyResp.StatusCode != http.StatusOK {
		log.Println("Failed to send data to proxy URL")
		return false
	}

	_, err = io.Copy(io.Discard, proxyResp.Body)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func main() {
	params := RequestParams{
		Url:      "http://example.com",
		ProxyUrl: "http://proxy.example.com",
	}

	result := _performProxyingIfAllowed(params)
	fmt.Println(result)
}