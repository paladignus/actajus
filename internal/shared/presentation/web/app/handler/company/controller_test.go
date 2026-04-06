package companyhandler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	companyread "github.com/paladignus/actajus/internal/module/company/application/readmodel"
	companyrepo "github.com/paladignus/actajus/internal/module/company/application/repository"
	companyusecase "github.com/paladignus/actajus/internal/module/company/application/usecase"
	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	handlerctx "github.com/paladignus/actajus/internal/shared/presentation/web/app/handlerctx"
	pagepresenter "github.com/paladignus/actajus/internal/shared/presentation/web/app/presenter/page"
	webhtml "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/html"
	navigation "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/navigation"
	webtemplate "github.com/paladignus/actajus/internal/shared/presentation/web/template"
)

func newCompanyTestHTMLResponder(t *testing.T) *webhtml.Responder {
	t.Helper()
	renderer, err := webtemplate.NewRenderer(webtemplate.Source{
		TemplateFS: fstest.MapFS{
			"templates/pages/companies.gohtml": {
				Data: []byte(`{{define "pages/companies"}}<html><body><h1>{{.Page.Title}}</h1></body></html>{{end}}`),
			},
			"templates/pages/company_form.gohtml": {
				Data: []byte(`{{define "pages/company-form"}}<html><body><h1>{{.Page.Title}}</h1></body></html>{{end}}`),
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

func authenticatedCompanyContext(t *testing.T) handlerctx.Context {
	t.Helper()
	return handlerctx.Context{
		HTML:      newCompanyTestHTMLResponder(t),
		Navigator: navigation.New(false),
		RequireAuth: func(http.ResponseWriter, *http.Request) (*handlerctx.State, error) {
			return &handlerctx.State{}, nil
		},
		RespondNotFound: func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		},
	}
}

func TestCompaniesPageRedirectsWhenUnauthenticated(t *testing.T) {
	controller := Controller{
		Context: handlerctx.Context{
			Navigator: navigation.New(false),
			RequireAuth: func(http.ResponseWriter, *http.Request) (*handlerctx.State, error) {
				http.Redirect(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/companies", nil), "/login", http.StatusSeeOther)
				return nil, errors.New("missing auth")
			},
		},
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/companies", nil)
	controller.Context.RequireAuth = func(w http.ResponseWriter, r *http.Request) (*handlerctx.State, error) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil, errors.New("missing auth")
	}
	controller.CompaniesPage().ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	if location := rec.Header().Get("Location"); location != "/login" {
		t.Fatalf("expected redirect to /login, got %q", location)
	}
}

func TestCompanyNewPageRendersForm(t *testing.T) {
	controller := Controller{Context: authenticatedCompanyContext(t)}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/companies/new", nil)
	controller.CompanyNewPage().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Company Form") {
		t.Fatalf("expected company form page, got %q", rec.Body.String())
	}
}

func TestCompanyDeleteActionRespondsNotFoundOnInvalidID(t *testing.T) {
	controller := Controller{Context: authenticatedCompanyContext(t)}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/companies/invalid/delete", nil)
	req.SetPathValue("id", "invalid")
	controller.CompanyDeleteAction().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestCompaniesPageRendersOnSuccess(t *testing.T) {
	controller := Controller{
		Context: authenticatedCompanyContext(t),
		List: companyusecase.NewListCompanies(stubCompanyReadRepository{
			list: &companyread.CompanyListReadModel{
				Data: []companyread.CompanyReadModel{{ID: 1, Name: "Acme", TradeName: "Acme", CNPJ: "123"}},
			},
		}),
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/companies", nil)
	controller.CompaniesPage().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Companies") {
		t.Fatalf("expected companies page, got %q", rec.Body.String())
	}
}

func TestCompanyCreateActionReturnsBadRequestOnInvalidNumber(t *testing.T) {
	controller := Controller{
		Context: handlerctx.Context{
			Navigator: navigation.New(false),
			RequireAuth: func(http.ResponseWriter, *http.Request) (*handlerctx.State, error) {
				return &handlerctx.State{Claims: &identityread.AccessTokenClaims{IDUser: 7, IDSession: 9}}, nil
			},
		},
	}

	body := strings.NewReader("name=Acme&trade_name=Acme&cnpj=123&zip=79000-000&title=Rua&street=Central&number=abc&neighborhood=Centro&city=Campo+Grande&state=MS&country=BR&phone_number=123&phone_kind=commercial&email_address=user%40example.com")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/companies", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	controller.CompanyCreateAction().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

type stubCompanyReadRepository struct {
	list *companyread.CompanyListReadModel
	err  error
}

func (s stubCompanyReadRepository) List(context.Context, companyrepo.CompanyListFilter, *string, *string, int, string) (*companyread.CompanyListReadModel, error) {
	return s.list, s.err
}
