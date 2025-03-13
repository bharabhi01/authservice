package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Env string
	ReadTimeout time.Duration
	WriteTimeout time.Duration
	ShutdownTimeout time.Duration

	JWTSecret string
	JWTExpirationHours int
	RefreshExpirationDays int

	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string
	DBSSLMode string

	RedisAddr string
	RedisPassword string
	RedisDB int
	
	CORSAllowOrigins []string	
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using env variables")
	}

	config := &Config{
		Port: getEnv("PORT", "8080"),
		Env:                  getEnv("ENV", "development"),
		ReadTimeout:          time.Duration(getEnvAsInt("READ_TIMEOUT", 5)) * time.Second,
		WriteTimeout:         time.Duration(getEnvAsInt("WRITE_TIMEOUT", 10)) * time.Second,
		ShutdownTimeout:      time.Duration(getEnvAsInt("SHUTDOWN_TIMEOUT", 5)) * time.Second,
		JWTSecret:            getEnv("JWT_SECRET", "your_jwt_secret_key_here"),
		JWTExpirationHours:   getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		RefreshExpirationDays: getEnvAsInt("REFRESH_EXPIRATION_DAYS", 7),
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", "postgres"),
		DBName:               getEnv("DB_NAME", "authservice"),
		DBSSLMode:            getEnv("DB_SSLMODE", "disable"),
		RedisAddr:            getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:        getEnv("REDIS_PASSWORD", ""),
		RedisDB:              getEnvAsInt("REDIS_DB", 0),
		CORSAllowOrigins:     getEnvAsSlice("CORS_ALLOW_ORIGINS", []string{"*"}),
	}

	// Validate required configuration
	if config.JWTSecret == "your_jwt_secret_key_here" && config.Env == "production" {
		return nil, fmt.Errorf("JWT_SECRET must be set in production environment")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	return []string{valueStr}
}
