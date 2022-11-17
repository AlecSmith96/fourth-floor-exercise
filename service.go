package main

import (
	"context"
	"log"
	"net/http"
)

// RESTServer wraps the http server functionality for startup and graceful shutdown
type RESTServer interface {
	Start() error
	Shutdown(context.Context) error
}

// Service top level struct containing any injected dependencies
type Service struct {
	Server *http.Server
	Logger *log.Logger
}

func (s *Service) Start() error {
	s.Logger.Println("Starting the service")
	return s.Server.ListenAndServe()
} 

func (s *Service) Shutdown(ctx context.Context) error {
	s.Logger.Println("Stopping the service")
	return s.Server.Shutdown(ctx)
}

// NewService creates a new Service instance with server and logger
func NewService(server *http.Server, logger *log.Logger) RESTServer {
	return &Service{
		Server: server,
		Logger: logger,
	}
}
