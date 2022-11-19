package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// RESTServer wraps the http server functionality for startup and graceful shutdown
type RESTServer interface {
	Start() error
	Shutdown(context.Context) error
}

func (s *Service) Start() error {
	s.Logger.Info("Starting the service")
	return s.Server.ListenAndServe()
}

func (s *Service) Shutdown(ctx context.Context) error {
	s.Logger.Info("Stopping the service")
	return s.Server.Shutdown(ctx)
}

// Service top level struct containing any injected dependencies
type Service struct {
	Server *http.Server
	Logger *zap.Logger
}

func main() {
	service, err := InitialiseService()
	if err != nil {
		panic(err)
	}

	// start REST server in separate go routine to allow for execution to continue
	go func() {
		if err := service.Start(); err != nil {
			service.Logger.Error("failed to start rest server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal)	// channel is used to register a system signall used to stop the service
	signal.Notify(quit, syscall.SIGINT)

	// wait for a Ctrl+C signal before shutting down server
	<-quit
	service.Logger.Info("cancel signal registered")
	if err := service.Shutdown(context.Background()); err != nil {
		service.Logger.Error("gracefully shutting down server", zap.Error(err))
	}
}
