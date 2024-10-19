package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Monkman08/megaport-prometheus-exporter/pkg/handlers"
	"github.com/Monkman08/megaport-prometheus-exporter/pkg/megaport"
	"github.com/Monkman08/megaport-prometheus-exporter/pkg/metrics"
)

func main() {
	// Initialize custom metrics
	metrics.RegisterMetrics()

	// Initialize Megaport client
	megaportClient, err := megaport.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Megaport client: %v", err)
	}

	// Setup HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.LandingPageHandler)
	mux.HandleFunc("/healthz", handlers.HealthzHandler)
	mux.HandleFunc("/readiness", handlers.ReadinessHandler)
	mux.Handle("/metrics", metrics.MetricsHandler(megaportClient))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Run the server in a goroutine
	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :8080: %v\n", err)
		}
	}()

	// Wait for termination signal
	<-stop
	log.Println("Shutting down server...")

	// Create a deadline to wait for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
