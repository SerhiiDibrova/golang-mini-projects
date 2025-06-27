package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type URLList struct {
	Urls []string `json:"urls"`
}

func _executeExternalRequest(urlString string, maxRetries int, retryDelay time.Duration) (*URLList, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, errors.New("invalid URL")
	}

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	var resp *http.Response
	var retryCount int
	for {
		resp, err = client.Do(req)
		if err != nil {
			if retryCount < maxRetries {
				retryCount++
				time.Sleep(retryDelay)
				continue
			}
			return nil, err
		}
		break
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve data. Status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var urlList URLList
	err = json.Unmarshal(body, &urlList)
	if err != nil {
		return nil, err
	}

	return &urlList, nil
}

func main() {
	url := "http://example.com/urls"
	urlList, err := _executeExternalRequest(url, 3, 500*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}

	for _, url := range urlList.Urls {
		fmt.Println(url)
	}
}