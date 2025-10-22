// Package adapter
package adapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockLogger implements repository.Logger for testing
// type MockLogger struct {
// 	mock.Mock
// }
//
// func (m *MockLogger) Debug(ctx context.Context, msg string, args ...any) {
// 	m.Called(ctx, msg, args)
// }
//
// func (m *MockLogger) Info(ctx context.Context, msg string, args ...any) {
// 	m.Called(ctx, msg, args)
// }
//
// func (m *MockLogger) Warn(ctx context.Context, msg string, args ...any) {
// 	m.Called(ctx, msg, args)
// }
//
// func (m *MockLogger) Error(ctx context.Context, msg string, args ...any) {
// 	m.Called(ctx, msg, args)
// }
//
// func (m *MockLogger) With(args ...any) repository.Logger {
// 	callArgs := m.Called(args)
// 	return callArgs.Get(0).(repository.Logger)
// }
//
// func (m *MockLogger) WithError(err error) repository.Logger {
// 	callArgs := m.Called(err)
// 	return callArgs.Get(0).(repository.Logger)
// }

// Helper to capture log output
// type logCapture struct {
// 	*bytes.Buffer
// }
//
// func newLogCapture() *logCapture {
// 	return &logCapture{Buffer: &bytes.Buffer{}}
// }

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "should create a json format with stdout and level info",
			config: Config{
				Level:     LevelInfo,
				AddSource: false,
				Format:    "json",
				Output:    "stdout",
			},
		},
		{
			name: "should create a text format with stderr and level debug",
			config: Config{
				Level:     LevelDebug,
				AddSource: true,
				Format:    "text",
				Output:    "stderr",
			},
		},
		{
			name: "should create a default format and output and level warn",
			config: Config{
				Level:     LevelWarn,
				AddSource: false,
				Format:    "",
				Output:    "",
			},
		},
		{
			name: "should create a json format with stdout and level error",
			config: Config{
				Level:     LevelError,
				AddSource: false,
				Format:    "json",
				Output:    "stdout",
			},
		},
		{
			name:   "should create a default format and output and level",
			config: Config{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.config)
			assert.NotNil(t, logger)
			assert.IsType(t, &slogAdapter{}, logger)
		})
	}
}
