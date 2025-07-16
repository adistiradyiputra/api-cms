package config

import (
	"os"
)

type OPAConfig struct {
	APIBaseURL string
}

var OPA OPAConfig

func LoadOPAConfig() {
	OPA = OPAConfig{
		APIBaseURL: getEnv("OPA_API_BASE_URL", "https://api.opa.co.id"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}