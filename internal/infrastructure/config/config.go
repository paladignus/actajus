// Package config
package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	SMTP     SMTPConfig
	NATS     NATSConfig
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
	ResetSecret   string
	Issuer        string
	AccessExpire  time.Duration
	RefreshExpire time.Duration
	ResetExpire   time.Duration
}

type SMTPConfig struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

type NATSConfig struct {
	URL           string
	StreamName    string
	Subjects      []string
	DLQSubject    string
	MaxReconnects int
	ReconnectWait time.Duration
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
			ResetSecret:   getEnv("RESET_SECRET", "your-secret-key-change-in-production"),
			Issuer:        getEnv("JWT_ISSUER", "prod.example.com.br"),
			AccessExpire:  time.Duration(getEnvAsInt("ACCESS_EXPIRE", 15)) * time.Minute,
			RefreshExpire: time.Duration(getEnvAsInt("REFRESH_EXPIRE", 7*24)) * time.Hour,
			ResetExpire:   time.Duration(getEnvAsInt("RESET_EXPIRE", 30)) * time.Minute,
		},
		SMTP: SMTPConfig{
			Host: getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port: getEnv("SMTP_PORT", "587"),
			// User: getEnv("SMTP_USER", "your-email"),
			User: getEnv("SMTP_USER", "tidofsejuspms@gmail.com"),
			// Pass: getEnv("SMTP_PASS", "your-password"),
			Pass: getEnv("SMTP_PASS", "vrea ubpt gwxs oirw"),
			// From: getEnv("SMTP_FROM", "ActaJus <your-email>"),
			From: getEnv("SMTP_FROM", "tidofsejuspms@gmail.com"),
		},
		NATS: NATSConfig{
			URL:           getEnv("NATS_URL", "nats://localhost:4222"),
			StreamName:    getEnv("NATS_STREAM_NAME", "EVENTS"),
			Subjects:      strings.Split(getEnv("NATS_SUBJECTS", "auth.>,user.>"), ","),
			MaxReconnects: getEnvAsInt("NATS_MAX_RECONNECTS", 5),
			ReconnectWait: time.Duration(getEnvAsInt("NATS_RECONNECT_WAIT", 2)) * time.Second,
			DLQSubject:    getEnv("NATS_DLQ_SUBJECT", "events.dlq"),
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
