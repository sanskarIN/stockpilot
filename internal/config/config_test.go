package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadUsesDevelopmentDefaults(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("SESSION_SECRET", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if cfg.Environment != "development" {
		t.Fatalf("Environment = %q, want development", cfg.Environment)
	}
	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("HTTPAddr = %q, want :8080", cfg.HTTPAddr)
	}
	if cfg.SessionTTL != 12*time.Hour {
		t.Fatalf("SessionTTL = %v, want 12h", cfg.SessionTTL)
	}
	if cfg.AutoMigrate {
		t.Fatal("AutoMigrate = true, want false by default")
	}
}

func TestLoadRequiresStrongProductionSessionSecret(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("APP_ENV", "production")
	t.Setenv("SESSION_SECRET", "too-short")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "at least 32") {
		t.Fatalf("Load() error = %v, want production session secret validation", err)
	}
}

func TestLoadParsesOriginsAndMigrationFlag(t *testing.T) {
	clearConfigEnv(t)
	t.Setenv("CORS_ORIGINS", "https://inventory.example, https://admin.example ")
	t.Setenv("AUTO_MIGRATE", "true")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}
	if !cfg.AutoMigrate {
		t.Fatal("AutoMigrate = false, want true")
	}
	if got, want := len(cfg.CORSOrigins), 2; got != want {
		t.Fatalf("len(CORSOrigins) = %d, want %d", got, want)
	}
	if cfg.CORSOrigins[1] != "https://admin.example" {
		t.Fatalf("second origin = %q", cfg.CORSOrigins[1])
	}
}

func clearConfigEnv(t *testing.T) {
	t.Helper()
	keys := []string{
		"APP_ENV", "HTTP_ADDR", "DATABASE_URL", "SESSION_SECRET", "SESSION_TTL", "CORS_ORIGINS",
		"LOG_LEVEL", "BACKUP_DIR", "AUTO_MIGRATE", "MIGRATIONS_DIR", "STATIC_DIR",
	}
	for _, key := range keys {
		t.Setenv(key, "")
	}
	// Environment variables intentionally set to an empty string override defaults,
	// so restore the variables that Load expects to have useful defaults.
	t.Setenv("APP_ENV", "development")
	t.Setenv("HTTP_ADDR", ":8080")
	t.Setenv("DATABASE_URL", "postgres://stockpilot:stockpilot@localhost:5432/stockpilot?sslmode=disable")
	t.Setenv("SESSION_TTL", "12h")
	t.Setenv("CORS_ORIGINS", "http://localhost:5173")
	t.Setenv("LOG_LEVEL", "info")
	t.Setenv("BACKUP_DIR", "./backups")
	t.Setenv("AUTO_MIGRATE", "false")
	t.Setenv("MIGRATIONS_DIR", "./migrations")
	t.Setenv("STATIC_DIR", "./web/dist")
}
