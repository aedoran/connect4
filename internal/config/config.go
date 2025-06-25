package config

import (
	"os"
)

// Config holds application configuration values.
type Config struct {
	// HTTPPort is the port the API server listens on.
	HTTPPort string
}

// Load reads configuration from environment variables or defaults.
func Load() Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	return Config{HTTPPort: port}
}
