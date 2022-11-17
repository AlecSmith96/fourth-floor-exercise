package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.New(log.Writer(), "", 0)
	adapter := NewTwitchAdapter(logger)
	r := NewRouter(adapter)
	server := NewHTTPServer(r)
	service := NewService(server, logger)

	go func() {
		if err := service.Start(); err != nil {
			logger.Fatalf("failed to start rest server: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT)

	<-quit
	logger.Println("cancel signal registered")
	if err := service.Shutdown(context.Background()); err != nil {
		logger.Fatalf("gracefully shutting down server: %v", err)
	}
}
