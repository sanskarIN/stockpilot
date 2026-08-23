package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sanskarIN/stockpilot/internal/auth"
	"github.com/sanskarIN/stockpilot/internal/config"
	"github.com/sanskarIN/stockpilot/internal/postgres"
)

func main() {
	if len(os.Args) != 2 || os.Args[1] != "bootstrap" {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/admin bootstrap")
		fmt.Fprintln(os.Stderr, "set STOCKPILOT_ADMIN_EMAIL, STOCKPILOT_ADMIN_NAME, and STOCKPILOT_ADMIN_PASSWORD first")
		os.Exit(2)
	}

	email := strings.TrimSpace(os.Getenv("STOCKPILOT_ADMIN_EMAIL"))
	name := strings.TrimSpace(os.Getenv("STOCKPILOT_ADMIN_NAME"))
	password := os.Getenv("STOCKPILOT_ADMIN_PASSWORD")
	if email == "" || name == "" || password == "" {
		fmt.Fprintln(os.Stderr, "STOCKPILOT_ADMIN_EMAIL, STOCKPILOT_ADMIN_NAME, and STOCKPILOT_ADMIN_PASSWORD are required")
		os.Exit(2)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid configuration: %v\n", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	store, err := postgres.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "database connection failed: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()
	if err := store.Migrate(ctx, cfg.MigrationsDir); err != nil {
		fmt.Fprintf(os.Stderr, "database migration failed: %v\n", err)
		os.Exit(1)
	}

	service := auth.New(store, cfg.SessionTTL, cfg.SessionSecret)
	user, err := service.BootstrapAdmin(ctx, email, name, password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "administrator bootstrap failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("StockPilot administrator created: %s (%s)\n", user.DisplayName, user.Email)
}
