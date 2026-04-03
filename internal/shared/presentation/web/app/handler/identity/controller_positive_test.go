package identityhandler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	identityrepo "github.com/paladignus/actajus/internal/module/identity/application/repository"
	identityusecase "github.com/paladignus/actajus/internal/module/identity/application/usecase"
	handlerctx "github.com/paladignus/actajus/internal/shared/presentation/web/app/handlerctx"
	pagepresenter "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/page"
	webhtml "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/html"
	navigation "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/navigation"
	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

func newIdentityTestHTMLResponder(t *testing.T) *webhtml.Responder {
	t.Helper()
	renderer, err := webtemplate.NewRenderer(webtemplate.Source{
		TemplateFS: fstest.MapFS{
			"templates/pages/sessions.gohtml": {
				Data: []byte(`{{define "pages/sessions"}}<html><body><h1>{{.Page.Title}}</h1></body></html>{{end}}`),
			},
		},
		TemplatePatterns: []string{"templates/pages/*.gohtml"},
		DevServerURL:     "http://localhost:5173",
	})
	if err != nil {
		t.Fatalf("new renderer: %v", err)
	}
	return webhtml.New(pagepresenter.New(renderer), nil)
}

func TestSessionsPageRendersOnSuccess(t *testing.T) {
	controller := Controller{
		Context: handlerctx.Context{
			HTML:        newIdentityTestHTMLResponder(t),
			Navigator:   navigation.New(false),
			FlashCookie: "actajus_flash",
			RequireAuth: func(http.ResponseWriter, *http.Request) (*handlerctx.State, error) {
				return &handlerctx.State{
					Claims:            &identityread.AccessTokenClaims{IDUser: 7, IDSession: 9},
					CanManageSessions: true,
				}, nil
			},
		},
		ViewSessions: identityusecase.NewViewSessions(stubSessionQueryRepository{
			own: []identityread.SessionReadModel{{
				IDSession: 9, IDUser: 7, UserEmail: "user@example.com", ExpiresAt: time.Now(),
			}},
			admin: []identityread.SessionReadModel{{
				IDSession: 10, IDUser: 8, UserEmail: "other@example.com", ExpiresAt: time.Now(),
			}},
		}),
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/sessions", nil)
	controller.SessionsPage().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Sessões") && !strings.Contains(rec.Body.String(), "Sess") {
		t.Fatalf("expected sessions page, got %q", rec.Body.String())
	}
}

type stubSessionQueryRepository struct {
	own   []identityread.SessionReadModel
	admin []identityread.SessionReadModel
	byID  *identityread.SessionReadModel
	err   error
}

func (s stubSessionQueryRepository) GetByID(context.Context, int64) (*identityread.SessionReadModel, error) {
	return s.byID, s.err
}

func (s stubSessionQueryRepository) ListByUser(context.Context, int64) ([]identityread.SessionReadModel, error) {
	return s.own, s.err
}

func (s stubSessionQueryRepository) ListAll(context.Context, identityrepo.SessionQueryFilter) ([]identityread.SessionReadModel, error) {
	return s.admin, s.err
}
