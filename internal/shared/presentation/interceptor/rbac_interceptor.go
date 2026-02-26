// Package interceptor
package interceptor

import (
	"context"

	"connectrpc.com/connect"
	"github.com/paladignus/actajus/internal/shared/presentation/authctx"
)

type PermissionChecker interface {
	HasPermission(ctx context.Context, userID int64, perm string) (bool, error)
}

type RBACInterceptor struct {
	checker PermissionChecker
	rules   map[string]string // procedure -> perm
}

func NewRBACInterceptor(checker PermissionChecker, rules map[string]string) *RBACInterceptor {
	return &RBACInterceptor{checker: checker, rules: rules}
}

func (i *RBACInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if err := i.authorize(ctx, req.Spec().Procedure); err != nil {
			return nil, err
		}
		return next(ctx, req)
	}
}

func (i *RBACInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i *RBACInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, shc connect.StreamingHandlerConn) error {
		if err := i.authorize(ctx, shc.Spec().Procedure); err != nil {
			return err
		}
		return next(ctx, shc)
	}
}

func (i *RBACInterceptor) authorize(ctx context.Context, procedure string) error {
	perm, ok := i.rules[procedure]
	if !ok || perm == "" {
		return nil
	}
	claims, ok := authctx.GetClaims(ctx)
	if !ok {
		return connect.NewError(connect.CodeUnauthenticated, nil)
	}
	allowed, err := i.checker.HasPermission(ctx, claims.IDUser, perm)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}
	if !allowed {
		return connect.NewError(connect.CodePermissionDenied, nil)
	}
	return nil
}
