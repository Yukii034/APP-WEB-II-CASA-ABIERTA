package config

import "os"

// Puerto returns the configured service port, defaulting to 8083 for local use.
func Puerto() string {
	if puerto := os.Getenv("PORT"); puerto != "" {
		return puerto
	}
	return "8083"
}
