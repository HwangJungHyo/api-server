package config

import "os"

// Config holds the application configuration.
type Config struct {
	Port      string
	JWTSecret string
}

// Load reads configuration from environment variables with defaults.
func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-in-production"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
