package wps

import (
	"encoding/xml"
	"log"
	"net/http"
	"time"

	"helpers"
)

type ExecutionStatusResponse struct {
	XMLName xml.Name `xml:"wps:ExecuteResponse"`
	Status  Status   `xml:"wps:Status"`
}

type Status struct {
	Value string `xml:"ows:Value"`
}

func GetExecutionStatusResponse(url string) (*ExecutionStatusResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		helpers.LogError(err)
		return nil, err
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		helpers.LogError(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		helpers.LogError(err)
		return nil, err
	}

	var executionStatusResponse ExecutionStatusResponse
	err = xml.NewDecoder(resp.Body).Decode(&executionStatusResponse)
	if err != nil {
		helpers.LogError(err)
		return nil, err
	}

	return &executionStatusResponse, nil
}