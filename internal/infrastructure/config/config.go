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
	Tracing  TracingConfig
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
	Subjects          []string
	URL               string
	StreamName        string
	ConsumerName      string
	DurableName       string
	ReplayPolicy      string
	MaxBytes          int64
	Replicas          int
	MaxDeliver        int
	MaxAckPending     int
	MaxAge            time.Duration
	AckWait           time.Duration
	ConnectionTimeout time.Duration
	RequestTimeout    time.Duration
}

type TracingConfig struct {
	ServiceName string
	AgentHost   string
	AgentPort   string
	Enabled     bool
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
			Subjects:          strings.Split(getEnv("NATS_SUBJECTS", "events.>"), ","),
			URL:               getEnv("NATS_URL", "nats://localhost:4222"),
			StreamName:        getEnv("NATS_STREAM_NAME", "EVENTS"),
			ConsumerName:      getEnv("NATS_CONSUMER_NAME", "actajus-consumer"),
			DurableName:       getEnv("NATS_DURABLE_NAME", "actajus-durable"),
			ReplayPolicy:      getEnv("NATS_REPLAY_POLICY", "last"),
			MaxBytes:          int64(getEnvAsInt("NATS_MAX_BYTES", 1024*1024*1024)),
			Replicas:          getEnvAsInt("NATS_REPLICAS", 1),
			MaxDeliver:        getEnvAsInt("NATS_MAX_DELIVER", 3),
			MaxAckPending:     getEnvAsInt("NATS_MAX_ACK_PENDING", 100),
			MaxAge:            time.Duration(getEnvAsInt("NATS_MAX_AGE", 7*24)) * time.Hour,
			AckWait:           time.Duration(getEnvAsInt("NATS_ACK_WAIT", 30)) * time.Second,
			ConnectionTimeout: time.Duration(getEnvAsInt("NATS_CONNECTION_TIMEOUT", 10)) * time.Second,
			RequestTimeout:    time.Duration(getEnvAsInt("NATS_REQUEST_TIMEOUT", 5)) * time.Second,
		},
		Tracing: TracingConfig{
			ServiceName: getEnv("TRACING_SERVICE_NAME", "actajus-api"),
			AgentHost:   getEnv("TRACING_AGENT_HOST", "localhost"),
			AgentPort:   getEnv("TRACING_AGENT_PORT", "6831"),
			Enabled:     getEnvAsBool("TRACING_ENABLED", true), // Default to true for development
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
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

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
