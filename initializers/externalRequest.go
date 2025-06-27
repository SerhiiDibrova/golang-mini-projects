package initializers

import (
	"io"
	"net/http"
)

type StreamProcessor func(io.Reader) (io.Reader, error)

func urlListStreamProcessor(reader io.Reader) (io.Reader, error) {
	return reader, nil
}

func _executeExternalRequest(url string, streamProcessor StreamProcessor, resultStream io.Writer) error {
	if streamProcessor == nil {
		return io.ErrUnexpectedEOF
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return http.ErrFailedLoad
	}

	processedStream, err := streamProcessor(resp.Body)
	if err != nil {
		return err
	}

	_, err = io.Copy(resultStream, processedStream)
	if err != nil {
		return err
	}

	return nil
}