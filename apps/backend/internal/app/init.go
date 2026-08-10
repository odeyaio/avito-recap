package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	httpadapter "avito-recap/internal/adapter/in/http"
	postgresrepo "avito-recap/internal/adapter/out/repository/postgres"
	"avito-recap/internal/engine"
	"avito-recap/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func New(ctx context.Context, config Config, logger *slog.Logger) (*App, error) {
	const op = "app.New"

	recapEngine, err := engine.New(config.Engine)
	if err != nil {
		return nil, fmt.Errorf("%s: configure engine: %w", op, err)
	}

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

	profileRepository := postgresrepo.NewProfileRepository(pool)
	datasetRepository := postgresrepo.NewDatasetRepository(pool)
	catalogRepository := postgresrepo.NewCatalogRepository(pool)
	recapRepository := postgresrepo.NewRecapRepository(pool)
	profileService := service.NewProfileService(profileRepository)
	recapService := service.NewRecapService(
		profileRepository,
		datasetRepository,
		catalogRepository,
		recapRepository,
		recapEngine,
		service.DefaultActionResolver{},
	)

	router := echo.New()
	router.Logger = logger
	router.Use(middleware.RequestID(), middleware.RequestLogger(), middleware.Recover())
	httpadapter.RegisterHandlers(router, pool, profileService, recapService)

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

	return &App{
		server:   server,
		router:   router,
		pool:     pool,
		profiles: profileService,
		recaps:   recapService,
		logger:   logger,
	}, nil
}
