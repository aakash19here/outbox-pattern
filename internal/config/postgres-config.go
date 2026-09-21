package config

import (
	"fmt"
	"os"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func GetDefaultConnString() string {
	return GetConnString(getEnv("POSTGRES_DB", "events"))
}

func GetConnString(dbname string) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("POSTGRES_HOST", "localhost"),
		getEnv("POSTGRES_PORT", "5432"),
		getEnv("POSTGRES_USER", "user"),
		getEnv("POSTGRES_PASSWORD", "pass"),
		dbname,
	)
}
