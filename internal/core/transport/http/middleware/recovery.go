package core_transport_http_middleware

import (
	"net/http"

	core_logger "github.com/romreign/AuthService/internal/core/logger"
	core_transport_http_response "github.com/romreign/AuthService/internal/core/transport/http/response"
)

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_transport_http_response.NewResponseWriter(w)
			rh := core_transport_http_response.NewHTTPResponseHandler(rw, log)

			defer func() {
				if p := recover(); p != nil {
					rh.PanicResponse(p, "http request got unexpected panic")
				}
			}()

			next.ServeHTTP(w, r)
		},
		)
	}
}
