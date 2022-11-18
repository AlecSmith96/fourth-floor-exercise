package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"github.com/AlecSmith96/fourth-floor-exercise/adapters"
	"github.com/AlecSmith96/fourth-floor-exercise/rest"
)

func main() {
	logger := log.New(log.Writer(), "", 0)
	client := &http.Client{}
	adapter := adapters.NewTwitchAdapter(client, logger)
	r := rest.NewRouter(adapter)
	server := rest.NewHTTPServer(r)
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
