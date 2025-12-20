package config

import "os"

type Config struct {
	DatabaseURL string
	JWTSecret   string
	CORSOrigins string
	Port        string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://linkbio:linkbio_dev_password@localhost:5432/linkbio?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:5173"),
		Port:        getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
