package core_transport_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/romreign/AuthService/docs"
	core_logger "github.com/romreign/AuthService/internal/core/logger"
	core_transport_http_metrics "github.com/romreign/AuthService/internal/core/transport/http/metrics"
	core_transport_http_middleware "github.com/romreign/AuthService/internal/core/transport/http/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type HTTPServer struct {
	mux *http.ServeMux
	cfg Config
	log *core_logger.Logger

	middleware []core_transport_http_middleware.Middleware
}

func NewHTTPServer(
	cfg Config,
	log *core_logger.Logger,
	middleware ...core_transport_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		cfg:        cfg,
		log:        log,
		middleware: middleware,
	}
}

func (h *HTTPServer) RegisterAPIRoutes(routes ...*APIVersionRouter) {
	for _, route := range routes {
		prefix := "/api/" + string(route.apiVersion)
		h.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, route),
		)
	}
}

func (h *HTTPServer) RegisterSwagger() {
	h.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		),
	)

	h.mux.HandleFunc(
		"/swagger/doc.json",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)
}

func (h *HTTPServer) Listen(ctx context.Context) error {
	mux := core_transport_http_middleware.ChainMiddleware(h.mux, h.middleware...)

	server := http.Server{
		Addr:    h.cfg.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		if err := core_transport_http_metrics.Listen("0.0.0.0:8081"); err != nil {
			h.log.Warn("failed to start metrics server", "err", err)
		}
	}()

	go func() {
		defer close(ch)
		h.log.Warn("start HTTP server", "addr", h.cfg.Addr)

		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and server HTTP: %w", err)
		}
	case <-ctx.Done():
		h.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			h.cfg.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
	}

	h.log.Warn("http server stopped")
	return nil
}
