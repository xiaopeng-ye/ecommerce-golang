package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	DBHost                 string
	DBPort                 string
	JWTSecret              string
	JWTExpirationInSeconds int64
	Environment            string // development, staging, production
	LogLevel               string // debug, info, warn, error
	LogFormat              string // json, text
}

var Envs = initConfig()

func initConfig() Config {
	// Only load .env file in development
	// In production, environment variables should be set by the orchestration platform
	if os.Getenv("ENVIRONMENT") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not found, using environment variables and defaults")
		}
	}

	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")

	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8080"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnvRequired("DB_PASSWORD"),
		DBHost:                 dbHost,
		DBPort:                 dbPort,
		DBAddress:              fmt.Sprintf("%s:%s", dbHost, dbPort),
		DBName:                 getEnv("DB_NAME", "ecom"),
		JWTSecret:              getEnvRequired("JWT_SECRET"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
		Environment:            getEnv("ENVIRONMENT", "development"),
		LogLevel:               getEnv("LOG_LEVEL", "info"),
		LogFormat:              getEnv("LOG_FORMAT", "json"),
	}
}

// getEnv gets the env by key or returns fallback
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvRequired gets a required environment variable or uses test default
func getEnvRequired(key string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	// In test/development environment, provide safe defaults
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "development" || env == "test" {
		log.Printf("Warning: Required environment variable %s is not set, using test default", key)
		switch key {
		case "DB_PASSWORD":
			return "test-password"
		case "JWT_SECRET":
			return "test-secret-key-for-testing-only-do-not-use-in-production"
		}
	}

	log.Fatalf("Required environment variable %s is not set in production", key)
	return ""
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Printf("Warning: invalid value for %s, using fallback %d", key, fallback)
			return fallback
		}
		return i
	}
	return fallback
}
