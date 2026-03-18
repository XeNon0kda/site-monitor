package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port          string
	CheckInterval time.Duration
	HTTPTimeout   time.Duration
}

func Load() *Config {
	return &Config{
		Port:          ":" + getEnv("PORT", "8080"),
		CheckInterval: time.Duration(getEnvAsInt("CHECK_INTERVAL_SEC", 60)) * time.Second,
		HTTPTimeout:   time.Duration(getEnvAsInt("HTTP_TIMEOUT_SEC", 5)) * time.Second,
	}
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func getEnvAsInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}