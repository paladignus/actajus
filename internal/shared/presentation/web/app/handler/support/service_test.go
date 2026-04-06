package support

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	securityheaders "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/securityheaders"
	handlerctx "github.com/paladignus/actajus/internal/shared/presentation/web/app/handlerctx"
	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

type stubHTMLRenderer struct {
	status       int
	templateName string
	page         webtemplate.Page
	renderCalls  int
}

func (s *stubHTMLRenderer) Render(_ http.ResponseWriter, status int, templateName string, page webtemplate.Page) error {
	s.status = status
	s.templateName = templateName
	s.page = page
	s.renderCalls++
	return nil
}

func TestServiceRenderNotFoundForGuest(t *testing.T) {
	renderer := &stubHTMLRenderer{}
	service := Service{
		HTML: renderer,
		AuthFromRequest: func(http.ResponseWriter, *http.Request) (*handlerctx.State, error) {
			return nil, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/invalida/rota", nil)
	rec := httptest.NewRecorder()
	service.RenderNotFound(rec, req)

	if renderer.renderCalls != 1 {
		t.Fatalf("expected one render call, got %d", renderer.renderCalls)
	}
	if renderer.status != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, renderer.status)
	}
	if renderer.templateName != "pages/error" {
		t.Fatalf("unexpected template %q", renderer.templateName)
	}

	data, ok := renderer.page.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected map page data")
	}
	if data["PrimaryHref"] != "/login" {
		t.Fatalf("expected guest primary href /login, got %#v", data["PrimaryHref"])
	}
	if data["PrimaryLabel"] != "Ir para login" {
		t.Fatalf("unexpected primary label %#v", data["PrimaryLabel"])
	}
	shortcuts, ok := data["Shortcuts"].([]map[string]string)
	if !ok || len(shortcuts) != 2 {
		t.Fatalf("expected guest shortcuts, got %#v", data["Shortcuts"])
	}
	breadcrumbs, ok := data["Breadcrumbs"].([]map[string]string)
	if !ok || len(breadcrumbs) != 3 {
		t.Fatalf("expected breadcrumbs for /invalida/rota, got %#v", data["Breadcrumbs"])
	}
}

func TestServiceRenderForbiddenForAuthenticatedUser(t *testing.T) {
	renderer := &stubHTMLRenderer{}
	service := Service{
		HTML: renderer,
		AuthFromRequest: func(http.ResponseWriter, *http.Request) (*handlerctx.State, error) {
			return &handlerctx.State{
				Claims: &identityread.AccessTokenClaims{
					IDUser:    7,
					IDSession: 14,
				},
			}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/roles", nil)
	rec := httptest.NewRecorder()
	service.RenderForbidden(rec, req)

	data := renderer.page.Data.(map[string]any)
	if data["PrimaryHref"] != "/sessions" {
		t.Fatalf("expected authenticated primary href /sessions, got %#v", data["PrimaryHref"])
	}
	if data["PrimaryLabel"] != "Voltar ao painel" {
		t.Fatalf("unexpected primary label %#v", data["PrimaryLabel"])
	}
	shortcuts := data["Shortcuts"].([]map[string]string)
	if len(shortcuts) != 4 {
		t.Fatalf("expected authenticated shortcuts, got %#v", shortcuts)
	}
}

func TestBuildErrorBreadcrumbs(t *testing.T) {
	got := buildErrorBreadcrumbs("/companies/42/edit")
	if len(got) != 4 {
		t.Fatalf("expected 4 breadcrumb items, got %d", len(got))
	}
	if got[0]["Href"] != "/" || got[1]["Href"] != "/companies" || got[2]["Href"] != "/companies/42" || got[3]["Href"] != "/companies/42/edit" {
		t.Fatalf("unexpected breadcrumbs %#v", got)
	}
}

func TestServiceStaticFileServesAssetWithContentType(t *testing.T) {
	assetFS := fstest.MapFS{
		"favicon.svg": {Data: []byte(`<svg></svg>`)},
	}
	root, err := fs.Sub(assetFS, ".")
	if err != nil {
		t.Fatalf("fs sub: %v", err)
	}
	service := Service{
		AssetsHandler: http.FileServerFS(root),
		Headers:       securityheaders.New(""),
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/favicon.svg", nil)
	service.StaticFile("favicon.svg").ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if contentType := rec.Header().Get("Content-Type"); contentType == "" {
		t.Fatalf("expected content type header")
	}
	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("expected nosniff header")
	}
}

func TestServiceRenderErrorAppliesDocumentHeaders(t *testing.T) {
	renderer := &stubHTMLRenderer{}
	service := Service{
		HTML:    renderer,
		Headers: securityheaders.New(""),
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/nao-existe", nil)
	service.RenderNotFound(rec, req)

	if rec.Header().Get("Content-Security-Policy") == "" {
		t.Fatalf("expected content security policy header")
	}
	if rec.Header().Get("X-Frame-Options") != "DENY" {
		t.Fatalf("expected frame options header")
	}
}
