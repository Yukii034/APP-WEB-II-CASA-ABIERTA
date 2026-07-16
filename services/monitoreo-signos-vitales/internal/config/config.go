package config

import "os"

// Puerto returns the configured service port, defaulting to 8080 inside Docker.
func Puerto() string {
	if puerto := os.Getenv("PORT"); puerto != "" {
		return puerto
	}
	return "8080"
}
