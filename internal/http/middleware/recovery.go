package middleware

import (
	"net/http"

	"github.com/bkiran6398/library/internal/http/response"
	"github.com/rs/zerolog"
)

func Recovery(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error().Interface("panic", rec).Str("request_id", GetRequestID(r.Context())).Msg("panic recovered")
					response.Error(w, http.StatusInternalServerError, "internal_error", "Internal server error", nil)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
