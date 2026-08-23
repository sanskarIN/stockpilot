package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment   string
	HTTPAddr      string
	DatabaseURL   string
	SessionSecret string
	SessionTTL    time.Duration
	CORSOrigins   []string
	LogLevel      string
	BackupDir     string
	AutoMigrate   bool
	MigrationsDir string
	StaticDir     string
}

func Load() (Config, error) {
	cfg := Config{
		Environment:   env("APP_ENV", "development"),
		HTTPAddr:      env("HTTP_ADDR", ":8080"),
		DatabaseURL:   env("DATABASE_URL", "postgres://stockpilot:stockpilot@localhost:5432/stockpilot?sslmode=disable"),
		SessionSecret: os.Getenv("SESSION_SECRET"),
		LogLevel:      env("LOG_LEVEL", "info"),
		BackupDir:     env("BACKUP_DIR", "./backups"),
		MigrationsDir: env("MIGRATIONS_DIR", "./migrations"),
		StaticDir:     env("STATIC_DIR", "./web/dist"),
	}

	var err error
	cfg.SessionTTL, err = time.ParseDuration(env("SESSION_TTL", "12h"))
	if err != nil || cfg.SessionTTL <= 0 {
		return Config{}, fmt.Errorf("SESSION_TTL must be a positive duration")
	}

	cfg.AutoMigrate, err = strconv.ParseBool(env("AUTO_MIGRATE", "false"))
	if err != nil {
		return Config{}, fmt.Errorf("AUTO_MIGRATE must be true or false")
	}

	for _, origin := range strings.Split(env("CORS_ORIGINS", "http://localhost:5173"), ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			cfg.CORSOrigins = append(cfg.CORSOrigins, origin)
		}
	}

	if strings.TrimSpace(cfg.HTTPAddr) == "" {
		return Config{}, fmt.Errorf("HTTP_ADDR cannot be empty")
	}
	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		return Config{}, fmt.Errorf("DATABASE_URL cannot be empty")
	}
	if cfg.Environment == "production" && len(cfg.SessionSecret) < 32 {
		return Config{}, fmt.Errorf("SESSION_SECRET must contain at least 32 characters in production")
	}
	return cfg, nil
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
