package initializers

import (
	"errors"
	"log"
)

type ServerClassFactory struct{}

type ServerClass interface {
	Start() error
	Stop() error
}

type ServerType string

const (
	HTTPServer  ServerType = "http"
	GRPCServer  ServerType = "grpc"
	WebServer   ServerType = "web"
	UnknownType ServerType = "unknown"
)

type HTTPServer struct{}

func (h *HTTPServer) Start() error {
	log.Println("HTTP server started")
	return nil
}

func (h *HTTPServer) Stop() error {
	log.Println("HTTP server stopped")
	return nil
}

type GRPCServer struct{}

func (g *GRPCServer) Start() error {
	log.Println("GRPC server started")
	return nil
}

func (g *GRPCServer) Stop() error {
	log.Println("GRPC server stopped")
	return nil
}

type WebServer struct{}

func (w *WebServer) Start() error {
	log.Println("Web server started")
	return nil
}

func (w *WebServer) Stop() error {
	log.Println("Web server stopped")
	return nil
}

func NewServerClassFactory() *ServerClassFactory {
	return &ServerClassFactory{}
}

func (s *ServerClassFactory) CreateServerClass(serverType ServerType) (ServerClass, error) {
	if serverType == UnknownType {
		return nil, errors.New("unknown server type")
	}
	switch serverType {
	case HTTPServer:
		return &HTTPServer{}, nil
	case GRPCServer:
		return &GRPCServer{}, nil
	case WebServer:
		return &WebServer{}, nil
	default:
		return nil, errors.New("unknown server type")
	}
}