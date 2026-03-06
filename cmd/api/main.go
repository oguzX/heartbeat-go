package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oguzx/devpulse/internal/config"
	"github.com/oguzx/devpulse/internal/db"
	"github.com/oguzx/devpulse/internal/http/routes"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		slog.Error("Failed to load config", "error", err)
	}

	ctx := context.Background()

	dbPool, err := db.NewPool(ctx, cfg)

	if err != nil {
		slog.Error("Failed to connect to db", "error", err)
	}

	defer dbPool.Close()

	router := routes.NewRouter(dbPool)

	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		slog.Info("Server has started at", "port", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	stopCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-stopCtx.Done()
	slog.Info("shutdown signal recieved")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("server stoped")
}
