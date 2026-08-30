package core_transport_http_metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	core_transport_http_middleware "github.com/romreign/AuthService/internal/core/transport/http/middleware"
	core_transport_http_response "github.com/romreign/AuthService/internal/core/transport/http/response"
)

var requestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests.",
	},
	[]string{"method", "path", "status"},
)

func init() {
	prometheus.MustRegister(requestsTotal)
}

func Metrics() core_transport_http_middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := core_transport_http_response.NewResponseWriter(w)

			next.ServeHTTP(rw, r)

			requestsTotal.WithLabelValues(
				r.Method,
				r.URL.Path,
				strconv.Itoa(rw.StatusCode()),
			).Inc()
		})
	}
}

func Listen(host string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(host, mux)
}
