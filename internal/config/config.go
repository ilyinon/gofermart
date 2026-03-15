package config

import "os"

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
}

func Load() *Config {

	return &Config{
		RunAddress:           getEnv("RUN_ADDRESS", ":8080"),
		DatabaseURI:          getEnv("DATABASE_URI", ""),
		AccrualSystemAddress: getEnv("ACCRUAL_SYSTEM_ADDRESS", ""),
	}
}

func getEnv(key string, fallback string) string {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
