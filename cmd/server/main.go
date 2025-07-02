package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ru/mvideo/com/gitlab/token-exporter/internal/config"
	"ru/mvideo/com/gitlab/token-exporter/internal/gitlab"
	"ru/mvideo/com/gitlab/token-exporter/internal/metrics"
	"ru/mvideo/com/gitlab/token-exporter/internal/scraper"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		cancel()
	}()

	gitlabClient, err := gitlab.NewClient(cfg.Gitlab.Token, cfg.Gitlab.BaseURL)
	if err != nil {
		log.Fatalf("Failed to create GitLab client: %v", err)
	}

	metricsHandler := metrics.NewHandler()
	tokenScraper := scraper.NewTokenScraper(gitlabClient, metricsHandler, []int(cfg.Gitlab.ProjectIDs), []int(cfg.Gitlab.GroupIDs))

	go func() {
		tokenScraper.Start(ctx, cfg.Scraper.Interval)
	}()

	mux := http.NewServeMux()
	mux.Handle("/metrics", metricsHandler.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: mux,
	}

	go func() {
		log.Printf("Starting HTTP server on :%d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Server stopped gracefully")
}
