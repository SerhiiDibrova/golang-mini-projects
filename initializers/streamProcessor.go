package initializers

import (
	"io"
	"log"
	"os"
)

type StreamProcessor struct {
	Streams map[string]io.Closer
}

func (sp *StreamProcessor) ProcessStream(requestStream io.Reader, resultStream io.Writer) error {
	if requestStream == nil {
		return io.ErrUnexpectedEOF
	}
	if resultStream == nil {
		return io.ErrShortWrite
	}
	_, err := io.Copy(resultStream, requestStream)
	if err != nil {
		log.Printf("Error processing stream: %v", err)
		return err
	}
	return nil
}

func (sp *StreamProcessor) AddStream(name string, stream io.Closer) {
	if sp.Streams == nil {
		sp.Streams = make(map[string]io.Closer)
	}
	sp.Streams[name] = stream
}

func (sp *StreamProcessor) CloseStreams() error {
	if sp.Streams == nil {
		return nil
	}
	for _, stream := range sp.Streams {
		if err := stream.Close(); err != nil {
			log.Printf("Error closing stream: %v", err)
			return err
		}
	}
	return nil
}

func NewStreamProcessor() *StreamProcessor {
	return &StreamProcessor{}
}

func main() {
	streamProcessor := NewStreamProcessor()
	requestStream, err := os.Open("request.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer requestStream.Close()
	resultStream, err := os.Create("result.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer resultStream.Close()
	streamProcessor.AddStream("request", requestStream)
	streamProcessor.AddStream("result", resultStream)
	err = streamProcessor.ProcessStream(requestStream, resultStream)
	if err != nil {
		log.Fatal(err)
	}
	if err := streamProcessor.CloseStreams(); err != nil {
		log.Fatal(err)
	}
}