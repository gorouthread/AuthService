package core_logger

import (
	"context"
	"log/slog"
	"os"
)

type contextKey string

const loggerKey contextKey = "logger"

type Logger struct {
	*slog.Logger
}

func NewLogger(cfg Config) *Logger {
	var log *slog.Logger

	switch cfg.Env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level:     slog.LevelDebug,
					AddSource: true,
				},
			),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level:     slog.LevelDebug,
					AddSource: true,
				},
			),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level:     slog.LevelInfo,
					AddSource: false,
				},
			),
		)
	default:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level:     slog.LevelDebug,
					AddSource: true,
				},
			),
		)
	}

	log = log.With(
		"service", cfg.ServiceName,
		"environment", cfg.Env,
		"version", cfg.Version,
	)

	return &Logger{
		Logger: log,
	}
}

func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value(loggerKey).(*Logger)
	if !ok {
		return NewLogger(Config{
			Env:         envLocal,
			ServiceName: "unknown",
			Version:     "unknown",
		})
	}

	return log
}
