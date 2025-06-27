package wps

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"helpers"
)

type WpsService struct {
	url string
}

type ExecutionStatusResponse struct {
	XMLName xml.Name `xml:"wps:ExecuteResponse"`
	Status  string   `xml:"wps:Status>ows:Value"`
}

func NewWpsService(url string) *WpsService {
	return &WpsService{url: url}
}

func (s *WpsService) _getExecutionStatusResponse(executionId string) (*ExecutionStatusResponse, error) {
	req, err := http.NewRequest("GET", s.url+"/jobs/"+executionId, nil)
	if err != nil {
		helpers.LogError(err)
		return nil, err
	}

	req.Header.Set("Accept", "application/xml")

	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			helpers.LogError(errors.New("connection timeout"))
			return nil, errors.New("connection timeout")
		}
		helpers.LogError(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		helpers.LogError(errors.New("non-200 status code: " + resp.Status))
		return nil, errors.New("non-200 status code: " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.LogError(err)
		return nil, err
	}

	var executionStatusResponse ExecutionStatusResponse
	err = xml.Unmarshal(body, &executionStatusResponse)
	if err != nil {
		helpers.LogError(err)
		return nil, err
	}

	return &executionStatusResponse, nil
}

func GetExecutionStatusResponse(executionId string) (*ExecutionStatusResponse, error) {
	wpsService := NewWpsService("http://localhost:8080/wps")
	return wpsService._getExecutionStatusResponse(executionId)
}