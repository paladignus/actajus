// Package config
package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	Redis         RedisConfig
	JWT           JWTConfig
	Session       SessionConfig
	PasswordReset PasswordResetConfig
	Auth          AuthConfig
	SMTP          SMTPConfig
	NATS          NATSConfig
	Tracing       TracingConfig
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

type RedisConfig struct {
	Addr     string
	Password string
	Prefix   string
	Fallback string
}

type JWTConfig struct {
	AccessSecret string
	Issuer       string
	Audience     string
	AccessTTL    time.Duration
}

type SessionConfig struct {
	MaxSessions int
	RefreshTTL  time.Duration
}

type PasswordResetConfig struct {
	MaxResetAttempts int
	ResetTTL         time.Duration
}

type AuthConfig struct {
	JWTConfig
	SessionConfig
	PasswordResetConfig
}

type SMTPConfig struct {
	Host          string
	Port          string
	User          string
	Pass          string
	From          string
	PublicBaseURL string
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
	// Carregar variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}
	// Verificar se as variáveis obrigatórias estão definidas (apenas fora dos testes)
	isTest := strings.HasSuffix(os.Args[0], ".test")
	// Determinar se estamos em modo de produção (não desenvolvimento nem teste)
	isProduction := !isTest && getEnv("ENV", "development") != "development"
	// Obter os valores das variáveis de ambiente
	password := getEnv("DB_PASSWORD", "") // Sem valor padrão sensível
	accessSecret := getEnv("ACCESS_SECRET", "")
	refreshSecret := getEnv("REFRESH_SECRET", "")
	resetSecret := getEnv("RESET_SECRET", "")
	// Verificar se as variáveis obrigatórias estão definidas em produção
	if isProduction {
		if password == "" {
			log.Fatal("DB_PASSWORD environment variable must be set in production")
		}
		if accessSecret == "" {
			log.Fatal("ACCESS_SECRET environment variable must be set in production")
		}
		if refreshSecret == "" {
			log.Fatal("REFRESH_SECRET environment variable must be set in production")
		}
		if resetSecret == "" {
			log.Fatal("RESET_SECRET environment variable must be set in production")
		}
	}
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
			Password: password,                   // Usando a variável verificada acima
			DBName:   getEnv("DB_NAME", "sidof"), // Nome do banco de dados padrão
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			Prefix:   getEnv("REDIS_PREFIX", "actajus:"),
			Fallback: getEnv("REDIS_FALLBACK", "true"),
		},
		JWT: JWTConfig{
			AccessSecret: accessSecret,                           // Usando a variável verificada acima
			Issuer:       getEnv("JWT_ISSUER", "example.com.br"), // Exemplo de valor
			AccessTTL:    time.Duration(getEnvAsInt("ACCESS_TTL", 15)) * time.Minute,
			Audience:     getEnv("JWT_AUDIENCE", "example.com.br"),
		},
		Session: SessionConfig{
			MaxSessions: getEnvAsInt("MAX_SESSIONS", 10),
			RefreshTTL:  time.Duration(getEnvAsInt("REFRESH_TTL", 24)) * time.Hour,
		},
		PasswordReset: PasswordResetConfig{
			MaxResetAttempts: getEnvAsInt("MAX_RESET_ATTEMPTS", 5),
			ResetTTL:         time.Duration(getEnvAsInt("RESET_TTL", 30)) * time.Minute,
		},
		Auth: AuthConfig{
			JWTConfig{
				AccessSecret: accessSecret,                           // Usando a variável verificada acima
				Issuer:       getEnv("JWT_ISSUER", "example.com.br"), // Exemplo de valor
				AccessTTL:    time.Duration(getEnvAsInt("ACCESS_TTL", 15)) * time.Minute,
				Audience:     getEnv("JWT_AUDIENCE", "example.com.br"),
			},
			SessionConfig{
				MaxSessions: getEnvAsInt("MAX_SESSIONS", 10),
				RefreshTTL:  time.Duration(getEnvAsInt("REFRESH_TTL", 24)) * time.Hour,
			},
			PasswordResetConfig{
				MaxResetAttempts: getEnvAsInt("MAX_RESET_ATTEMPTS", 5),
				ResetTTL:         time.Duration(getEnvAsInt("RESET_TTL", 30)) * time.Minute,
			},
		},
		SMTP: SMTPConfig{
			Host:          getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:          getEnv("SMTP_PORT", "587"),
			User:          getEnv("SMTP_USER", ""), // Sem valor padrão sensível
			Pass:          getEnv("SMTP_PASS", ""), // Sem valor padrão sensível
			From:          getEnv("SMTP_FROM", ""), // Sem valor padrão sensível
			PublicBaseURL: getEnv("PUBLIC_BASE_URL", "http://localhost:5173"),
		},
		NATS: NATSConfig{
			Subjects:          strings.Split(getEnv("NATS_SUBJECTS", "events.>"), ","),
			URL:               getEnv("NATS_URL", "nats://localhost:4222"),
			StreamName:        getEnv("NATS_STREAM_NAME", "EVENTS"),
			ConsumerName:      getEnv("NATS_CONSUMER_NAME", "sidof-consumer"),
			DurableName:       getEnv("NATS_DURABLE_NAME", "sidof-durable"),
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
			ServiceName: getEnv("TRACING_SERVICE_NAME", "sidof-api"),
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
