// Package authctx
package authctx

import (
	"context"

	"github.com/paladignus/actajus/internal/module/identity/application/dto"
)

type key struct{}

func With(ctx context.Context, ac dto.AuthContext) context.Context {
	return context.WithValue(ctx, key{}, ac)
}

func Get(ctx context.Context) (dto.AuthContext, bool) {
	v := ctx.Value(key{})
	if v == nil {
		return dto.AuthContext{}, false
	}
	ac, ok := v.(dto.AuthContext)
	return ac, ok
}

func Must(ctx context.Context) dto.AuthContext {
	ac, ok := Get(ctx)
	if !ok {
		panic("authctx: missing AuthContext in context")
	}
	return ac
}

func UserID(ctx context.Context) (int64, bool) {
	ac, ok := Get(ctx)
	if !ok {
		return 0, false
	}
	return ac.IDUser, true
}

func SessionID(ctx context.Context) (int64, bool) {
	ac, ok := Get(ctx)
	if !ok {
		return 0, false
	}
	return ac.IDSession, true
}
