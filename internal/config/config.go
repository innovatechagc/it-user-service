package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	Port              string
	LogLevel          string
	Environment       string
	RateLimitRPS      int
	RateLimitBurst    int
	AuthServiceURL    string // URL del servicio de autenticación
}

func LoadConfig() Config {
	return Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "itapp"),
		Port:           getEnv("PORT", "8083"), // Puerto diferente para user service
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		Environment:    getEnv("ENVIRONMENT", "development"),
		RateLimitRPS:   getEnvAsInt("RATE_LIMIT_RPS", 100),
		RateLimitBurst: getEnvAsInt("RATE_LIMIT_BURST", 200),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8082"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}