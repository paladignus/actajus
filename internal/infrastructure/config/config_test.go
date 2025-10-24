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
		"ACCESS_SECRET", "REFRESH_SECRET", "JWT_ISSUER", "ACCESS_EXPIRE", "REFRESH_EXPIRE",
	}
	for _, envVar := range envVars {
		originalEnv[envVar] = os.Getenv(envVar)
		err := os.Unsetenv(envVar)
		assert.NoError(t, err)
	}
	defer func() {
		for key, value := range originalEnv {
			if value != "" {
				err := os.Setenv(key, value)
				assert.NoError(t, err)
			} else {
				err := os.Unsetenv(key)
				assert.NoError(t, err)
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
					Port:     "5432",
					User:     "postgres",
					Password: "M4rc3l0",
					DBName:   "actajus",
					SSLMode:  "disable",
				},
				JWT: JWTConfig{
					AccessSecret:  "your-secret-key-change-in-production",
					RefreshSecret: "your-secret-key-change-in-production",
					Issuer:        "prod.example.com.br",
					AccessExpire:  15 * time.Minute,
					RefreshExpire: 168 * time.Hour, // 7 * 24
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
				"JWT_ISSUER":     "example.com.br",
				"ACCESS_EXPIRE":  "30",
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
					Port:     "5433",
					User:     "myuser",
					Password: "mypassword",
					DBName:   "mydatabase",
					SSLMode:  "require",
				},
				JWT: JWTConfig{
					AccessSecret:  "custom-access-secret",
					RefreshSecret: "custom-refresh-secret",
					Issuer:        "example.com.br",
					AccessExpire:  30 * time.Minute,
					RefreshExpire: 720 * time.Hour,
				},
			},
		},
		{
			name: "should partial custom values",
			envVars: map[string]string{
				"SERVER_PORT":   "3000",
				"DB_HOST":       "127.0.0.1",
				"DB_USER":       "testuser",
				"ACCESS_EXPIRE": "5",
			},
			expected: Config{
				Server: ServerConfig{
					Port:     "3000",
					GRPCPort: "50051",
					Env:      "development",
				},
				Database: DatabaseConfig{
					Host:     "127.0.0.1",
					Port:     "5432",
					User:     "testuser",
					Password: "M4rc3l0",
					DBName:   "actajus",
					SSLMode:  "disable",
				},
				JWT: JWTConfig{
					AccessSecret:  "your-secret-key-change-in-production",
					RefreshSecret: "your-secret-key-change-in-production",
					Issuer:        "prod.example.com.br",
					AccessExpire:  5 * time.Minute,
					RefreshExpire: 168 * time.Hour, // default
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				err := os.Setenv(key, value)
				assert.NoError(t, err)
			}
			config := Load()
			assert.Equal(t, tt.expected, config)
			for key := range tt.envVars {
				err := os.Unsetenv(key)
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "should return value from environment variable set",
			envValue:     "custom-value",
			defaultValue: "default-value",
			expected:     "custom-value",
		},
		{
			name:         "should return default value when environment variable not set",
			envValue:     "",
			defaultValue: "default-value",
			expected:     "default-value",
		},
		{
			name:         "should return default value when environment variable empty",
			envValue:     "",
			defaultValue: "fallback",
			expected:     "fallback",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalValue := os.Getenv("TEST_ENV_VAR")
			err := os.Setenv("TEST_ENV_VAR", originalValue)
			assert.NoError(t, err)
			if tt.envValue != "" {
				err := os.Setenv("TEST_ENV_VAR", tt.envValue)
				assert.NoError(t, err)
			} else {
				err := os.Unsetenv("TEST_ENV_VAR")
				assert.NoError(t, err)
			}
			result := getEnv("TEST_ENV_VAR", tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnvAsInt(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue int
		expected     int
	}{
		{
			name:         "should return a valid integer when valid integer in environment",
			envValue:     "42",
			defaultValue: 100,
			expected:     42,
		},
		{
			name:         "should return default value when environment variable not set",
			envValue:     "",
			defaultValue: 100,
			expected:     100,
		},
		{
			name:         "should return default value when invalid integer in environment",
			envValue:     "not-a-number",
			defaultValue: 100,
			expected:     100,
		},
		{
			name:         "should return a valid integer when negative integer",
			envValue:     "-10",
			defaultValue: 100,
			expected:     -10,
		},
		{
			name:         "should return a valid integer when zero value",
			envValue:     "0",
			defaultValue: 100,
			expected:     0,
		},
		{
			name:         "should return a valid integer when large integer",
			envValue:     "999999",
			defaultValue: 100,
			expected:     999999,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalValue := os.Getenv("TEST_INT_VAR")
			err := os.Setenv("TEST_INT_VAR", originalValue)
			assert.NoError(t, err)
			if tt.envValue != "" {
				err := os.Setenv("TEST_INT_VAR", tt.envValue)
				assert.NoError(t, err)
			} else {
				err := os.Unsetenv("TEST_INT_VAR")
				assert.NoError(t, err)
			}
			result := getEnvAsInt("TEST_INT_VAR", tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}
