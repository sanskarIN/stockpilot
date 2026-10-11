package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sanskarIN/stockpilot/internal/auth"
	"github.com/sanskarIN/stockpilot/internal/config"
	"github.com/sanskarIN/stockpilot/internal/httpapi"
	"github.com/sanskarIN/stockpilot/internal/postgres"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	startupCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	store, err := postgres.Open(startupCtx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database startup failed", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	if cfg.AutoMigrate {
		if err := store.Migrate(startupCtx, cfg.MigrationsDir); err != nil {
			logger.Error("database migration failed", "error", err)
			os.Exit(1)
		}
	}

	authService := auth.New(store, cfg.SessionTTL, cfg.SessionSecret)
	coreAPI := httpapi.NewCore(store, store, store, store.Ping)
	protectedAPI := httpapi.WithAccess(coreAPI, authService, store, cfg.Environment == "production")

	root := http.NewServeMux()
	root.Handle("/api/", protectedAPI)
	root.Handle("/healthz", coreAPI)
	root.Handle("/readyz", coreAPI)
	root.Handle("/", http.FileServer(http.Dir(cfg.StaticDir)))

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           httpapi.WrapCommon(root, cfg.CORSOrigins, logger),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-stopCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("graceful shutdown failed", "error", err)
		}
	}()

	logger.Info("StockPilot started", "addr", cfg.HTTPAddr, "environment", cfg.Environment, "credit", "Made by the Sanskar")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("http server failed", "error", err)
		os.Exit(1)
	}
	logger.Info("StockPilot stopped")
}
