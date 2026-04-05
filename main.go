package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"weatherservice/api"
	weatherclient "weatherservice/weatherClient"
)

func main() {
	addr := ":8080"

	// Initialize weather client and server
	weatherClient := weatherclient.New()
	server := api.NewServer(addr, weatherClient)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to bind %s (is another process using this port?): %v", addr, err)
	}

	log.Printf("weather service listening on %s", addr)

	// Start goroutine for serving traffic
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// Handle graceful shutdown to avoid hanging port binding
	shutdownSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-shutdownSignal.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}

	log.Println("server stopped cleanly")
}
