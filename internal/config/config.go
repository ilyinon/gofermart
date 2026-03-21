package config

import (
	"crypto/rand"
	"encoding/hex"
	"os"
)

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	JWTSecret            string
}

func Load() *Config {

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = generateSecret(32)
	}

	return &Config{
		RunAddress:           getEnv("RUN_ADDRESS", ":8080"),
		DatabaseURI:          getEnv("DATABASE_URI", ""),
		AccrualSystemAddress: getEnv("ACCRUAL_SYSTEM_ADDRESS", ""),
		JWTSecret:            secret,
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func generateSecret(n int) string {
	b := make([]byte, n)

	if _, err := rand.Read(b); err != nil {
		panic("failed to generate JWT secret")
	}

	return hex.EncodeToString(b)
}
