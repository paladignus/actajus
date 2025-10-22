// Package config
package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	originalEnv := make(map[string]string)
	envVars := []string{
		"SERVER_PORT", "GRPC_PORT", "ENV",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE",
		"ACCESS_SECRET", "REFRESH_SECRET", "ACCESS_EXPIRY", "REFRESH_EXPIRE",
	}
	for _, envVar := range envVars {
		originalEnv[envVar] = os.Getenv(envVar)
		os.Unsetenv(envVar)
	}
	defer func() {
		for key, value := range originalEnv {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
	}()
	tests := []struct {
		name     string
		envVars  map[string]string
		expected Config
	}{
		{
			name:    "should default values",
			envVars: map[string]string{},
			expected: Config{
				Server: ServerConfig{
					Port:     "8080",
					GRPCPort: "50051",
					Env:      "development",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "postgres",
					Password: "M4rc3l0",
					DBName:   "actajus",
					SSLMode:  "disable",
				},
				JWT: JWTConfig{
					AccessSecret:  "your-secret-key-change-in-production",
					RefreshSecret: "your-secret-key-change-in-production",
					AccessExpiry:  15 * time.Minute,
					RefreshExpiry: 168 * time.Hour, // 7 * 24
				},
			},
		},
		{
			name: "should custom values from environment",
			envVars: map[string]string{
				"SERVER_PORT":    "9090",
				"GRPC_PORT":      "50052",
				"ENV":            "production",
				"DB_HOST":        "db.example.com",
				"DB_PORT":        "5433",
				"DB_USER":        "myuser",
				"DB_PASSWORD":    "mypassword",
				"DB_NAME":        "mydatabase",
				"DB_SSLMODE":     "require",
				"ACCESS_SECRET":  "custom-access-secret",
				"REFRESH_SECRET": "custom-refresh-secret",
				"ACCESS_EXPIRY":  "30",
				"REFRESH_EXPIRE": "720", // 30 days in hours
			},
			expected: Config{
				Server: ServerConfig{
					Port:     "9090",
					GRPCPort: "50052",
					Env:      "production",
				},
				Database: DatabaseConfig{
					Host:     "db.example.com",
					Port:     5433,
					User:     "myuser",
					Password: "mypassword",
					DBName:   "mydatabase",
					SSLMode:  "require",
				},
				JWT: JWTConfig{
					AccessSecret:  "custom-access-secret",
					RefreshSecret: "custom-refresh-secret",
					AccessExpiry:  30 * time.Minute,
					RefreshExpiry: 720 * time.Hour,
				},
			},
		},
		{
			name: "should partial custom values",
			envVars: map[string]string{
				"SERVER_PORT":   "3000",
				"DB_HOST":       "127.0.0.1",
				"DB_USER":       "testuser",
				"ACCESS_EXPIRY": "5",
			},
			expected: Config{
				Server: ServerConfig{
					Port:     "3000",
					GRPCPort: "50051",
					Env:      "development",
				},
				Database: DatabaseConfig{
					Host:     "127.0.0.1",
					Port:     5432,
					User:     "testuser",
					Password: "M4rc3l0",
					DBName:   "actajus",
					SSLMode:  "disable",
				},
				JWT: JWTConfig{
					AccessSecret:  "your-secret-key-change-in-production",
					RefreshSecret: "your-secret-key-change-in-production",
					AccessExpiry:  5 * time.Minute,
					RefreshExpiry: 168 * time.Hour, // default
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			config := Load()
			assert.Equal(t, tt.expected, config)
			for key := range tt.envVars {
				os.Unsetenv(key)
			}
		})
	}
}
