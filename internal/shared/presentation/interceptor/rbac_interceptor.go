// Package interceptor
package interceptor

import (
	"context"

	"connectrpc.com/connect"
	sharedAdapter "github.com/paladignus/actajus/internal/shared/presentation/adapter"
	"github.com/paladignus/actajus/internal/shared/presentation/authctx"
)

type PermissionChecker interface {
	HasPermission(ctx context.Context, idUser int64, permission string) (bool, error)
}

type RBACInterceptor struct {
	checker PermissionChecker
	rules   map[string]string
}

type RBACRule map[string]string

func NewRBACInterceptor(checker PermissionChecker, rules RBACRule) *RBACInterceptor {
	return &RBACInterceptor{
		checker: checker,
		rules:   rules,
	}
}

func (i *RBACInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		perm, ok := i.rules[req.Spec().Procedure]
		if !ok || perm == "" {
			return next(ctx, req)
		}
		claims, ok := authctx.GetClaims(ctx)
		if !ok {
			return nil, sharedAdapter.ToConnectIdentityError(errUnauthenticated())
		}
		allowed, err := i.checker.HasPermission(ctx, claims.IDUser, perm)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		if !allowed {
			return nil, sharedAdapter.ToConnectIdentityError(errPermissionDenied())
		}
		return next(ctx, req)
	}
}

func errUnauthenticated() error { return connect.NewError(connect.CodeUnauthenticated, nil) }

func errPermissionDenied() error {
	return connect.NewError(connect.CodePermissionDenied, connect.NewError(connect.CodePermissionDenied, nil))
}

func (i *RBACInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i *RBACInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}
