// Package interceptor
package interceptor

import (
	"context"
	"log/slog"
	"strings"

	"connectrpc.com/connect"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
	identityUsecase "github.com/paladignus/actajus/internal/module/identity/application/usecase"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/presentation/adapter"
	"github.com/paladignus/actajus/internal/shared/presentation/authctx"
)

type AuthInterceptor struct {
	validate            identityUsecase.ValidateAccess
	logger              *slog.Logger
	whitelistProcedures map[string]struct{}
	whitelistPrefixes   []string
	logProcedures       bool
}

type AuthInterceptorOption func(*AuthInterceptor)

func WithLogger(l *slog.Logger) AuthInterceptorOption {
	return func(a *AuthInterceptor) { a.logger = l }
}

func WithWhitelistProcedures(procedures ...string) AuthInterceptorOption {
	return func(a *AuthInterceptor) {
		if a.whitelistProcedures == nil {
			a.whitelistProcedures = map[string]struct{}{}
		}
		for _, p := range procedures {
			a.whitelistProcedures[p] = struct{}{}
		}
	}
}

func WithWhitelistPrefixes(prefixes ...string) AuthInterceptorOption {
	return func(a *AuthInterceptor) {
		a.whitelistPrefixes = append(a.whitelistPrefixes, prefixes...)
	}
}

func WithProcedureLogging(enabled bool) AuthInterceptorOption {
	return func(a *AuthInterceptor) { a.logProcedures = enabled }
}

func NewAuthInterceptor(validate identityUsecase.ValidateAccess, opts ...AuthInterceptorOption) *AuthInterceptor {
	ai := &AuthInterceptor{
		validate:            validate,
		whitelistProcedures: map[string]struct{}{},
		whitelistPrefixes:   []string{},
	}
	for _, opt := range opts {
		opt(ai)
	}
	return ai
}

func (a *AuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		proc := req.Spec().Procedure
		if a.logProcedures {
			a.log("connect procedure", slog.String("procedure", proc))
		}
		if a.isWhitelisted(proc) {
			return next(ctx, req)
		}
		token := extractBearer(req.Header().Get("Authorization"))
		if token == "" {
			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
		}
		ac, err := a.validate.Execute(ctx, dto.ValidateAccessCommand{AccessToken: token})
		if err != nil {
			return nil, sharedAdapter.ToConnectIdentityError(err)
		}
		ctx = authctx.With(ctx, *ac)
		return next(ctx, req)
	}
}

func (a *AuthInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (a *AuthInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, shc connect.StreamingHandlerConn) error {
		proc := shc.Spec().Procedure
		if a.logProcedures {
			a.log("connect procedure", slog.String("procedure", proc))
		}
		if a.isWhitelisted(proc) {
			return next(ctx, shc)
		}
		token := extractBearer(shc.RequestHeader().Get("Authorization"))
		if token == "" {
			return connect.NewError(connect.CodeUnauthenticated, nil)
		}
		ac, err := a.validate.Execute(ctx, dto.ValidateAccessCommand{AccessToken: token})
		if err != nil {
			return sharedAdapter.ToConnectIdentityError(err)
		}
		ctx = authctx.With(ctx, *ac)
		return next(ctx, shc)
	}
}

func (a *AuthInterceptor) isWhitelisted(procedure string) bool {
	if _, ok := a.whitelistProcedures[procedure]; ok {
		return true
	}
	for _, p := range a.whitelistPrefixes {
		if strings.HasPrefix(procedure, p) {
			return true
		}
	}
	return false
}

func extractBearer(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	token, ok := strings.CutPrefix(v, strings.ToLower("Bearer "))
	if ok {
		return token
	}
	// low := strings.ToLower(v)
	// if strings.HasPrefix(low, "bearer ") {
	// 	return strings.TrimSpace(v[7:])
	// }
	return v
}

func (a *AuthInterceptor) log(msg string, attrs ...slog.Attr) {
	if a.logger != nil {
		a.logger.LogAttrs(context.Background(), slog.LevelDebug, msg, attrs...)
	}
}

// type AuthInterceptor struct {
// 	validate  identityUsecase.ValidateAccess
// 	whitelist map[string]struct{}
// }
//
// type AuthInterceptorOption func(*AuthInterceptor)
//
// func WithWhitelist(procedures ...string) AuthInterceptorOption {
// 	return func(a *AuthInterceptor) {
// 		if a.whitelist == nil {
// 			a.whitelist = map[string]struct{}{}
// 		}
// 		for _, p := range procedures {
// 			a.whitelist[p] = struct{}{}
// 		}
// 	}
// }
//
// func NewAuthInterceptor(validate identityUsecase.ValidateAccess, opts ...AuthInterceptorOption) *AuthInterceptor {
// 	ai := &AuthInterceptor{
// 		validate:  validate,
// 		whitelist: map[string]struct{}{},
// 	}
// 	for _, opt := range opts {
// 		opt(ai)
// 	}
// 	return ai
// }
//
// func (a *AuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
// 	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
// 		if a.isWhitelisted(req.Spec().Procedure) {
// 			return next(ctx, req)
// 		}
// 		token := extractBearer(req.Header().Get("Authorization"))
// 		if token == "" {
// 			return nil, connect.NewError(connect.CodeUnauthenticated, nil)
// 		}
// 		ac, err := a.validate.Execute(ctx, dto.ValidateAccessCommand{AccessToken: token})
// 		if err != nil {
// 			return nil, sharedAdapter.ToConnectIdentityError(err)
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
// 		if a.isWhitelisted(shc.Spec().Procedure) {
// 			return next(ctx, shc)
// 		}
// 		token := extractBearer(shc.RequestHeader().Get("Authorization"))
// 		if token == "" {
// 			return connect.NewError(connect.CodeUnauthenticated, nil)
// 		}
// 		ac, err := a.validate.Execute(ctx, dto.ValidateAccessCommand{AccessToken: token})
// 		if err != nil {
// 			return sharedAdapter.ToConnectIdentityError(err)
// 		}
// 		ctx = authctx.With(ctx, *ac)
// 		return next(ctx, shc)
// 	}
// }
//
// func (a *AuthInterceptor) isWhitelisted(procedure string) bool {
// 	_, ok := a.whitelist[procedure]
// 	return ok
// }
//
// func extractBearer(v string) string {
// 	v = strings.TrimSpace(v)
// 	if v == "" {
// 		return ""
// 	}
// 	token, ok := strings.CutPrefix(v, strings.ToLower("Bearer "))
// 	if ok {
// 		return token
// 	}
// 	return v
// }
