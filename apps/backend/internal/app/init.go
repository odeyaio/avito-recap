package app

import (
	httpadapter "avito-recap/internal/adapter/in/http"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func New(ctx context.Context, config Config, logger *slog.Logger) (*App, error) {
	const op = "app.New"

	poolConfig, err := pgxpool.ParseConfig(config.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	router := echo.New()
	router.Logger = logger
	router.Use(middleware.RequestID(), middleware.RequestLogger(), middleware.Recover())
	httpadapter.RegisterHandlers(router, pool)

	server := echo.StartConfig{
		Address:         config.HTTP.Address,
		HideBanner:      true,
		HidePort:        true,
		GracefulTimeout: config.HTTP.ShutdownTimeout,
		OnShutdownError: func(shutdownErr error) {
			logger.Error("graceful shutdown failed", "error", shutdownErr)
		},
		BeforeServeFunc: func(server *http.Server) error {
			server.ReadHeaderTimeout = config.HTTP.ReadHeaderTimeout
			server.ReadTimeout = config.HTTP.ReadTimeout
			server.WriteTimeout = config.HTTP.WriteTimeout
			server.IdleTimeout = config.HTTP.IdleTimeout
			return nil
		},
	}

	return &App{server: server, router: router, pool: pool, logger: logger}, nil
}
