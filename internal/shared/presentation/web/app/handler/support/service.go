package support

import (
	"mime"
	"net/http"
	"path"
	"strings"

	securityheaders "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/securityheaders"
	handlerctx "github.com/paladignus/actajus/internal/shared/presentation/web/app/handlerctx"
	httperror "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/error"
	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

type htmlRenderer interface {
	Render(http.ResponseWriter, int, string, webtemplate.Page) error
}

type Service struct {
	AssetsHandler   http.Handler
	HTML            htmlRenderer
	AuthFromRequest func(http.ResponseWriter, *http.Request) (*handlerctx.State, error)
	Headers         *securityheaders.Adapter
}

func (s Service) Assets() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.AssetsHandler == nil {
			http.NotFound(w, r)
			return
		}
		if s.Headers != nil {
			s.Headers.ApplyAsset(w)
		}
		r2 := r.Clone(r.Context())
		r2.URL.Path = "/assets/" + strings.TrimPrefix(strings.TrimSpace(r.URL.Path), "/")
		s.AssetsHandler.ServeHTTP(w, r2)
	})
}

func (s Service) StaticFile(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.AssetsHandler == nil {
			http.NotFound(w, r)
			return
		}
		if s.Headers != nil {
			s.Headers.ApplyAsset(w)
		}
		if ext := strings.ToLower(path.Ext(strings.TrimSpace(name))); ext != "" {
			if contentType := mime.TypeByExtension(ext); contentType != "" {
				w.Header().Set("Content-Type", contentType)
			}
		}
		r2 := r.Clone(r.Context())
		r2.URL.Path = "/" + strings.TrimPrefix(name, "/")
		s.AssetsHandler.ServeHTTP(w, r2)
	})
}

func (s Service) RenderNotFound(w http.ResponseWriter, r *http.Request) {
	s.renderErrorPage(w, r, http.StatusNotFound, "404", "A pagina solicitada nao foi encontrada.")
}

func (s Service) RenderForbidden(w http.ResponseWriter, r *http.Request) {
	s.renderErrorPage(w, r, http.StatusForbidden, "403", "Voce nao tem permissao para acessar esta pagina.")
}

func (s Service) RespondForbidden(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.RenderForbidden(w, r)
		return
	}
	httperror.Forbidden(w)
}

func (s Service) RespondNotFound(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.RenderNotFound(w, r)
		return
	}
	http.NotFound(w, r)
}

func (s Service) renderErrorPage(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	if s.Headers != nil {
		s.Headers.ApplyDocument(w)
	}
	primaryHref := "/login"
	primaryLabel := "Ir para login"
	shortcuts := []map[string]string{
		{"Label": "Pagina inicial", "Href": "/"},
		{"Label": "Login", "Href": "/login"},
	}
	if s.AuthFromRequest != nil {
		if state, err := s.AuthFromRequest(w, r); err == nil && state != nil && state.Claims != nil {
			primaryHref = "/sessions"
			primaryLabel = "Voltar ao painel"
			shortcuts = []map[string]string{
				{"Label": "Sessoes", "Href": "/sessions"},
				{"Label": "Companies", "Href": "/companies"},
				{"Label": "Usuarios", "Href": "/users"},
				{"Label": "Roles", "Href": "/roles"},
			}
		}
	}
	if err := s.HTML.Render(w, status, "pages/error", webtemplate.Page{
		Title:       code + " | Actajus",
		Description: "Erro de navegacao do painel WEB.",
		NavKey:      "login",
		Entry:       "src/main.ts",
		Data: map[string]any{
			"Code":         code,
			"Message":      message,
			"CurrentPath":  r.URL.Path,
			"Breadcrumbs":  buildErrorBreadcrumbs(r.URL.Path),
			"PrimaryHref":  primaryHref,
			"PrimaryLabel": primaryLabel,
			"Shortcuts":    shortcuts,
		},
	}); err != nil {
		httperror.Render(w, status, message)
	}
}

func buildErrorBreadcrumbs(pathValue string) []map[string]string {
	trimmed := strings.Trim(pathValue, "/")
	if trimmed == "" {
		return []map[string]string{{"Label": "inicio", "Href": "/"}}
	}
	parts := strings.Split(trimmed, "/")
	items := make([]map[string]string, 0, len(parts)+1)
	items = append(items, map[string]string{"Label": "inicio", "Href": "/"})
	current := ""
	for _, part := range parts {
		current += "/" + part
		items = append(items, map[string]string{
			"Label": part,
			"Href":  current,
		})
	}
	return items
}
