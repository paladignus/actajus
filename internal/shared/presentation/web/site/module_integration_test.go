package site

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/paladignus/actajus/internal/module/company"
	companyread "github.com/paladignus/actajus/internal/module/company/application/readmodel"
	companyrepo "github.com/paladignus/actajus/internal/module/company/application/repository"
	companyusecase "github.com/paladignus/actajus/internal/module/company/application/usecase"
	companydomain "github.com/paladignus/actajus/internal/module/company/domain"
	"github.com/paladignus/actajus/internal/module/identity"
	identitymapper "github.com/paladignus/actajus/internal/module/identity/application/mapper"
	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	identityrepo "github.com/paladignus/actajus/internal/module/identity/application/repository"
	identityusecase "github.com/paladignus/actajus/internal/module/identity/application/usecase"
	identitydomain "github.com/paladignus/actajus/internal/module/identity/domain"
	socialMediaDomain "github.com/paladignus/actajus/internal/module/social_media/domain"
	"github.com/paladignus/actajus/internal/shared/application/pagination"
	sharedrepo "github.com/paladignus/actajus/internal/shared/application/repository"
	"github.com/paladignus/actajus/internal/shared/application/uow"
	"github.com/paladignus/actajus/internal/shared/presentation/validation"
	webcsrf "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/csrf"
)

type stubLogger struct{}

func (stubLogger) Debug(context.Context, string, ...any) {}
func (stubLogger) Info(context.Context, string, ...any)  {}
func (stubLogger) Warn(context.Context, string, ...any)  {}
func (stubLogger) Error(context.Context, string, ...any) {}
func (stubLogger) With(...any) sharedrepo.Logger         { return stubLogger{} }
func (stubLogger) WithError(error) sharedrepo.Logger     { return stubLogger{} }

type manifestEntry struct {
	File string   `json:"file"`
	CSS  []string `json:"css"`
}

type permissionCheckerStub struct{}

func (permissionCheckerStub) HasPermission(_ context.Context, _ int64, perm string) (bool, error) {
	return perm == "session:manage" || perm == "rbac:manage", nil
}

type accessTokenServiceStub struct{}

func (accessTokenServiceStub) Sign(identityread.AccessTokenClaims) (string, error) {
	return "", nil
}

func (accessTokenServiceStub) Verify(token string) (identityread.AccessTokenClaims, error) {
	if token != "valid-token" {
		return identityread.AccessTokenClaims{}, identitydomain.ErrInvalidToken
	}
	return identityread.AccessTokenClaims{IDUser: 7, IDSession: 11}, nil
}

type clockStub struct{}

func (clockStub) Now() time.Time { return time.Unix(0, 0) }

type sessionRepositoryStub struct {
	revokedIDs        []int64
	revokedAllUserIDs []int64
}

func (*sessionRepositoryStub) Create(context.Context, *identitydomain.Session) error { return nil }
func (*sessionRepositoryStub) GetByID(context.Context, int64) (*identitydomain.Session, error) {
	return nil, nil
}
func (*sessionRepositoryStub) RotateRefreshToken(context.Context, int64, [32]byte, time.Time) error {
	return nil
}
func (s *sessionRepositoryStub) Revoke(_ context.Context, id int64) error {
	s.revokedIDs = append(s.revokedIDs, id)
	return nil
}
func (s *sessionRepositoryStub) RevokeAllByUser(_ context.Context, idUser int64) error {
	s.revokedAllUserIDs = append(s.revokedAllUserIDs, idUser)
	return nil
}
func (*sessionRepositoryStub) CountActiveByUser(context.Context, int64) (int, error) {
	return 0, nil
}
func (*sessionRepositoryStub) IsActive(context.Context, int64, int64, time.Time) (bool, error) {
	return true, nil
}
func (*sessionRepositoryStub) RotateRefreshTokenAtomic(context.Context, int64, [32]byte, [32]byte, time.Time, time.Time) (bool, error) {
	return true, nil
}

type sessionQueryRepositoryStub struct{}

func (sessionQueryRepositoryStub) GetByID(context.Context, int64) (*identityread.SessionReadModel, error) {
	return &identityread.SessionReadModel{IDSession: 11, IDUser: 7, UserEmail: "user@mail.com"}, nil
}
func (sessionQueryRepositoryStub) ListByUser(context.Context, int64) ([]identityread.SessionReadModel, error) {
	return []identityread.SessionReadModel{
		{IDSession: 11, IDUser: 7, UserEmail: "user@mail.com", UserAgent: "GoTest"},
	}, nil
}
func (sessionQueryRepositoryStub) ListAll(context.Context, identityrepo.SessionQueryFilter) ([]identityread.SessionReadModel, error) {
	return nil, nil
}

type catalogQueryRepositoryStub struct{}

func (catalogQueryRepositoryStub) ListUsers(context.Context, identityrepo.CatalogUserFilter) ([]identityread.UserListItemReadModel, error) {
	return []identityread.UserListItemReadModel{{ID: 7, Email: "user@mail.com"}}, nil
}
func (catalogQueryRepositoryStub) ListPermissions(context.Context, identityrepo.CatalogPermissionFilter) ([]identityread.PermissionListItemReadModel, error) {
	return []identityread.PermissionListItemReadModel{{ID: 1, Resource: "session", Action: "manage", Description: "Gerenciar sessoes"}}, nil
}
func (catalogQueryRepositoryStub) ListRoles(context.Context) ([]identityread.RoleListItemReadModel, error) {
	return []identityread.RoleListItemReadModel{{ID: 1, Name: "admin", Description: "Administrador"}}, nil
}
func (catalogQueryRepositoryStub) ListRolesByUser(context.Context, int64) ([]identityread.UserRoleAssignmentReadModel, error) {
	return nil, nil
}
func (catalogQueryRepositoryStub) ListPermissionsByRole(context.Context, int16) ([]identityread.RolePermissionAssignmentReadModel, error) {
	return nil, nil
}

type companyReadRepositoryStub struct{}

func (companyReadRepositoryStub) List(context.Context, companyrepo.CompanyListFilter, *string, *string, int, string) (*companyread.CompanyListReadModel, error) {
	return &companyread.CompanyListReadModel{
		Data: []companyread.CompanyReadModel{
			{ID: 21, Name: "Acme LTDA", CNPJ: "12345678000199"},
		},
		PageInfo: pagination.PageInfo{},
	}, nil
}

type roleAdminRepositoryStub struct {
	createdName string
	updatedID   int16
	updatedName string
	deletedID   int16
}

func (s *roleAdminRepositoryStub) Create(_ context.Context, name, _ string) error {
	s.createdName = name
	return nil
}
func (s *roleAdminRepositoryStub) Update(_ context.Context, id int16, name, _ string) error {
	s.updatedID = id
	s.updatedName = name
	return nil
}
func (s *roleAdminRepositoryStub) Delete(_ context.Context, id int16) error {
	s.deletedID = id
	return nil
}

type permissionAdminRepositoryStub struct {
	resource        string
	action          string
	updatedID       int16
	updatedResource string
	updatedAction   string
	deletedID       int16
}

func (s *permissionAdminRepositoryStub) Create(_ context.Context, resource, action, _ string) error {
	s.resource = resource
	s.action = action
	return nil
}
func (s *permissionAdminRepositoryStub) Update(_ context.Context, id int16, resource, action, _ string) error {
	s.updatedID = id
	s.updatedResource = resource
	s.updatedAction = action
	return nil
}
func (s *permissionAdminRepositoryStub) Delete(_ context.Context, id int16) error {
	s.deletedID = id
	return nil
}

type roleUserAdminRepositoryStub struct {
	assignedUser int64
	assignedRole int16
	removedUser  int64
	removedRole  int16
}

func (s *roleUserAdminRepositoryStub) AssignRole(_ context.Context, uid int64, rid int16, _ int64) error {
	s.assignedUser = uid
	s.assignedRole = rid
	return nil
}
func (s *roleUserAdminRepositoryStub) RemoveRole(_ context.Context, uid int64, rid int16) error {
	s.removedUser = uid
	s.removedRole = rid
	return nil
}
func (*roleUserAdminRepositoryStub) ListUserIDsByRole(context.Context, int16) ([]int64, error) {
	return []int64{7}, nil
}

type permissionRoleAdminRepositoryStub struct {
	roleID              int16
	permissionID        int16
	revokedRoleID       int16
	revokedPermissionID int16
}

func (s *permissionRoleAdminRepositoryStub) GrantPermission(_ context.Context, rid int16, pid int16) error {
	s.roleID = rid
	s.permissionID = pid
	return nil
}
func (s *permissionRoleAdminRepositoryStub) RevokePermission(_ context.Context, rid int16, pid int16) error {
	s.revokedRoleID = rid
	s.revokedPermissionID = pid
	return nil
}

type roleUserQueryRepositoryStub struct{}

func (roleUserQueryRepositoryStub) ListUserIDsByRole(context.Context, int16) ([]int64, error) {
	return []int64{7}, nil
}

type cacheInvalidatorStub struct{}

func (cacheInvalidatorStub) InvalidateUser(context.Context, int64) {}

type roleUsersIndexStub struct{}

func (roleUsersIndexStub) AddUserToRole(context.Context, int16, int64)      {}
func (roleUsersIndexStub) RemoveUserFromRole(context.Context, int16, int64) {}
func (roleUsersIndexStub) ListUsersByRole(context.Context, int16) ([]int64, error) {
	return []int64{7}, nil
}
func (roleUsersIndexStub) AddUsersToRole(context.Context, int16, []int64) {}

type uowStub struct{}

func (uowStub) Do(ctx context.Context, fn func(tx uow.Tx) error) error {
	return fn(struct{}{})
}

type companyRepositoryStub struct {
	deletedCompanyID int64
}

func (*companyRepositoryStub) Create(context.Context, *companydomain.Company) error { return nil }
func (*companyRepositoryStub) Update(context.Context, *companydomain.Company) error { return nil }
func (s *companyRepositoryStub) Delete(_ context.Context, company *companydomain.Company) error {
	s.deletedCompanyID = company.ID().Value()
	return nil
}
func (*companyRepositoryStub) FindByCNPJ(context.Context, string) (*companydomain.Company, error) {
	return nil, nil
}
func (*companyRepositoryStub) FindByID(_ context.Context, id int64) (*companydomain.Company, error) {
	return companydomain.NewCompanyBuilder().
		WithID(id).
		WithRegisteredBy(7).
		WithName("Acme LTDA").
		Build()
}

type companyLinkRepositoryStub struct {
	deletedCompanyID int64
}

func (*companyLinkRepositoryStub) Create(context.Context, int64, int64) error { return nil }
func (s *companyLinkRepositoryStub) DeleteByIDCompany(_ context.Context, idCompany int64) error {
	s.deletedCompanyID = idCompany
	return nil
}

type socialMediaRepositoryDeleteStub struct {
	deletedCompanyID int64
}

func (*socialMediaRepositoryDeleteStub) Create(context.Context, *socialMediaDomain.SocialMedia) error {
	return nil
}
func (*socialMediaRepositoryDeleteStub) Update(context.Context, socialMediaDomain.SocialMedia) error {
	return nil
}
func (*socialMediaRepositoryDeleteStub) Delete(context.Context, socialMediaDomain.SocialMedia) error {
	return nil
}
func (*socialMediaRepositoryDeleteStub) FindByIDCompany(context.Context, int64) ([]*socialMediaDomain.SocialMedia, error) {
	return nil, nil
}
func (s *socialMediaRepositoryDeleteStub) DeleteByIDCompany(_ context.Context, idCompany int64) error {
	s.deletedCompanyID = idCompany
	return nil
}

type companyFactoryDeleteStub struct {
	company     *companyRepositoryStub
	addressLink *companyLinkRepositoryStub
	phoneLink   *companyLinkRepositoryStub
	emailLink   *companyLinkRepositoryStub
	socialMedia *socialMediaRepositoryDeleteStub
}

func (f *companyFactoryDeleteStub) WithTx(uow.Tx) companyrepo.Factory      { return f }
func (f *companyFactoryDeleteStub) Company() companyrepo.CompanyRepository { return f.company }
func (f *companyFactoryDeleteStub) Address() companyrepo.AddressRepository { return nil }
func (f *companyFactoryDeleteStub) CompanyAddress() companydomain.CompanyAddressRepository {
	return f.addressLink
}
func (f *companyFactoryDeleteStub) Phone() companyrepo.PhoneRepository { return nil }
func (f *companyFactoryDeleteStub) CompanyPhone() companydomain.CompanyPhoneRepository {
	return f.phoneLink
}
func (f *companyFactoryDeleteStub) Email() companyrepo.EmailRepository { return nil }
func (f *companyFactoryDeleteStub) CompanyEmail() companydomain.CompanyEmailRepository {
	return f.emailLink
}
func (f *companyFactoryDeleteStub) SocialMedia() companyrepo.SocialMediaRepository {
	return f.socialMedia
}

func newTestModule(t *testing.T) Module {
	t.Helper()
	module, err := NewModule(Dependencies{
		Logger:   stubLogger{},
		Company:  company.Module{},
		Identity: identity.Module{},
		Secure:   false,
	})
	if err != nil {
		t.Fatalf("new web module: %v", err)
	}
	return module
}

func newAuthenticatedTestModule(t *testing.T) Module {
	t.Helper()
	sessionRepo := &sessionRepositoryStub{}
	module, err := NewModule(Dependencies{
		Logger:  stubLogger{},
		Company: company.Module{List: companyusecase.NewListCompanies(companyReadRepositoryStub{})},
		Identity: identity.Module{
			ValidateAccess:  identityusecase.NewValidateAccess(accessTokenServiceStub{}, sessionRepo, clockStub{}, true),
			ViewSessions:    identityusecase.NewViewSessions(sessionQueryRepositoryStub{}),
			ListUsers:       identityusecase.NewListUsers(catalogQueryRepositoryStub{}),
			ListPermissions: identityusecase.NewListPermissions(catalogQueryRepositoryStub{}),
			ListRoles:       identityusecase.NewListRoles(catalogQueryRepositoryStub{}),
		},
		Secure: false,
	})
	if err != nil {
		t.Fatalf("new authenticated web module: %v", err)
	}
	return module
}

func newAuthenticatedActionTestModule(t *testing.T) (Module, *sessionRepositoryStub) {
	t.Helper()
	sessionRepo := &sessionRepositoryStub{}
	authMapper := identitymapper.NewAuthMapper(validation.New())
	module, err := NewModule(Dependencies{
		Logger:  stubLogger{},
		Company: company.Module{List: companyusecase.NewListCompanies(companyReadRepositoryStub{})},
		Identity: identity.Module{
			ValidateAccess: identityusecase.NewValidateAccess(accessTokenServiceStub{}, sessionRepo, clockStub{}, true),
			Logout:         identityusecase.NewLogout(sessionRepo, *authMapper),
			LogoutAll:      identityusecase.NewLogoutAll(sessionRepo, *authMapper),
			RevokeSession:  identityusecase.NewRevokeSession(sessionRepo, *authMapper),
			ViewSessions:   identityusecase.NewViewSessions(sessionQueryRepositoryStub{}),
		},
		Secure: false,
	})
	if err != nil {
		t.Fatalf("new authenticated action web module: %v", err)
	}
	return module, sessionRepo
}

func newAdminActionTestModule(t *testing.T) (Module, *roleAdminRepositoryStub, *permissionAdminRepositoryStub, *roleUserAdminRepositoryStub, *permissionRoleAdminRepositoryStub) {
	t.Helper()
	roleRepo := &roleAdminRepositoryStub{}
	permRepo := &permissionAdminRepositoryStub{}
	roleUserRepo := &roleUserAdminRepositoryStub{}
	permRoleRepo := &permissionRoleAdminRepositoryStub{}
	authMapper := identitymapper.NewAuthMapper(validation.New())
	rbacMapper := identitymapper.NewRBACAdminMapper(validation.New())

	module, err := NewModule(Dependencies{
		Logger:  stubLogger{},
		Company: company.Module{},
		Identity: identity.Module{
			ValidateAccess:           identityusecase.NewValidateAccess(accessTokenServiceStub{}, &sessionRepositoryStub{}, clockStub{}, true),
			ListUsers:                identityusecase.NewListUsers(catalogQueryRepositoryStub{}),
			ListPermissions:          identityusecase.NewListPermissions(catalogQueryRepositoryStub{}),
			ListRoles:                identityusecase.NewListRoles(catalogQueryRepositoryStub{}),
			CreateRole:               identityusecase.NewCreateRole(roleRepo, validation.New()),
			UpdateRole:               identityusecase.NewUpdateRole(roleRepo, validation.New()),
			DeleteRole:               identityusecase.NewDeleteRole(roleRepo),
			CreatePermission:         identityusecase.NewCreatePermission(permRepo, validation.New()),
			UpdatePermission:         identityusecase.NewUpdatePermission(permRepo, validation.New()),
			DeletePermission:         identityusecase.NewDeletePermission(permRepo),
			AssignRoleToUser:         identityusecase.NewAssignRoleToUser(roleUserRepo, cacheInvalidatorStub{}, roleUsersIndexStub{}, rbacMapper),
			RemoveRoleFromUser:       identityusecase.NewRemoveRoleFromUser(roleUserRepo, cacheInvalidatorStub{}, roleUsersIndexStub{}, rbacMapper),
			GrantPermissionToRole:    identityusecase.NewGrantPermissionToRole(roleUserQueryRepositoryStub{}, permRoleRepo, cacheInvalidatorStub{}, rbacMapper),
			RevokePermissionFromRole: identityusecase.NewRevokePermissionFromRole(roleUserQueryRepositoryStub{}, permRoleRepo, cacheInvalidatorStub{}, rbacMapper),
			RBACChecker:              permissionCheckerStub{},
			Logout:                   identityusecase.NewLogout(&sessionRepositoryStub{}, *authMapper),
		},
		Secure: false,
	})
	if err != nil {
		t.Fatalf("new admin action web module: %v", err)
	}
	return module, roleRepo, permRepo, roleUserRepo, permRoleRepo
}

func newCompanyDeleteActionTestModule(t *testing.T) (Module, *companyFactoryDeleteStub) {
	t.Helper()
	factory := &companyFactoryDeleteStub{
		company:     &companyRepositoryStub{},
		addressLink: &companyLinkRepositoryStub{},
		phoneLink:   &companyLinkRepositoryStub{},
		emailLink:   &companyLinkRepositoryStub{},
		socialMedia: &socialMediaRepositoryDeleteStub{},
	}
	module, err := NewModule(Dependencies{
		Logger: stubLogger{},
		Company: company.Module{
			Delete: companyusecase.NewDeleteCompany(uowStub{}, factory),
		},
		Identity: identity.Module{
			ValidateAccess: identityusecase.NewValidateAccess(accessTokenServiceStub{}, &sessionRepositoryStub{}, clockStub{}, true),
		},
		Secure: false,
	})
	if err != nil {
		t.Fatalf("new company delete action web module: %v", err)
	}
	return module, factory
}

func currentAssetPaths(t *testing.T) (string, string) {
	t.Helper()
	raw, err := contentFS.ReadFile("content/dist/.vite/manifest.json")
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	manifest := map[string]manifestEntry{}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatalf("unmarshal manifest: %v", err)
	}
	entry := manifest["src/main.ts"]
	if entry.File == "" || len(entry.CSS) == 0 {
		t.Fatalf("unexpected manifest entry %#v", entry)
	}
	return "/assets/" + path.Base(entry.File), "/assets/" + path.Base(entry.CSS[0])
}

func cookieValue(rec *httptest.ResponseRecorder, name string) string {
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return ""
}

func TestMountLoginPage(t *testing.T) {
	module := newTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(strings.ToLower(rec.Body.String()), "entrar") {
		t.Fatalf("expected login page body, got %q", rec.Body.String())
	}
	if cookie := rec.Result().Cookies(); len(cookie) == 0 || cookie[0].Name != webcsrf.CookieName {
		t.Fatalf("expected csrf cookie to be set")
	}
	if rec.Header().Get("Content-Security-Policy") == "" {
		t.Fatalf("expected content security policy header")
	}
}

func TestMountRootRedirectsToLoginWhenGuest(t *testing.T) {
	module := newTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	if location := rec.Header().Get("Location"); location != "/login" {
		t.Fatalf("expected redirect to /login, got %q", location)
	}
}

func TestMountServesFavicon(t *testing.T) {
	module := newTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/favicon.svg", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if contentType := rec.Header().Get("Content-Type"); !strings.Contains(contentType, "image/svg+xml") {
		t.Fatalf("expected svg content type, got %q", contentType)
	}
	if rec.Body.Len() == 0 {
		t.Fatalf("expected favicon body")
	}
}

func TestMountServesCurrentAssets(t *testing.T) {
	module := newTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)
	jsPath, cssPath := currentAssetPaths(t)

	for _, path := range []string{jsPath, cssPath} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status %d for %s, got %d", http.StatusOK, path, rec.Code)
		}
		if rec.Body.Len() == 0 {
			t.Fatalf("expected asset body for %s", path)
		}
		if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Fatalf("expected nosniff header for %s", path)
		}
	}
}

func TestRenderNotFoundPage(t *testing.T) {
	module := newTestModule(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/nao-existe", nil)
	module.RenderNotFound(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "404") {
		t.Fatalf("expected 404 body, got %q", rec.Body.String())
	}
}

func TestMountRejectsPostWithoutCSRF(t *testing.T) {
	module := newTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
}

func TestMountAllowsPostWithMatchingCSRF(t *testing.T) {
	module := newTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader("_csrf=test-token"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: "test-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	if location := rec.Header().Get("Location"); location != "/login" {
		t.Fatalf("expected redirect to /login, got %q", location)
	}
}

func TestMountSessionsPageRendersForAuthenticatedUser(t *testing.T) {
	module := newAuthenticatedTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/sessions", nil)
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Minhas sessoes") {
		t.Fatalf("expected sessions page, got %q", rec.Body.String())
	}
}

func TestMountUsersPageRendersForAuthenticatedUser(t *testing.T) {
	module := newAuthenticatedTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Usuarios") {
		t.Fatalf("expected users page, got %q", rec.Body.String())
	}
}

func TestMountCompaniesPageRendersForAuthenticatedUser(t *testing.T) {
	module := newAuthenticatedTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/companies", nil)
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Nova company") || !strings.Contains(body, "Acme LTDA") {
		t.Fatalf("expected companies page, got %q", body)
	}
}

func TestMountLogoutActionClearsAuthSession(t *testing.T) {
	module, sessionRepo := newAuthenticatedActionTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader("_csrf=test-token"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: "test-token"})
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	if location := rec.Header().Get("Location"); location != "/login" {
		t.Fatalf("expected redirect to /login, got %q", location)
	}
	if len(sessionRepo.revokedIDs) != 1 || sessionRepo.revokedIDs[0] != 11 {
		t.Fatalf("expected revoked session 11, got %#v", sessionRepo.revokedIDs)
	}
	if len(rec.Result().Cookies()) < 3 {
		t.Fatalf("expected auth clearing cookies")
	}
}

func TestMountRevokeCurrentSessionRedirectsToLogin(t *testing.T) {
	module, sessionRepo := newAuthenticatedActionTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/sessions/revoke", strings.NewReader("_csrf=test-token&id_session=11"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: "test-token"})
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	if location := rec.Header().Get("Location"); location != "/login" {
		t.Fatalf("expected redirect to /login, got %q", location)
	}
	if len(sessionRepo.revokedIDs) == 0 || sessionRepo.revokedIDs[len(sessionRepo.revokedIDs)-1] != 11 {
		t.Fatalf("expected revoked current session, got %#v", sessionRepo.revokedIDs)
	}
}

func TestMountRevokeAllOwnSessionsRedirectsToLogin(t *testing.T) {
	module, sessionRepo := newAuthenticatedActionTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/sessions/revoke-all", strings.NewReader("_csrf=test-token"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: "test-token"})
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}
	if location := rec.Header().Get("Location"); location != "/login" {
		t.Fatalf("expected redirect to /login, got %q", location)
	}
	if len(sessionRepo.revokedAllUserIDs) != 1 || sessionRepo.revokedAllUserIDs[0] != 7 {
		t.Fatalf("expected revoke-all for user 7, got %#v", sessionRepo.revokedAllUserIDs)
	}
}

func TestMountCreateRoleActionForAdmin(t *testing.T) {
	module, roleRepo, _, _, _ := newAdminActionTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/roles", strings.NewReader("_csrf=test-token&name=admin-web&description=Administrador"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: "test-token"})
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/roles" {
		t.Fatalf("expected redirect to /roles, got status=%d location=%q", rec.Code, rec.Header().Get("Location"))
	}
	if roleRepo.createdName != "admin-web" {
		t.Fatalf("expected created role name, got %q", roleRepo.createdName)
	}
}

func TestMountCreatePermissionActionForAdmin(t *testing.T) {
	module, _, permRepo, _, _ := newAdminActionTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/permissions", strings.NewReader("_csrf=test-token&resource=rbac&action=manage&description=Gerenciar%20RBAC"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: "test-token"})
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/permissions" {
		t.Fatalf("expected redirect to /permissions, got status=%d location=%q", rec.Code, rec.Header().Get("Location"))
	}
	if permRepo.resource != "rbac" || permRepo.action != "manage" {
		t.Fatalf("expected created permission, got resource=%q action=%q", permRepo.resource, permRepo.action)
	}
}

func TestMountAssignRoleToUserActionForAdmin(t *testing.T) {
	module, _, _, roleUserRepo, _ := newAdminActionTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users/roles", strings.NewReader("_csrf=test-token&id_user=7&id_role=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: "test-token"})
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/users" {
		t.Fatalf("expected redirect to /users, got status=%d location=%q", rec.Code, rec.Header().Get("Location"))
	}
	if roleUserRepo.assignedUser != 7 || roleUserRepo.assignedRole != 1 {
		t.Fatalf("expected assigned role to user, got user=%d role=%d", roleUserRepo.assignedUser, roleUserRepo.assignedRole)
	}
}

func TestMountGrantPermissionToRoleActionForAdmin(t *testing.T) {
	module, _, _, _, permRoleRepo := newAdminActionTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/roles/permissions", strings.NewReader("_csrf=test-token&id_role=1&id_permission=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: "test-token"})
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/permissions" {
		t.Fatalf("expected redirect to /permissions, got status=%d location=%q", rec.Code, rec.Header().Get("Location"))
	}
	if permRoleRepo.roleID != 1 || permRoleRepo.permissionID != 1 {
		t.Fatalf("expected granted permission to role, got role=%d permission=%d", permRoleRepo.roleID, permRoleRepo.permissionID)
	}
}

func TestMountUpdateRoleActionForAdmin(t *testing.T) {
	module, roleRepo, _, _, _ := newAdminActionTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/roles/2", strings.NewReader("_csrf=test-token&name=editor-web&description=Editor"))
	req.SetPathValue("id", "2")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: "test-token"})
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/roles" {
		t.Fatalf("expected redirect to /roles, got status=%d location=%q", rec.Code, rec.Header().Get("Location"))
	}
	if roleRepo.updatedID != 2 || roleRepo.updatedName != "editor-web" {
		t.Fatalf("expected updated role, got id=%d name=%q", roleRepo.updatedID, roleRepo.updatedName)
	}
}

func TestMountDeleteRoleActionForAdmin(t *testing.T) {
	module, roleRepo, _, _, _ := newAdminActionTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/roles/2/delete", strings.NewReader("_csrf=test-token"))
	req.SetPathValue("id", "2")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: "test-token"})
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/roles" {
		t.Fatalf("expected redirect to /roles, got status=%d location=%q", rec.Code, rec.Header().Get("Location"))
	}
	if roleRepo.deletedID != 2 {
		t.Fatalf("expected deleted role 2, got %d", roleRepo.deletedID)
	}
}

func TestMountUpdatePermissionActionForAdmin(t *testing.T) {
	module, _, permRepo, _, _ := newAdminActionTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/permissions/3", strings.NewReader("_csrf=test-token&resource=company&action=write&description=Editar%20company"))
	req.SetPathValue("id", "3")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: "test-token"})
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/permissions" {
		t.Fatalf("expected redirect to /permissions, got status=%d location=%q", rec.Code, rec.Header().Get("Location"))
	}
	if permRepo.updatedID != 3 || permRepo.updatedResource != "company" || permRepo.updatedAction != "write" {
		t.Fatalf("expected updated permission, got id=%d resource=%q action=%q", permRepo.updatedID, permRepo.updatedResource, permRepo.updatedAction)
	}
}

func TestMountDeletePermissionActionForAdmin(t *testing.T) {
	module, _, permRepo, _, _ := newAdminActionTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/permissions/3/delete", strings.NewReader("_csrf=test-token"))
	req.SetPathValue("id", "3")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: "test-token"})
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/permissions" {
		t.Fatalf("expected redirect to /permissions, got status=%d location=%q", rec.Code, rec.Header().Get("Location"))
	}
	if permRepo.deletedID != 3 {
		t.Fatalf("expected deleted permission 3, got %d", permRepo.deletedID)
	}
}

func TestMountRemoveRoleFromUserActionForAdmin(t *testing.T) {
	module, _, _, roleUserRepo, _ := newAdminActionTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users/roles/delete", strings.NewReader("_csrf=test-token&id_user=7&id_role=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: "test-token"})
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/users" {
		t.Fatalf("expected redirect to /users, got status=%d location=%q", rec.Code, rec.Header().Get("Location"))
	}
	if roleUserRepo.removedUser != 7 || roleUserRepo.removedRole != 1 {
		t.Fatalf("expected removed role from user, got user=%d role=%d", roleUserRepo.removedUser, roleUserRepo.removedRole)
	}
}

func TestMountRevokePermissionFromRoleActionForAdmin(t *testing.T) {
	module, _, _, _, permRoleRepo := newAdminActionTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/roles/permissions/delete", strings.NewReader("_csrf=test-token&id_role=1&id_permission=1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: "test-token"})
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/permissions" {
		t.Fatalf("expected redirect to /permissions, got status=%d location=%q", rec.Code, rec.Header().Get("Location"))
	}
	if permRoleRepo.revokedRoleID != 1 || permRoleRepo.revokedPermissionID != 1 {
		t.Fatalf("expected revoked permission from role, got role=%d permission=%d", permRoleRepo.revokedRoleID, permRoleRepo.revokedPermissionID)
	}
}

func TestMountCompanyDeleteAction(t *testing.T) {
	module, factory := newCompanyDeleteActionTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/companies/21/delete", strings.NewReader("_csrf=test-token"))
	req.SetPathValue("id", "21")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: "test-token"})
	req.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/companies" {
		t.Fatalf("expected redirect to /companies, got status=%d location=%q", rec.Code, rec.Header().Get("Location"))
	}
	if factory.company.deletedCompanyID != 21 {
		t.Fatalf("expected company delete for id 21, got %d", factory.company.deletedCompanyID)
	}
	if factory.addressLink.deletedCompanyID != 21 || factory.phoneLink.deletedCompanyID != 21 || factory.emailLink.deletedCompanyID != 21 || factory.socialMedia.deletedCompanyID != 21 {
		t.Fatalf("expected relationship cleanup for company 21")
	}
}

func TestSmokeJourneyAuthenticatedSessionManagement(t *testing.T) {
	module, sessionRepo := newAuthenticatedActionTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	sessionsRec := httptest.NewRecorder()
	sessionsReq := httptest.NewRequest(http.MethodGet, "/sessions", nil)
	sessionsReq.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(sessionsRec, sessionsReq)

	if sessionsRec.Code != http.StatusOK {
		t.Fatalf("expected sessions page status %d, got %d", http.StatusOK, sessionsRec.Code)
	}
	if !strings.Contains(sessionsRec.Body.String(), "Minhas sessoes") {
		t.Fatalf("expected sessions page body, got %q", sessionsRec.Body.String())
	}

	csrfToken := cookieValue(sessionsRec, webcsrf.CookieName)
	if csrfToken == "" {
		t.Fatalf("expected csrf cookie after sessions page render")
	}

	revokeRec := httptest.NewRecorder()
	revokeReq := httptest.NewRequest(http.MethodPost, "/sessions/revoke", strings.NewReader("_csrf="+csrfToken+"&id_session=11"))
	revokeReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	revokeReq.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: csrfToken})
	revokeReq.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(revokeRec, revokeReq)

	if revokeRec.Code != http.StatusSeeOther || revokeRec.Header().Get("Location") != "/login" {
		t.Fatalf("expected redirect to /login, got status=%d location=%q", revokeRec.Code, revokeRec.Header().Get("Location"))
	}
	if len(sessionRepo.revokedIDs) == 0 || sessionRepo.revokedIDs[len(sessionRepo.revokedIDs)-1] != 11 {
		t.Fatalf("expected revoke of session 11, got %#v", sessionRepo.revokedIDs)
	}
}

func TestSmokeJourneyAdminRBACManagement(t *testing.T) {
	module, roleRepo, _, roleUserRepo, _ := newAdminActionTestModule(t)
	mux := http.NewServeMux()
	module.Mount(mux)

	rolesRec := httptest.NewRecorder()
	rolesReq := httptest.NewRequest(http.MethodGet, "/roles", nil)
	rolesReq.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(rolesRec, rolesReq)

	if rolesRec.Code != http.StatusOK {
		t.Fatalf("expected roles page status %d, got %d", http.StatusOK, rolesRec.Code)
	}
	if !strings.Contains(rolesRec.Body.String(), "Catalogo de roles") {
		t.Fatalf("expected roles page body, got %q", rolesRec.Body.String())
	}

	csrfToken := cookieValue(rolesRec, webcsrf.CookieName)
	if csrfToken == "" {
		t.Fatalf("expected csrf cookie after roles page render")
	}

	createRoleRec := httptest.NewRecorder()
	createRoleReq := httptest.NewRequest(http.MethodPost, "/roles", strings.NewReader("_csrf="+csrfToken+"&name=operator&description=Operador"))
	createRoleReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	createRoleReq.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: csrfToken})
	createRoleReq.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(createRoleRec, createRoleReq)

	if createRoleRec.Code != http.StatusSeeOther || createRoleRec.Header().Get("Location") != "/roles" {
		t.Fatalf("expected redirect to /roles, got status=%d location=%q", createRoleRec.Code, createRoleRec.Header().Get("Location"))
	}
	if roleRepo.createdName != "operator" {
		t.Fatalf("expected created role operator, got %q", roleRepo.createdName)
	}

	assignRoleRec := httptest.NewRecorder()
	assignRoleReq := httptest.NewRequest(http.MethodPost, "/users/roles", strings.NewReader("_csrf="+csrfToken+"&id_user=7&id_role=1"))
	assignRoleReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	assignRoleReq.AddCookie(&http.Cookie{Name: webcsrf.CookieName, Value: csrfToken})
	assignRoleReq.AddCookie(&http.Cookie{Name: "actajus_access_token", Value: "valid-token"})
	mux.ServeHTTP(assignRoleRec, assignRoleReq)

	if assignRoleRec.Code != http.StatusSeeOther || assignRoleRec.Header().Get("Location") != "/users" {
		t.Fatalf("expected redirect to /users, got status=%d location=%q", assignRoleRec.Code, assignRoleRec.Header().Get("Location"))
	}
	if roleUserRepo.assignedUser != 7 || roleUserRepo.assignedRole != 1 {
		t.Fatalf("expected assigned user/role, got user=%d role=%d", roleUserRepo.assignedUser, roleUserRepo.assignedRole)
	}
}
