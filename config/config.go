package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	PORT                string
	DB_URL              string
	REDIS_URL           string
	JWT_SECRET          string
	JWT_EXPIRATION_HOURS int
	RATE_LIMIT          int
	LOG_LEVEL           string
}

func LoadConfig() (*Config, error) {
	config := &Config{}

	config.PORT = getEnv("PORT", "8080")
	config.JWT_SECRET = getEnv("JWT_SECRET", "")
	if config.JWT_SECRET == "" {
		return nil, fmt.Errorf("JWT_SECRET missing")
	}

	expHoursStr := getEnv("JWT_EXPIRATION_HOURS", "72")
	expHours, err := strconv.Atoi(expHoursStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION_HOURS: %v", err)
	}
	config.JWT_EXPIRATION_HOURS = expHours

	rateLimitStr := getEnv("RATE_LIMIT", "100")
	rateLimit, err := strconv.Atoi(rateLimitStr)
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT: %v", err)
	}
	config.RATE_LIMIT = rateLimit

	config.LOG_LEVEL = strings.ToLower(getEnv("LOG_LEVEL", "info"))
	level := config.LOG_LEVEL
	if level != "info"  && level != "debug" && level != "error"{
		return nil, fmt.Errorf("invalid LOG_LEVEL")
	}

	// Build DB URL
	dbHost := getEnv("DB_HOST", "")
	dbPort := getEnv("DB_PORT", "")
	dbUser := getEnv("DB_USER", "")
	dbPass := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "")
	if dbHost == "" || dbPort == "" || dbUser == "" || dbPass == "" || dbName == "" {
		return nil, fmt.Errorf("database env vars missing")
	}
	config.DB_URL = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	// Build Redis URL
	redisHost := getEnv("REDIS_HOST", "")
	redisPort := getEnv("REDIS_PORT", "")
	redisPass := getEnv("REDIS_PASSWORD", "")
	if redisHost == "" || redisPort == "" {
		return nil, fmt.Errorf("redis env vars missing")
	}
	if redisPass == "" {
		config.REDIS_URL = fmt.Sprintf("redis://%s:%s/0", redisHost, redisPort)
	} else {
		config.REDIS_URL = fmt.Sprintf("redis://:%s@%s:%s/0", redisPass, redisHost, redisPort)
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
