// Package middleware
package middleware

import (
	"net/http"

	"github.com/paladignus/actajus/internal/infrastructure/adapter"
	appctx "github.com/paladignus/actajus/internal/infrastructure/context"
)

const TraceIDHeader = "X-Trace-ID"

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Tentar pegar trace ID do header (útil para propagar entre serviços)
		traceID := r.Header.Get(TraceIDHeader)

		// Se não existir, gerar um novo
		if traceID == "" {
			traceID = adapter.GenerateTraceID()
		}

		// Adicionar ao contexto
		ctx := appctx.WithTraceID(r.Context(), traceID)

		// Adicionar ao response header para o cliente poder rastrear
		w.Header().Set(TraceIDHeader, traceID)

		// Continuar com o contexto enriquecido
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
