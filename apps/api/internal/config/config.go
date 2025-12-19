package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://linkbio:linkbio_dev_password@localhost:5432/linkbio?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Default DATABASE_URL for docker: postgres://linkbio:linkbio_dev_password@localhost:5432/linkbio?sslmode=disable
