package config

import "testing"

func TestGetEnvReturnsDefaultWhenEmpty(t *testing.T) {
	key := "TEST_GET_ENV_EMPTY"
	t.Setenv(key, "")

	got := getEnv(key, "default")

	if got != "default" {
		t.Fatalf("expected default value, got %q", got)
	}
}

func TestGetEnvReturnsValueWhenSet(t *testing.T) {
	key := "TEST_GET_ENV_VALUE"
	t.Setenv(key, "custom")

	got := getEnv(key, "default")

	if got != "custom" {
		t.Fatalf("expected env value, got %q", got)
	}
}

func TestLoadReturnsErrorWhenDSNMissing(t *testing.T) {
	t.Setenv("POSTGRES_DSN", "")
	t.Setenv("JWT_SECRET", "test-secret")

	_, err := Load()

	if err == nil {
		t.Fatal("expected error when POSTGRES_DSN is missing")
	}
}

func TestLoadUsesDefaultsAndEnvValues(t *testing.T) {
	t.Setenv("POSTGRES_DSN", "postgres://user:pass@localhost:5432/db")
	t.Setenv("HTTP_ADDR", ":9090")
	t.Setenv("JWT_SECRET", "test-secret")

	cfg, err := Load()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.HTTPAddress != ":9090" {
		t.Fatalf("expected HTTP address :9090, got %q", cfg.HTTPAddress)
	}
	if cfg.DSN != "postgres://user:pass@localhost:5432/db" {
		t.Fatalf("expected DSN to match env, got %q", cfg.DSN)
	}
}
