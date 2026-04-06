// Package interceptor
package interceptor

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	identityCommand "github.com/paladignus/actajus/internal/module/identity/application/command"
	"github.com/paladignus/actajus/internal/module/identity/application/usecase"
	"github.com/paladignus/actajus/internal/module/identity/domain"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/presentation/adapter"
	"github.com/paladignus/actajus/internal/shared/presentation/authctx"
)

type AuthInterceptor struct {
	validate  usecase.ValidateAccess
	whitelist map[string]struct{}
}

type AuthInterceptorOption func(*AuthInterceptor)

func WithWhitelistProcedures(procs ...string) AuthInterceptorOption {
	return func(a *AuthInterceptor) {
		for _, p := range procs {
			a.whitelist[p] = struct{}{}
		}
	}
}

func NewAuthInterceptor(
	validate usecase.ValidateAccess,
	opts ...AuthInterceptorOption,
) *AuthInterceptor {
	a := &AuthInterceptor{
		validate:  validate,
		whitelist: map[string]struct{}{},
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *AuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		proc := req.Spec().Procedure
		if _, ok := a.whitelist[proc]; ok {
			return next(ctx, req)
		}
		authz := req.Header().Get("Authorization")
		token := parseBearer(authz)
		if token == "" {
			return nil, sharedAdapter.ToConnectIdentityError(domain.ErrMissingAccessToken)
		}
		claims, err := a.validate.Execute(ctx, identityCommand.ValidateAccessCommand{AccessToken: token})
		if err != nil {
			return nil, sharedAdapter.ToConnectIdentityError(err)
		}
		ctx = authctx.WithClaims(ctx, authctx.Claims{
			IDUser:    claims.IDUser,
			IDSession: claims.IDSession,
		})
		return next(ctx, req)
	}
}

func parseBearer(h string) string {
	h = strings.TrimSpace(h)
	if h == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(h), "bearer ") {
		return strings.TrimSpace(h[7:])
	}
	return ""
}

func (a *AuthInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (a *AuthInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}

// func WithWhitelistProcedures(procedures ...string) AuthInterceptorOption {
// 	return func(a *AuthInterceptor) {
// 		if a.whitelistProcedures == nil {
// 			a.whitelistProcedures = map[string]struct{}{}
// 		}
// 		for _, p := range procedures {
// 			a.whitelistProcedures[p] = struct{}{}
// 		}
// 	}
// }
//
// func WithWhitelistPrefixes(prefixes ...string) AuthInterceptorOption {
// 	return func(a *AuthInterceptor) { a.whitelistPrefixes = append(a.whitelistPrefixes, prefixes...) }
// }
//
// func WithProcedureLogging(enabled bool) AuthInterceptorOption {
// 	return func(a *AuthInterceptor) { a.logProcedures = enabled }
// }
//
// func NewAuthInterceptor(
// 	validate identityUsecase.ValidateAccess,
// 	logger repository.Logger,
// 	opts ...AuthInterceptorOption,
// ) *AuthInterceptor {
// 	ai := &AuthInterceptor{
// 		validate:            validate,
// 		logger:              logger,
// 		whitelistProcedures: map[string]struct{}{},
// 		whitelistPrefixes:   []string{},
// 	}
// 	for _, opt := range opts {
// 		opt(ai)
// 	}
// 	return ai
// }
//
// func (a *AuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
// 	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
// 		proc := req.Spec().Procedure
// 		if a.logProcedures {
// 			a.logger.Info(ctx, "connect procedure", "procedure", proc)
// 		}
// 		if a.isWhitelisted(proc) {
// 			return next(ctx, req)
// 		}
// 		token := extractBearer(req.Header().Get("Authorization"))
// 		if token == "" {
// 			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
// 		}
// 		ac, err := a.validate.Execute(ctx, dto.ValidateAccessCommand{AccessToken: token})
// 		if err != nil {
// 			return nil, adapter.ToConnectIdentityError(err)
// 		}
// 		ctx = authctx.With(ctx, *ac)
// 		return next(ctx, req)
// 	}
// }
//
// func (a *AuthInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
// 	return next
// }
//
// func (a *AuthInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
// 	return func(ctx context.Context, shc connect.StreamingHandlerConn) error {
// 		proc := shc.Spec().Procedure
// 		if a.logProcedures {
// 			a.logger.Info(ctx, "connect procedure", "procedure", proc)
// 		}
// 		if a.isWhitelisted(proc) {
// 			return next(ctx, shc)
// 		}
// 		token := extractBearer(shc.RequestHeader().Get("Authorization"))
// 		if token == "" {
// 			return connect.NewError(connect.CodeUnauthenticated, nil)
// 		}
// 		ac, err := a.validate.Execute(ctx, dto.ValidateAccessCommand{AccessToken: token})
// 		if err != nil {
// 			return adapter.ToConnectIdentityError(err)
// 		}
// 		ctx = authctx.With(ctx, *ac)
// 		return next(ctx, shc)
// 	}
// }
//
// func (a *AuthInterceptor) isWhitelisted(procedure string) bool {
// 	if _, ok := a.whitelistProcedures[procedure]; ok {
// 		return true
// 	}
// 	for _, p := range a.whitelistPrefixes {
// 		if strings.HasPrefix(procedure, p) {
// 			return true
// 		}
// 	}
// 	return false
// }
//
// func extractBearer(v string) string {
// 	v = strings.TrimSpace(v)
// 	if v == "" {
// 		return ""
// 	}
// 	if strings.HasPrefix(strings.ToLower(v), "bearer ") {
// 		return strings.TrimSpace(v[7:])
// 	}
// 	return v
// }
