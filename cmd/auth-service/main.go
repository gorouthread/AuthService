package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/romreign/AuthService/internal/core/logger"
	core_repository_postgres "github.com/romreign/AuthService/internal/core/repository/postgres"
	core_repository_redis "github.com/romreign/AuthService/internal/core/repository/redis"
	core_transport_http_metrics "github.com/romreign/AuthService/internal/core/transport/http/metrics"
	core_transport_http_middleware "github.com/romreign/AuthService/internal/core/transport/http/middleware"
	core_transport_http_server "github.com/romreign/AuthService/internal/core/transport/http/server"
	auth_repository_postgres "github.com/romreign/AuthService/internal/repository/postgres"
	auth_repository_redis "github.com/romreign/AuthService/internal/repository/redis"
	auth_service "github.com/romreign/AuthService/internal/service/auth"
	idempotency_service "github.com/romreign/AuthService/internal/service/idempotency"
	auth_transport_http "github.com/romreign/AuthService/internal/transport/http"
	core_service_jwt "github.com/romreign/AuthService/pkg/jwt"

	_ "github.com/romreign/AuthService/docs"
)

// @title        Golang Authorization API
// @version      1.0
// @description  Authorization Application REST-API schema
// @host         127.0.0.1:8080
// @BasePath     /api/v1
func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	cfgLogger := core_logger.NewConfigMust()
	logger := core_logger.NewLogger(cfgLogger)

	logger.Debug("initializing postgres connection pool")
	cfgPgpool := core_repository_postgres.NewConfigMust()
	pgpool, err := core_repository_postgres.NewConnectionPool(ctx, cfgPgpool)
	if err != nil {
		log.Fatal("failed to init postgres connection pool", err)
	}
	defer pgpool.Close()

	logger.Debug("initializing redis connection pool")
	cfgRdpool := core_repository_redis.NewConfigMust()
	rdpool, err := core_repository_redis.NewConnectionPool(ctx, cfgRdpool)
	if err != nil {
		log.Fatal("failed to init redis connection pool", err)
	}
	defer rdpool.Close()

	logger.Debug("initializing jwt manager")
	cfgJWT := core_service_jwt.NewConfigMust()
	manager := core_service_jwt.NewJWTMaker(cfgJWT)

	logger.Debug("initializing layers")
	pgRepo := auth_repository_postgres.NewAuthRepositoryPostgres(pgpool)
	rdRepo := auth_repository_redis.NewAuthRepositoryRedis(rdpool)
	authSrvc := auth_service.NewAuthService(&pgRepo, manager)
	idempSrvc := idempotency_service.NewIdempotencyService(rdRepo, 24*time.Hour)
	authHandler := auth_transport_http.NewAuthHTTPHandler(authSrvc, idempSrvc)

	logger.Debug("initializing http server")
	cfgSrv := core_transport_http_server.NewConfigMust()
	httpServer := core_transport_http_server.NewHTTPServer(
		cfgSrv,
		logger,
		core_transport_http_metrics.Metrics(),
		core_transport_http_middleware.CORS(nil), // пока нет домена
		core_transport_http_middleware.RequestID(),
		core_transport_http_middleware.Logger(logger),
		core_transport_http_middleware.Trace(),
		core_transport_http_middleware.Panic(),
	)
	apiVersionRouter := core_transport_http_server.NewAPIVersionRouter(core_transport_http_server.APIVersion1)
	apiVersionRouter.RegisterRoutes(authHandler.Routes()...)
	httpServer.RegisterAPIRoutes(apiVersionRouter)
	httpServer.RegisterSwagger()

	if err := httpServer.Listen(ctx); err != nil {
		logger.Error("http server listen", "error:", err.Error())
	}
}
