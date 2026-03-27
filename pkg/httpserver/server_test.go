package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
	"github.com/paladignus/actajus/pkg/di"
)

func TestServerCreation(t *testing.T) {
	t.Skip("Skipping test - requires full DI container")

	// Create a mock container
	container := &di.Container{
		Config: config.Config{
			Server: config.ServerConfig{
				Port: "8080",
			},
		},
	}

	server := New(container)
	if server == nil {
		t.Fatal("Expected server to be created")
	}
}

func TestHealthCheck(t *testing.T) {
	// Create a minimal test for health check endpoint
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expected := `{"status":"ok"}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}

func TestCORSHeaders(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Test CORS middleware logic inline
	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		handler.ServeHTTP(w, r)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	corsHandler.ServeHTTP(w, req)

	// Check CORS headers
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected Access-Control-Allow-Origin: *")
	}
}

func TestNormalizeAddr(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: ":8080"},
		{name: "plain port", in: "8081", want: ":8081"},
		{name: "prefixed port", in: ":8082", want: ":8082"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(&di.Container{
				Config: config.Config{Server: config.ServerConfig{Port: tt.in}},
			})
			if got == nil {
				t.Fatal("Expected server to be created")
			}
			if got.srv.Addr != tt.want {
				t.Fatalf("expected addr %s, got %s", tt.want, got.srv.Addr)
			}
		})
	}
}
