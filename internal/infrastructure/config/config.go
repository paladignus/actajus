// Package config
package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port     string
	GRPCPort string
	Env      string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	Issuer        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

func Load() Config {
	return Config{
		Server: ServerConfig{
			Port:     getEnv("SERVER_PORT", "8080"),
			GRPCPort: getEnv("GRPC_PORT", "50051"),
			Env:      getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "M4rc3l0"),
			DBName:   getEnv("DB_NAME", "actajus"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			AccessSecret:  getEnv("ACCESS_SECRET", "your-secret-key-change-in-production"),
			RefreshSecret: getEnv("REFRESH_SECRET", "your-secret-key-change-in-production"),
			Issuer:        getEnv("JWT_ISSUER", "prod.example.com.br"),
			AccessExpiry:  time.Duration(getEnvAsInt("ACCESS_EXPIRE", 15)) * time.Minute,
			RefreshExpiry: time.Duration(getEnvAsInt("REFRESH_EXPIRE", 7*24)) * time.Hour,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	// if value, exists := syscall.Getenv(key); exists {
	// 	return value
	// }
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
