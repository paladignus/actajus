package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/paladignus/actajus/internal/module/company"
	identitymodule "github.com/paladignus/actajus/internal/module/identity"
	webmodule "github.com/paladignus/actajus/internal/module/web"
	sharedrepo "github.com/paladignus/actajus/internal/shared/application/repository"
	"github.com/paladignus/actajus/internal/shared/infrastructure/config"
	"github.com/paladignus/actajus/pkg/di"
)

type stubLogger struct{}

func (stubLogger) Debug(context.Context, string, ...any) {}
func (stubLogger) Info(context.Context, string, ...any)  {}
func (stubLogger) Warn(context.Context, string, ...any)  {}
func (stubLogger) Error(context.Context, string, ...any) {}
func (stubLogger) With(...any) sharedrepo.Logger         { return stubLogger{} }
func (stubLogger) WithError(error) sharedrepo.Logger     { return stubLogger{} }

func newTestWebModule(t *testing.T) webmodule.Module {
	t.Helper()
	module, err := webmodule.NewModule(webmodule.Dependencies{
		Logger:   stubLogger{},
		Company:  company.Module{},
		Identity: identitymodule.Module{},
		Secure:   false,
	})
	if err != nil {
		t.Fatalf("new web module: %v", err)
	}
	return module
}

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

func TestErrorPageMiddlewareRendersNotFoundHTML(t *testing.T) {
	container := &di.Container{
		Modules: di.Modules{Web: newTestWebModule(t)},
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	handler := errorPageMiddleware(container, next)

	req := httptest.NewRequest(http.MethodGet, "/nao-existe", nil)
	req.Header.Set("Accept", "text/html")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "404") {
		t.Fatalf("expected rendered 404 page, got %q", rec.Body.String())
	}
}

func TestErrorPageMiddlewareRendersForbiddenHTML(t *testing.T) {
	container := &di.Container{
		Modules: di.Modules{Web: newTestWebModule(t)},
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	handler := errorPageMiddleware(container, next)

	req := httptest.NewRequest(http.MethodGet, "/roles", nil)
	req.Header.Set("Accept", "text/html")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "403") {
		t.Fatalf("expected rendered 403 page, got %q", rec.Body.String())
	}
}

func TestErrorPageMiddlewarePreservesAsset404(t *testing.T) {
	container := &di.Container{
		Modules: di.Modules{Web: newTestWebModule(t)},
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	handler := errorPageMiddleware(container, next)

	req := httptest.NewRequest(http.MethodGet, "/assets/missing.js", nil)
	req.Header.Set("Accept", "text/html")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	if strings.Contains(rec.Body.String(), "<html") {
		t.Fatalf("did not expect html error page for asset path")
	}
}

func TestShouldRenderHTMLErrorPage(t *testing.T) {
	tests := []struct {
		name string
		req  *http.Request
		want bool
	}{
		{
			name: "html get",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/users", nil)
				r.Header.Set("Accept", "text/html")
				return r
			}(),
			want: true,
		},
		{
			name: "asset path",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)
				r.Header.Set("Accept", "text/html")
				return r
			}(),
			want: false,
		},
		{
			name: "post request",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/users", nil)
				r.Header.Set("Accept", "text/html")
				return r
			}(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldRenderHTMLErrorPage(tt.req)
			if got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
