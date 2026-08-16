package config

import "os"

// Config holds the application configuration.
type Config struct {
	Port      string
	JWTSecret string
	EnableTestEndpoints bool
}

// Load reads configuration from environment variables with defaults.
func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		EnableTestEndpoints: getEnv("ENABLE_TEST_ENDPOINTS", "false") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
