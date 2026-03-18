package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPPort    int
	DBURL       string
	RedisURL    string
	JWTSecret   string
	Environment string
	MaxDevices  int
}

func Load() *Config {
	port, _ := strconv.Atoi(getEnv("HTTP_PORT", "8080"))
	maxDevices, _ := strconv.Atoi(getEnv("MAX_DEVICES_PER_USER", "5"))
	return &Config{
		HTTPPort:    port,
		DBURL:       getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/vpndb?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production"),
		Environment: getEnv("ENV", "development"),
		MaxDevices:  maxDevices,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
