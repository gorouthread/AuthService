package core_transport_http_middleware

import (
	"net/http"
	"time"

	core_logger "github.com/romreign/AuthService/internal/core/logger"
	core_transport_http_response "github.com/romreign/AuthService/internal/core/transport/http/response"
)

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_transport_http_response.NewResponseWriter(w)

			before := time.Now()

			log.Debug(
				">>> incoming HTTP request",
				"time", before.UTC(),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				">>> done HTTP request",
				"status_code", rw.StatusCode(),
				"latency", time.Since(before),
			)
		})
	}
}
