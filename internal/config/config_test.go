package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Unset env vars to test defaults
	os.Unsetenv("PORT")
	os.Unsetenv("JWT_SECRET")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Port)
	}
	if cfg.JWTSecret != "dev-secret-change-in-production" {
		t.Errorf("expected default JWT secret, got %s", cfg.JWTSecret)
	}
}

func TestLoad_FromEnv(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("JWT_SECRET", "my-secret")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg := Load()

	if cfg.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Port)
	}
	if cfg.JWTSecret != "my-secret" {
		t.Errorf("expected JWT secret my-secret, got %s", cfg.JWTSecret)
	}
}
