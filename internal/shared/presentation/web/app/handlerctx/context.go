package handlerctx

import (
	"net/http"
	"time"

	identityread "github.com/paladignus/actajus/internal/module/identity/application/readmodel"
	sharedRepo "github.com/paladignus/actajus/internal/shared/application/repository"
	webauth "github.com/paladignus/actajus/internal/shared/presentation/web/app/adapter/auth"
	webhtml "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/html"
	navigation "github.com/paladignus/actajus/internal/shared/presentation/web/app/response/http/navigation"
)

type State = webauth.State

type Context struct {
	Logger             sharedRepo.Logger
	HTML               *webhtml.Responder
	Navigator          *navigation.Navigator
	AuthAdapter        *webauth.Adapter
	FlashCookie        string
	AuthFromRequest    func(http.ResponseWriter, *http.Request) (*State, error)
	RequireAuth        func(http.ResponseWriter, *http.Request) (*State, error)
	RequireRBACAdmin   func(http.ResponseWriter, *http.Request) (*State, error)
	SetAuthCookies     func(http.ResponseWriter, *identityread.AuthTokensReadModel)
	ClearAuthCookies   func(http.ResponseWriter)
	RespondForbidden   func(http.ResponseWriter, *http.Request)
	RespondNotFound    func(http.ResponseWriter, *http.Request)
	CheckLoginAttempt  func(*http.Request, string) (bool, time.Duration)
	RecordLoginFailure func(*http.Request, string)
	RecordLoginSuccess func(*http.Request, string)
	AuditSuccess       func(*http.Request, string, ...any)
	AuditFailure       func(*http.Request, string, ...any)
}
