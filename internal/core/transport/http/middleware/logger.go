package core_transport_http_middleware

import (
	"context"
	"net/http"

	core_logger "github.com/romreign/AuthService/internal/core/logger"
)

type contextKey string

const loggerKey contextKey = "logger"

func Logger(logger *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			log := logger.With(
				"request_id", requestID,
				"url", r.URL.String(),
				"method", r.Method,
			)

			ctx := context.WithValue(r.Context(), loggerKey, log)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
