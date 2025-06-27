package initializers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type ExternalRequestService struct{}

func (s *ExternalRequestService) ExecuteExternalRequest(url string, streamProcessor func(io.Reader) (interface{}, error), resultStream io.Writer) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return s.HandleError(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return s.HandleError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return s.HandleError(errors.New("invalid status code"))
	}

	result, err := streamProcessor(resp.Body)
	if err != nil {
		return s.HandleError(err)
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		return s.HandleError(err)
	}

	if resultStream == nil {
		return s.HandleError(errors.New("result stream is nil"))
	}

	_, err = resultStream.Write(jsonResult)
	if err != nil {
		return s.HandleError(err)
	}

	return nil
}

func (s *ExternalRequestService) HandleError(err error) error {
	log.Println(err)
	return err
}